package main

import (
	"os"

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
	serverFif "github.com/falabella-regulado/go-lib-server-fif"

	"mortgage-api-bfcl-mortgage-loans-catalogs/cmd/config"
	catalogApp "mortgage-api-bfcl-mortgage-loans-catalogs/internal/app/catalog"
	catalogGet "mortgage-api-bfcl-mortgage-loans-catalogs/internal/app/catalog/get"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/usecases"
	ginGet "mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/http/gin/catalog/get"
	catalogSvc "mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/rest/catalog_service"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/tracerMiddleware"
)

func main() {
	// ── 1. Configuración ──────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		panic(err) // antes del logger: no hay otra salida
	}

	// ── 2. Logger ─────────────────────────────────────────────────────────────
	logger := loggers.NewJsonSlogLogger(cfg.LoggingLevel, nil)

	// ── 3. DataDog Tracer ─────────────────────────────────────────────────────
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

	// ── 4. DataDog Profiler (opcional) ────────────────────────────────────────
	if cfg.DataDog.ProfilerEnabled {
		if err := profiler.Start(
			profiler.WithService(cfg.AppName),
			profiler.WithEnv(cfg.Env),
			profiler.WithVersion(cfg.Version),
			profiler.WithAgentAddr(cfg.DataDog.AgentAddr),
			profiler.WithProfileTypes(profiler.CPUProfile, profiler.HeapProfile),
		); err != nil {
			logger.Error("datadog profiler start failed", "error", err.Error())
		} else {
			defer profiler.Stop()
		}
	}

	// ── 5. Dependency Injection: repositorios ─────────────────────────────────

	// 5a. Dummy repository (datos estáticos, siempre disponible)
	dummyRepo := catalogSvc.NewDummyRepository()

	// 5b. Finnflow repository (GET + Basic Auth)
	finnflowClient := httpFIF.NewRestClient(
		httpFIF.BaseURL(cfg.FinnflowURL),
		httpFIF.Timeout(cfg.Timeout),
		httpFIF.Logger(logger),
		httpFIF.Header(config.HeaderContentType, "application/json"),
		httpFIF.BasicAuth(cfg.FinnflowKey, cfg.FinnflowSecret),
		httpFIF.WithDatadogTracer(),
	)
	realRepo := catalogSvc.NewRestRepository(finnflowClient, "/api/catalogo_detail/", true)

	// 5c. Java legacy proxy repository (POST + headers X-*)
	javaClient := httpFIF.NewRestClient(
		httpFIF.BaseURL(cfg.JavaLegacyURL),
		httpFIF.Timeout(cfg.Timeout),
		httpFIF.Logger(logger),
		httpFIF.Header(config.HeaderContentType, "application/json"),
		httpFIF.WithDatadogTracer(),
		httpFIF.WithMandatoryHeaders(),
	)
	javaRepo := catalogSvc.NewRestRepository(javaClient, "/v1/bfcl/mortgage-loan/catalogs/", false)

	// 5d. Decorar repositorios con tracing
	tracedDummyRepo := tracerMiddleware.NewCatalogRepositoryMiddleware(dummyRepo, logger)
	tracedRealRepo := tracerMiddleware.NewCatalogRepositoryMiddleware(realRepo, logger)
	tracedJavaRepo := tracerMiddleware.NewCatalogRepositoryMiddleware(javaRepo, logger)

	// 5e. Casos de uso
	dummyUsecase := usecases.NewCatalogUsecase(tracedDummyRepo)
	realUsecase := usecases.NewCatalogUsecase(tracedRealRepo)
	javaUsecase := usecases.NewCatalogUsecase(tracedJavaRepo)

	// 5f. Seleccionar backend por defecto
	defaultUsecase := realUsecase
	switch cfg.DefaultBackend {
	case "dummy":
		defaultUsecase = dummyUsecase
	case "java":
		defaultUsecase = javaUsecase
	}

	// 5g. Decorar usecases con tracing
	tracedDefaultUsecase := tracerMiddleware.NewCatalogUsecaseMiddleware(defaultUsecase, logger)
	tracedDummyUsecase := tracerMiddleware.NewCatalogUsecaseMiddleware(dummyUsecase, logger)
	tracedRealUsecase := tracerMiddleware.NewCatalogUsecaseMiddleware(realUsecase, logger)
	tracedJavaUsecase := tracerMiddleware.NewCatalogUsecaseMiddleware(javaUsecase, logger)

	// 5h. Controller y Handler
	catalogController := catalogGet.NewController(
		tracedDefaultUsecase,
		tracedDummyUsecase,
		tracedRealUsecase,
		tracedJavaUsecase,
		catalogGet.NewResponseMapper(),
	)
	catalogHandler := ginGet.NewHandler(catalogController, ginGet.NewDecoder())

	// ── 6. Gin Engine ─────────────────────────────────────────────────────────
	gin.SetMode(cfg.GinMode)
	engine := gin.New()

	// /health: fuera del URIPrefix y sin middlewares corporativos
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Middlewares globales (orden obligatorio ADR-API-001-C §2)
	engine.Use(
		gintrace.Middleware(cfg.AppName),
		recoveryMiddleware.NewFifErrorPanicRecoveryMiddleware(),
		loggerMiddleware.NewFifLoggerMiddleware(logger, loggerMiddleware.IgnorePath("/health")),
		headerForward.ForwardHeadersMiddleware(headerForward.HeaderKeys),
	)

	// ── 7. Rutas versionadas ──────────────────────────────────────────────────
	api := engine.Group(cfg.URIPrefix)
	api.Use(
		headerForward.ValidateMandatoryFieldsMiddlewareWithDefaultHeader(),
		errorHandler.NewFifErrorHandlerMiddleware(catalogApp.NewCatalogErrorHandler),
	)

	// POST /v1/bfcl/mortgage-loan/catalogs/:catalog  (paridad con el legado Java)
	api.POST("/catalogs/:catalog", catalogHandler)

	// ── 8. Servidor HTTP con graceful shutdown ────────────────────────────────
	if err := serverFif.StartServer(serverFif.ServerOptions{Logger: logger, Handler: engine}); err != nil {
		logger.Error("server stopped with error", "error", err.Error())
		os.Exit(1)
	}
}
