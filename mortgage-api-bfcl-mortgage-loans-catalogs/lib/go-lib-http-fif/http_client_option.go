package http_fif

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	loggerFif "github.com/falabella-regulado/go-lib-logger-fif"
)

const (
	mandatoryHeaders = "mandatory-headers"
)

// HttpClientOption is a function type that modifies a httpClient.
type HttpClientOption func(httpClient) httpClient

// DialContextFunc is a type alias for a function that establishes a network connection
type DialContextFunc func(ctx context.Context, network, addr string) (net.Conn, error)

// DialContextFuncMiddleware defines a middleware type that takes a DialContextFunc
type DialContextFuncMiddleware func(DialContextFunc) DialContextFunc

// BaseURL sets the base URL for the client.
// It panics if the provided baseURL is invalid.
func BaseURL(baseURL string) HttpClientOption {
	return func(client httpClient) httpClient {
		var err error
		client.baseURL, err = url.Parse(baseURL)
		if err != nil {
			panic(fmt.Errorf("%s: %w", err, ErrInvalidURLAddress))
		}
		return client
	}
}

// Header adds a header to the client's request.
func Header(key, value string) HttpClientOption {
	return func(client httpClient) httpClient {
		client.headers = client.headers.Clone()
		client.headers.Add(key, value)
		return client
	}
}

// Timeout sets the timeout duration for the client's requests.
func Timeout(timeout time.Duration) HttpClientOption {
	return func(client httpClient) httpClient {
		client.baseClient.Timeout = timeout
		return client
	}
}

func Transport(transport http.RoundTripper) HttpClientOption {
	return func(client httpClient) httpClient {
		client.baseClient.Transport = transport
		return client
	}
}

func Jar(jar http.CookieJar) HttpClientOption {
	return func(client httpClient) httpClient {
		client.baseClient.Jar = jar
		return client
	}
}

func CheckRedirect(redirect func(req *http.Request, via []*http.Request) error) HttpClientOption {
	return func(client httpClient) httpClient {
		client.baseClient.CheckRedirect = redirect
		return client
	}
}

func TLSConfig(tlsConfig *tls.Config) HttpClientOption {
	return func(client httpClient) httpClient {
		client.baseClient.Transport.(*http.Transport).TLSClientConfig = tlsConfig
		return client
	}
}

// Logger sets the logger_fif for the client.
func Logger(logger loggerFif.Logger, opts ...LoggerOption) HttpClientOption {
	return func(client httpClient) httpClient {
		client.logger = makeLoggerFifMiddleware(logger, opts...)
		return client
	}
}

func WithTokenProvider(provider AuthorizationProvider, headerName string) HttpClientOption {
	return func(client httpClient) httpClient {
		if headerName == "" {
			headerName = "Authorization"
		}
		mdw := makeTokenProviderMiddleware(provider, headerName)
		client.requester = mdw(client.requester)
		return client
	}
}

func WithDatadogTracer() HttpClientOption {
	return func(client httpClient) httpClient {
		client.tracer = makeDatadogTracerMiddleware("")
		return client
	}
}

func WithDatadogTracerWithName(spanName string) HttpClientOption {
	return func(client httpClient) httpClient {
		client.tracer = makeDatadogTracerMiddleware(spanName)
		return client
	}
}

func WithMandatoryHeaders() HttpClientOption {
	return func(client httpClient) httpClient {
		mdw := makeForwardedHeadersMiddleware(mandatoryHeaders)
		client.requester = mdw(client.requester)
		return client
	}
}

// TCPTraceDialContext returns an HttpClientOption that wraps the client's DialContext
// function with one or more middlewares. This is typically used to inject additional
// behavior like tracing, logging, or metrics around TCP connections.
func TCPTraceDialContext(mdws ...DialContextFuncMiddleware) HttpClientOption {
	return func(client httpClient) httpClient {
		dialContext := client.baseClient.Transport.(*http.Transport).DialContext
		for _, mdw := range mdws {
			dialContext = mdw(dialContext)
		}
		client.baseClient.Transport.(*http.Transport).DialContext = dialContext
		return client
	}
}
