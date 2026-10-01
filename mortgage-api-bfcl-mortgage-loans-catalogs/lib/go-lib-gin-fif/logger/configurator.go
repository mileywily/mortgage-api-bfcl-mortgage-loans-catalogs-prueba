package logger

import (
	configFif "github.com/falabella-regulado/go-lib-config-fif"
	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"
)

const (
	DefaultLevel                 = "logger_default_logging_level"
	ClientErrorLevel             = "logger_client_error_level"
	ServerErrorLevel             = "logger_server_error_level"
	LogWithRequestId             = "logger_with_request_id"
	LogWithRequestIdAutoComplete = "logger_with_request_auto_complete"
	LogWithRequestBody           = "logger_with_request_body"
	LogWithRequestHeader         = "logger_with_request_header"
	LogWithResponseBody          = "logger_with_response_body"
	BodyRequestMaxSize           = "logger_body_request_max_size"
	BodyResponseMaxSize          = "logger_body_response_max_size"
	BodyMaxSize                  = 64 * 1024 // 64KB
)

var defaultValues = map[string]interface{}{
	DefaultLevel:     loggerFif.InfoLevel,
	ClientErrorLevel: loggerFif.WarnLevel,
	ServerErrorLevel: loggerFif.ErrorLevel,
}

var configEntry = []configFif.ConfigEntry{
	{
		VariableName: DefaultLevel,
		DefaultValue: defaultValues[DefaultLevel],
		Description:  "Nivel por default del logger",
	},
	{
		VariableName: ClientErrorLevel,
		DefaultValue: defaultValues[ClientErrorLevel],
		Description:  "Nivel por default para los errores 4xx",
	},
	{
		VariableName: ServerErrorLevel,
		DefaultValue: defaultValues[ServerErrorLevel],
		Description:  "Nivel por default para los errores 5xx",
	},
	{
		VariableName: LogWithRequestId,
		DefaultValue: true,
		Description:  "Indica si agrega al el request id al log",
	},
	{
		VariableName: LogWithRequestIdAutoComplete,
		DefaultValue: true,
		Description:  "Indica si debe generar un request id en el caso de que este sea vacio",
	},
	{
		VariableName: LogWithRequestBody,
		DefaultValue: true,
		Description:  "Indica si debe loguear el body del request",
	},
	{
		VariableName: LogWithRequestHeader,
		DefaultValue: true,
		Description:  "Indica si debe loguear los header del request",
	},
	{
		VariableName: LogWithResponseBody,
		DefaultValue: true,
		Description:  "Indica si debe loguear el body del response",
	},
	{
		VariableName: BodyRequestMaxSize,
		DefaultValue: BodyMaxSize,
		Description:  "Tamaño maximo del request body",
	},
	{
		VariableName: BodyResponseMaxSize,
		DefaultValue: BodyMaxSize,
		Description:  "Tamaño maximo del response body",
	},
}

func GetConfigForLogger() Config {
	values := configFif.LoadConfig(configEntry)
	defaultLevel := getLevelFromConfigOrDefault(values, DefaultLevel)
	clientLevel := getLevelFromConfigOrDefault(values, ClientErrorLevel)
	serverLevel := getLevelFromConfigOrDefault(values, ServerErrorLevel)
	return Config{
		DefaultLevel:              defaultLevel,
		ClientErrorLevel:          clientLevel,
		ServerErrorLevel:          serverLevel,
		WithRequestID:             values[LogWithRequestId].(bool),
		WithRequestIDAutoComplete: values[LogWithRequestIdAutoComplete].(bool),
		WithRequestBody:           values[LogWithRequestBody].(bool),
		WithRequestHeader:         values[LogWithRequestHeader].(bool),
		WithResponseBody:          values[LogWithResponseBody].(bool),
		BodyResponseMaxSize:       values[BodyResponseMaxSize].(int),
		BodyRequestMaxSize:        values[BodyRequestMaxSize].(int),
	}
}

func getLevelFromConfigOrDefault(cfg map[string]interface{}, key string) loggerFif.LogLevel {
	param, isOk := cfg[key].(string)
	if !isOk {
		return defaultValues[key].(loggerFif.LogLevel)
	}
	sLevel, err := loggerFif.StringToLogLevel(param)
	if err != nil {
		return defaultValues[key].(loggerFif.LogLevel)
	}
	return sLevel
}
