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

	errorHandler "github.com/falabella-regulado/go-lib-gin-fif/error_handler"
	headerForward "github.com/falabella-regulado/go-lib-gin-fif/headers"
	loggerMiddleware "github.com/falabella-regulado/go-lib-gin-fif/logger"
	recoveryMiddleware "github.com/falabella-regulado/go-lib-gin-fif/panic_recovery"
	httpFIF "github.com/falabella-regulado/go-lib-http-fif"
	"github.com/falabella-regulado/go-lib-logger-fif/loggers"

	"api-hello-world/cmd/config"
	appExample "api-hello-world/internal/app/example"
	createApp "api-hello-world/internal/app/example/create"
	"api-hello-world/internal/core/usecases"
	createHandler "api-hello-world/internal/infra/http/gin/example/create"
	"api-hello-world/internal/infra/rest/external_service"
	"api-hello-world/internal/infra/tracerMiddleware"
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
		tracer.WithService(cfg.AppName),
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
			profiler.WithService(cfg.AppName),
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

	// Inicializa el cliente REST con configuraciones base y middlewares
	httpClient := httpFIF.NewRestClient(
		httpFIF.BaseURL(cfg.RestServiceCustomerURL),
		httpFIF.Timeout(cfg.Timeout),
		httpFIF.Logger(logger),
		httpFIF.Header(config.HeaderApplicationID, cfg.AppName),
		httpFIF.Header(config.HeaderContentType, "application/json"),
		httpFIF.WithDatadogTracer(),
		httpFIF.WithMandatoryHeaders(),
	)

	// Configuración del motor Gin
	gin.SetMode(cfg.GinMode)
	engine := gin.New()

	// /health fuera del URIPrefix y antes de middlewares globales: no pasa por gateway ni logger
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "timestamp": time.Now().UTC()})
	})

	// Middlewares globales, en el orden obligatorio (ADR-API-001-C §2)
	engine.Use(
		gintrace.Middleware(cfg.AppName),
		recoveryMiddleware.NewFifErrorPanicRecoveryMiddleware(),
		loggerMiddleware.NewFifLoggerMiddleware(logger, loggerMiddleware.IgnorePath("/health")),
		headerForward.ForwardHeadersMiddleware(headerForward.HeaderKeys),
	)

	// Grupo versionado: acá se validan los headers obligatorios del gateway y se registran errores
	api := engine.Group(cfg.URIPrefix)
	api.Use(
		headerForward.ValidateMandatoryFieldsMiddlewareWithDefaultHeader(),
		errorHandler.NewFifErrorHandlerMiddleware(appExample.NewExampleErrorHandler),
	)

	// Wire-up del feature example (instanciar → decorar → inyectar → controller → handler → rutas)
	externalRepository := external_service.NewExternalServiceRepository(
		httpClient,
		external_service.NewRequestMapper(),
		external_service.NewResponseMapper(),
	)
	tracedExternalRepository := tracerMiddleware.NewExternalServiceRepositoryMiddleware(externalRepository, logger)

	usecase := usecases.NewExampleUsecase(tracedExternalRepository)
	tracedUsecase := tracerMiddleware.NewExampleUsecaseMiddleware(usecase, logger)

	createController := createApp.NewController(tracedUsecase, createApp.NewRequestMapper(), createApp.NewResponseMapper())
	createH := createHandler.NewHandler(createController, createHandler.NewDecoder())

	route := api.Group("/anything")
	route.POST("/:something/example", createH)

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
