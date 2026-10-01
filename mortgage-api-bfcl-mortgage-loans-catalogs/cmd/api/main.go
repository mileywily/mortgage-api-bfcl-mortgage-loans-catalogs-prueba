package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	gintrace "github.com/DataDog/dd-trace-go/contrib/gin-gonic/gin/v2"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/DataDog/dd-trace-go/v2/profiler"

	headerForward "github.com/falabella-regulado/go-lib-gin-fif/headers"
	loggerMiddleware "github.com/falabella-regulado/go-lib-gin-fif/logger"
	recoveryMiddleware "github.com/falabella-regulado/go-lib-gin-fif/panic_recovery"
	"github.com/falabella-regulado/go-lib-logger-fif/loggers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"mortgage-api-bfcl-mortgage-loans-catalogs/cmd/config"
	appCatalog "mortgage-api-bfcl-mortgage-loans-catalogs/internal/app/catalog/get"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/usecases"
	catalogHandler "mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/http/gin/catalog/get"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/rest/catalog_repository"

	_ "mortgage-api-bfcl-mortgage-loans-catalogs/docs"
)

func main() {
	// Carga y valida la configuración de la API
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("error cargando configuracion: %v", err))
	}

	// Inicializa el logger estructurado en formato JSON
	logger := loggers.NewJsonSlogLogger(cfg.LoggingLevel, nil)

	// Configuración e inicio del Tracer (DataDog v2)
	if err := tracer.Start(
		tracer.WithService(cfg.DDServiceName),
		tracer.WithEnv(cfg.Env),
		tracer.WithServiceVersion(cfg.Version),
		tracer.WithAgentAddr(cfg.DataDog.AgentAddr),
		tracer.WithRuntimeMetrics(),
	); err != nil {
		logger.Error("datadog tracer start failed", "error", err.Error())
	}
	defer tracer.Stop()

	// Configuración del Profiler de DataDog si está habilitado
	if cfg.DataDog.ProfilerEnabled {
		if err := profiler.Start(
			profiler.WithService(cfg.DDServiceName),
			profiler.WithEnv(cfg.Env),
			profiler.WithVersion(cfg.Version),
			profiler.WithAgentAddr(cfg.DataDog.AgentAddr),
			profiler.WithProfileTypes(
				profiler.CPUProfile,
				profiler.HeapProfile,
			),
		); err != nil {
			logger.Error("datadog profiler start failed", "error", err.Error())
		} else {
			defer profiler.Stop()
		}
	}

	// Inyección de los backends disponibles; el caso de uso aplica default y override por request.
	dummyRepository := catalog_repository.NewDummyRepository()
	finnflowRepository := catalog_repository.NewRestRepositoryWithTimeout(
		cfg.FinnflowURL,
		"/api/catalogo_detail/",
		true,
		cfg.FinnflowKey,
		cfg.FinnflowSecret,
		cfg.Timeout,
	)
	javaRepository := catalog_repository.NewRestRepositoryWithTimeout(
		cfg.JavaLegacyURL,
		"/v1/bfcl/mortgage-loan/catalogs/",
		false,
		"",
		"",
		cfg.Timeout,
	)
	catalogUsecase := usecases.NewCatalogUsecase(
		domain.CatalogBackend(cfg.DefaultBackend),
		dummyRepository,
		finnflowRepository,
		javaRepository,
	)
	catalogController := appCatalog.NewController(
		catalogUsecase,
		appCatalog.NewRequestMapper(),
		appCatalog.NewResponseMapper(),
	)
	getCatalogHandler := catalogHandler.NewHandler(catalogController, catalogHandler.NewDecoder())

	// Configuración del motor Gin
	gin.SetMode(cfg.GinMode)
	engine := gin.New()

	// Health check fuera del flujo legacy.
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "timestamp": time.Now().UTC()})
	})

	// Middlewares globales, en el orden obligatorio (ADR-API-001-C §2)
	engine.Use(
		gintrace.Middleware(cfg.DDServiceName),
		recoveryMiddleware.NewFifErrorPanicRecoveryMiddleware(),
		loggerMiddleware.NewFifLoggerMiddleware(logger, loggerMiddleware.IgnorePath("/health")),
		headerForward.ForwardHeadersMiddleware(headerForward.HeaderKeys),
	)

	// Esta ruta conserva el contrato Java; sus headers se propagan sin validación obligatoria.
	engine.POST("/v1/bfcl/mortgage-loan/catalogs/:catalog", getCatalogHandler)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Inicio del servidor HTTP con lectura dinámica de PORT (Nullplatform) y Graceful Shutdown
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("Servidor escuchando en el puerto %s", port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped with error", "error", err.Error())
			os.Exit(1)
		}
	}()

	// Captura de señales SIGTERM/SIGINT para Graceful Shutdown en Nullplatform
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Señal de apagado recibida. Iniciando Graceful Shutdown...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Error("Forzando apagado del servidor tras timeout", "error", err.Error())
	} else {
		logger.Info("Servidor detenido limpiamente.")
	}
}
