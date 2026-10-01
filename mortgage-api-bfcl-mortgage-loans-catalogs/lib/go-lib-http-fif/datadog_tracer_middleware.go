package http_fif

import (
	"net/http"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

type datadogTracerMiddleware struct {
	next     HTTPClient
	spanName string
}

func (m *datadogTracerMiddleware) Do(req *http.Request) (res *http.Response, err error) {
	span, ctx := tracer.StartSpanFromContext(req.Context(), m.spanName)
	defer func() {
		span.Finish(tracer.WithError(err))
	}()

	req = req.Clone(ctx)

	if err := tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(req.Header)); err != nil {
		return nil, err
	}

	return m.next.Do(req)
}

func makeDatadogTracerMiddleware(spanName string) HTTPClientMiddleware {
	return func(next HTTPClient) HTTPClient {
		if spanName == "" {
			spanName = "http.client.request"
		}
		return &datadogTracerMiddleware{next: next, spanName: spanName}
	}
}
