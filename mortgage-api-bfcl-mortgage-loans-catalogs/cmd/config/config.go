package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	configFif "github.com/falabella-regulado/go-lib-config-fif"
	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"
)

const (
	AppName          = "app_name"
	AppEnv           = "app_env"
	AppVersion       = "app_version"
	DdServiceName    = "dd_service_name"
	FinnflowURL      = "finnflow_url"
	FinnflowKey      = "finnflow_key"
	FinnflowSecret   = "finnflow_secret"
	JavaLegacyURL    = "java_legacy_url"
	DefaultBackend   = "default_backend"
	GinMode          = "gin_mode"
	Country          = "country"
	UriPrefix        = "uri_prefix"
	Timeout          = "timeout"
	LoggingLevel     = "logging_level"
	DdProfileEnabled = "dd_profile_enabled"
	DdAgentPort      = "dd_agent_port"
	DDAgentHost      = "dd_agent_host"
	Port             = "port"
)

// configEntries retorna configuraciones iniciales (env y flags)
func configEntries() []configFif.ConfigEntry {

	return []configFif.ConfigEntry{
		{
			VariableName: AppName,
			Description:  "nombre de la api",
			DefaultValue: "mortgage-api-bfcl-mortgage-loans-catalogs",
		},
		{
			VariableName: AppEnv,
			Description:  "Ambiente de ejecución (dev|qa|uat|prod)",
			DefaultValue: "dev",
		},
		{
			VariableName: AppVersion,
			Description:  "Versión de la api, inyectada por CI (APP_VERSION)",
			DefaultValue: "0.0.0",
		},
		{
			VariableName: DdServiceName,
			Description:  "nombre del servicio en datadog",
			DefaultValue: "mortgage-api-bfcl-mortgage-loans-catalogs",
		},
		{
			VariableName: GinMode,
			Description:  "Modo en que inicializa GIN DEBUG|RELEASE",
			DefaultValue: "DEBUG",
		},
		{
			VariableName: LoggingLevel,
			Description:  "Level de detalle de logs",
			DefaultValue: "info",
		},
		{
			VariableName: Timeout,
			Description:  "timeout por defecto ",
			DefaultValue: 10,
		},
		{
			VariableName: UriPrefix,
			Description:  "Prefijo de URL con version",
			DefaultValue: "/fifcl/v1",
		},
		{
			VariableName: Country,
			Description:  "Pais de la api",
			DefaultValue: "CL",
		},
		{VariableName: FinnflowURL, Description: "URL base del upstream Finnflow", DefaultValue: ""},
		{VariableName: FinnflowKey, Description: "Client ID de Finnflow", DefaultValue: ""},
		{VariableName: FinnflowSecret, Description: "Client secret de Finnflow", DefaultValue: ""},
		{VariableName: JavaLegacyURL, Description: "URL base del proxy Java legacy", DefaultValue: ""},
		{VariableName: DefaultBackend, Description: "Backend por defecto: real, java o dummy", DefaultValue: "real"},
		{
			VariableName: DDAgentHost,
			Description:  "Especifica el host de datadog",
			DefaultValue: "localhost",
		},
		{
			VariableName: DdAgentPort,
			Description:  "Especifica el ports de datadog",
			DefaultValue: "8126",
		},
		{
			VariableName: DdProfileEnabled,
			Description:  "Especifica si está activado el profiler de datadog",
			DefaultValue: false,
		},
		{
			VariableName: Port,
			Description:  "Puerto TCP en el que escucha la API (Nullplatform)",
			DefaultValue: "8080",
		},
	}
}

// APIConfig struct
type APIConfig struct {
	Port           string
	URIPrefix      string
	Timeout        time.Duration
	AppName        string
	Env            string
	Version        string
	FinnflowURL    string
	FinnflowKey    string
	FinnflowSecret string
	JavaLegacyURL  string
	DefaultBackend string
	DDServiceName  string
	DataDog        DataDogConfig
	Country        string
	GinMode        string
	LoggingLevel   loggerFif.LogLevel
}

// DataDogConfig agrupa la configuración de tracing/profiling de DataDog
type DataDogConfig struct {
	AgentAddr       string
	ProfilerEnabled bool
}

// Load obtiene y valida la configuración de la API.
// Nunca arranca con variables obligatorias vacías: retorna error en su lugar.
func Load() (*APIConfig, error) {
	cfg := configFif.LoadConfig(configEntries())
	ginMode := strings.ToLower(cfg[GinMode].(string))

	logLevel, err := loggerFif.StringToLogLevel(cfg[LoggingLevel].(string))
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %w", LoggingLevel, err)
	}

	finnflowURL := strings.TrimRight(cfg[FinnflowURL].(string), "/")
	javaLegacyURL := strings.TrimRight(cfg[JavaLegacyURL].(string), "/")
	defaultBackend := strings.ToLower(cfg[DefaultBackend].(string))
	if defaultBackend != "dummy" && defaultBackend != "java" && defaultBackend != "real" {
		defaultBackend = "real"
	}
	if defaultBackend == "real" && finnflowURL == "" {
		return nil, errors.New(FinnflowURL + " is required when DEFAULT_BACKEND=real")
	}
	if defaultBackend == "java" && javaLegacyURL == "" {
		return nil, errors.New(JavaLegacyURL + " is required when DEFAULT_BACKEND=java")
	}

	port := "8080"
	if p, ok := cfg[Port].(string); ok && p != "" {
		port = p
	}

	return &APIConfig{
		Port:           port,
		AppName:        cfg[AppName].(string),
		Env:            cfg[AppEnv].(string),
		Version:        cfg[AppVersion].(string),
		FinnflowURL:    finnflowURL,
		FinnflowKey:    cfg[FinnflowKey].(string),
		FinnflowSecret: cfg[FinnflowSecret].(string),
		JavaLegacyURL:  javaLegacyURL,
		DefaultBackend: defaultBackend,
		DDServiceName:  cfg[DdServiceName].(string),
		URIPrefix:      cfg[UriPrefix].(string),
		LoggingLevel:   logLevel,
		Timeout:        time.Duration(cfg[Timeout].(int)) * time.Second,
		DataDog: DataDogConfig{
			AgentAddr:       cfg[DDAgentHost].(string) + ":" + cfg[DdAgentPort].(string),
			ProfilerEnabled: cfg[DdProfileEnabled].(bool),
		},
		Country: cfg[Country].(string),
		GinMode: ginMode,
	}, nil
}
