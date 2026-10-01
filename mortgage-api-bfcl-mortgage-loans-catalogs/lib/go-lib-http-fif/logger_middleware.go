package http_fif

import (
	"bytes"
	"encoding/json"
	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"
	"io"
	"net/http"
	"time"
)

type LoggerConfig struct {
	Enabled           bool
	LogRequestHeaders bool
	LogRequestBody    bool
	LogResponse       bool
	LogResponseBody   bool
	ErrorStatusCode   int
}

type LoggerOption func(*LoggerConfig)

func WithLoggerEnabled(value bool) LoggerOption {
	return func(c *LoggerConfig) {
		c.Enabled = value
	}
}

func WithRequestHeaders(value bool) LoggerOption {
	return func(c *LoggerConfig) {
		c.LogRequestHeaders = value
	}
}

func WithRequestBody(value bool) LoggerOption {
	return func(c *LoggerConfig) {
		c.LogRequestBody = value
	}
}

func WithResponseBody(value bool) LoggerOption {
	return func(c *LoggerConfig) {
		c.LogResponseBody = value
	}
}

func WithErrorStatusCode(value int) LoggerOption {
	return func(c *LoggerConfig) {
		c.ErrorStatusCode = value
	}
}

type loggerFifMiddleware struct {
	logger loggerFif.Logger
	next   HTTPClient
	config LoggerConfig
}

func (m *loggerFifMiddleware) Do(req *http.Request) (res *http.Response, err error) {
	if !m.config.Enabled {
		return m.next.Do(req)
	}
	logFields := []interface{}{
		"method", req.Method,
		"url", req.URL.String(),
	}

	if m.config.LogRequestHeaders && req.Header != nil {
		headers := make(map[string]string)
		for k, v := range req.Header {
			headers[k] = v[0]
		}
		logFields = append(logFields, "headers", headers)
	}

	var requestBodyMap map[string]interface{}
	if m.config.LogRequestBody && req.Body != nil {
		// Leer el cuerpo de la petición
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		if err := json.Unmarshal(bodyBytes, &requestBodyMap); err == nil {
			logFields = append(logFields, "request_body", requestBodyMap)
		} else {
			// Si no es un JSON válido, guardar como string
			logFields = append(logFields, "request_body", string(bodyBytes))
		}
	}

	start := time.Now()
	res, err = m.next.Do(req)

	// Calcular duración
	duration := time.Since(start)
	logFields = append(logFields, "elapsed_time", duration.Milliseconds())

	// Registrar información de la respuesta
	if err != nil {
		m.logger.Error("HTTP request failed", append(logFields, "error", err.Error())...)
		return
	}

	logFields = append(logFields, "status_code", res.StatusCode)

	if m.config.LogResponseBody && res.Body != nil {
		// Leer el cuerpo de la respuesta
		bodyBytes, _ := io.ReadAll(res.Body)
		res.Body.Close()
		res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		var responseBodyMap map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &responseBodyMap); err == nil {
			logFields = append(logFields, "response_body", responseBodyMap)
		} else {
			logFields = append(logFields, "response_body", string(bodyBytes))
		}
	}

	if res.StatusCode >= m.config.ErrorStatusCode {
		m.logger.Error("HTTP request error", logFields...)
	} else {
		m.logger.Info("HTTP request completed", logFields...)
	}

	return
}

func makeLoggerFifMiddleware(logger loggerFif.Logger, opts ...LoggerOption) HTTPClientMiddleware {
	return func(next HTTPClient) HTTPClient {
		config := LoggerConfig{
			Enabled:           false,
			LogRequestHeaders: true,
			LogRequestBody:    true,
			LogResponse:       true,
			LogResponseBody:   true,
			ErrorStatusCode:   400,
		}

		for _, opt := range opts {
			opt(&config)
		}

		return &loggerFifMiddleware{
			logger: logger,
			next:   next,
			config: config,
		}
	}
}
