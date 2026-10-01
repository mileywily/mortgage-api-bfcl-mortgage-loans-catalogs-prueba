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
	GinMode          = "gin_mode"
	Country          = "country"
	UriPrefix        = "uri_prefix"
	Timeout          = "timeout"
	LoggingLevel     = "logging_level"
	DdProfileEnabled = "dd_profile_enabled"
	DdAgentPort      = "dd_agent_port"
	DDAgentHost      = "dd_agent_host"

	// Backends de catálogos
	FinnflowURL    = "finnflow_url"
	FinnflowKey    = "finnflow_key"
	FinnflowSecret = "finnflow_secret"
	JavaLegacyURL  = "java_legacy_url"
	DefaultBackend = "default_backend"
)

func configEntries() []configFif.ConfigEntry {
	return []configFif.ConfigEntry{
		{VariableName: AppName, Description: "nombre de la api", DefaultValue: "mortgage-api-bfcl-mortgage-loans-catalogs"},
		{VariableName: AppEnv, Description: "Ambiente de ejecución (dev|qa|uat|prod)", DefaultValue: "dev"},
		{VariableName: AppVersion, Description: "Versión de la api, inyectada por CI (APP_VERSION)", DefaultValue: "0.0.0"},
		{VariableName: DdServiceName, Description: "nombre del servicio en datadog", DefaultValue: "mortgage-api-bfcl-mortgage-loans-catalogs"},
		{VariableName: GinMode, Description: "Modo en que inicializa GIN DEBUG|RELEASE", DefaultValue: "DEBUG"},
		{VariableName: LoggingLevel, Description: "Level de detalle de logs", DefaultValue: "info"},
		{VariableName: Timeout, Description: "Timeout HTTP en segundos", DefaultValue: 10},
		{VariableName: UriPrefix, Description: "Prefijo de URL con versión", DefaultValue: "/v1/bfcl/mortgage-loan"},
		{VariableName: Country, Description: "País de la api", DefaultValue: "CL"},
		{VariableName: DDAgentHost, Description: "Host del agente DataDog", DefaultValue: "localhost"},
		{VariableName: DdAgentPort, Description: "Puerto del agente DataDog", DefaultValue: "8126"},
		{VariableName: DdProfileEnabled, Description: "Activa el profiler de DataDog", DefaultValue: false},
		// Backends
		{VariableName: FinnflowURL, Description: "URL base de Finnflow (catálogos reales)", DefaultValue: ""},
		{VariableName: FinnflowKey, Description: "Client ID / usuario de Finnflow (Basic Auth)", DefaultValue: ""},
		{VariableName: FinnflowSecret, Description: "Client Secret / password de Finnflow (Basic Auth)", DefaultValue: ""},
		{VariableName: JavaLegacyURL, Description: "URL del proxy Java legado", DefaultValue: ""},
		{VariableName: DefaultBackend, Description: "Backend por defecto: real | dummy | java", DefaultValue: "real"},
	}
}

// APIConfig contiene la configuración completa de la API.
type APIConfig struct {
	URIPrefix     string
	Timeout       time.Duration
	AppName       string
	Env           string
	Version       string
	DDServiceName string
	DataDog       DataDogConfig
	Country       string
	GinMode       string
	LoggingLevel  loggerFif.LogLevel

	// Backends de catálogos
	FinnflowURL    string
	FinnflowKey    string
	FinnflowSecret string
	JavaLegacyURL  string
	DefaultBackend string
}

// DataDogConfig agrupa la configuración de tracing/profiling de DataDog.
type DataDogConfig struct {
	AgentAddr       string
	ProfilerEnabled bool
}

// Load obtiene y valida la configuración de la API.
// Retorna error si hay variables obligatorias vacías.
func Load() (*APIConfig, error) {
	cfg := configFif.LoadConfig(configEntries())
	ginMode := strings.ToLower(cfg[GinMode].(string))

	logLevel, err := loggerFif.StringToLogLevel(cfg[LoggingLevel].(string))
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %w", LoggingLevel, err)
	}

	finnflowURL := cfg[FinnflowURL].(string)
	defaultBackend := strings.ToLower(cfg[DefaultBackend].(string))

	// Si el backend por defecto es "real", Finnflow URL es obligatoria
	if defaultBackend == "real" && finnflowURL == "" {
		return nil, errors.New(FinnflowURL + " is required when default_backend=real")
	}

	return &APIConfig{
		AppName:        cfg[AppName].(string),
		Env:            cfg[AppEnv].(string),
		Version:        cfg[AppVersion].(string),
		DDServiceName:  cfg[DdServiceName].(string),
		URIPrefix:      cfg[UriPrefix].(string),
		LoggingLevel:   logLevel,
		Timeout:        time.Duration(cfg[Timeout].(int)) * time.Second,
		Country:        cfg[Country].(string),
		GinMode:        ginMode,
		DataDog: DataDogConfig{
			AgentAddr:       cfg[DDAgentHost].(string) + ":" + cfg[DdAgentPort].(string),
			ProfilerEnabled: cfg[DdProfileEnabled].(bool),
		},
		FinnflowURL:    finnflowURL,
		FinnflowKey:    cfg[FinnflowKey].(string),
		FinnflowSecret: cfg[FinnflowSecret].(string),
		JavaLegacyURL:  cfg[JavaLegacyURL].(string),
		DefaultBackend: defaultBackend,
	}, nil
}
