package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	RequestIDHeaderKey = "X-Request-Id"
	requestIDCtx       = "loggers-gin.request-id"
	RequestIDKey       = "id"
)

var (
	HiddenRequestHeaders = map[string]struct{}{
		"authorization": {},
		"cookie":        {},
		"set-cookie":    {},
		"x-auth-token":  {},
		"x-csrf-token":  {},
		"x-xsrf-token":  {},
	}

	HiddenResponseHeaders = map[string]struct{}{
		"set-cookie": {},
	}
)

type Filter func(ctx *gin.Context) bool

type Config struct {
	DefaultLevel     loggerFif.LogLevel
	ClientErrorLevel loggerFif.LogLevel
	ServerErrorLevel loggerFif.LogLevel

	WithRequestID             bool
	WithRequestIDAutoComplete bool
	WithRequestBody           bool
	WithRequestHeader         bool
	WithResponseBody          bool

	BodyRequestMaxSize  int
	BodyResponseMaxSize int
}

func NewFifLoggerMiddleware(logger loggerFif.Logger, filters ...Filter) gin.HandlerFunc {
	config := GetConfigForLogger()
	return func(c *gin.Context) {
		span, _ := tracer.SpanFromContext(c.Request.Context())
		defer span.Finish()
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Extract parameters from the request
		params := map[string]string{}
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}

		// Handle request ID
		requestID := c.GetHeader(RequestIDHeaderKey)

		if config.WithRequestIDAutoComplete && requestID == "" {
			requestID = uuid.New().String()
			c.Header(RequestIDHeaderKey, requestID)
		}

		if config.WithRequestID {
			c.Set(requestIDCtx, requestID)
		}

		method := c.Request.Method

		span.SetTag(ext.HTTPMethod, method)
		span.SetTag(ext.HTTPURL, path)

		baseAttributes := []any{
			"timestamp", start.UTC().Format("2006-01-02T15:04:05.0000"),
			"date", start.UnixNano() / int64(time.Millisecond),
			ext.LogKeyTraceID, span.Context().TraceID(),
			ext.LogKeySpanID, span.Context().SpanID(),
		}

		httpAttributes := map[string]interface{}{

			//"status_code": status,
			"method": method,
			"url":    path,
			"query":  query,
			"params": params,
		}

		if config.WithRequestID {
			baseAttributes = append(baseAttributes, RequestIDKey, requestID)
		}

		if config.WithRequestHeader {
			headers := map[string]string{}
			for k, v := range c.Request.Header {
				if _, found := HiddenRequestHeaders[strings.ToLower(k)]; !found {
					headers[k] = v[0]
				}
			}
			httpAttributes["headers"] = headers
		}

		// Request body
		if config.WithRequestBody {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) //dejo el buffer devuelta al principio
			if len(bodyBytes) > 0 && len(bodyBytes) < config.BodyRequestMaxSize {
				var requestBody map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
					fmt.Println(err.Error())
				}
				httpAttributes["request"] = requestBody
			}
		}

		recorder := &responseRecorder{
			ResponseWriter: c.Writer,
			body:           new(bytes.Buffer),
			status:         http.StatusOK,
		}
		c.Writer = recorder

		//ACA MANDO A EJECUTAR
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		status := recorder.status
		span.SetTag(ext.HTTPCode, status)
		httpAttributes["status_code"] = status
		httpAttributes["elapsed"] = fmt.Sprintf("%dms", latency.Milliseconds())

		// Response body
		if config.WithResponseBody {

			var responseBody map[string]interface{}

			if len(recorder.body.Bytes()) > 0 && len(recorder.body.Bytes()) < config.BodyResponseMaxSize {
				if err := json.Unmarshal(recorder.body.Bytes(), &responseBody); err != nil {
					fmt.Println(err.Error())
				}
				httpAttributes["response"] = responseBody
			}
		}

		for _, filter := range filters {
			if !filter(c) {
				return
			}
		}

		baseAttributes = append(baseAttributes, "contextMap", map[string]interface{}{"http": httpAttributes})

		// Determine log level and message
		level := config.DefaultLevel
		msg := fmt.Sprintf("%s %s %d %dms", method, path, status, latency.Milliseconds())
		if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
			level = config.ClientErrorLevel
		} else if status >= http.StatusInternalServerError {
			level = config.ServerErrorLevel
		}

		// Log the request
		switch level {
		case loggerFif.ErrorLevel:
			logger.Error(msg, baseAttributes...)
		case loggerFif.InfoLevel:
			logger.Info(msg, baseAttributes...)
		case loggerFif.DebugLevel:
			logger.Debug(msg, baseAttributes...)
		case loggerFif.WarnLevel:
			logger.Warn(msg, baseAttributes...)
		}

	}
}
