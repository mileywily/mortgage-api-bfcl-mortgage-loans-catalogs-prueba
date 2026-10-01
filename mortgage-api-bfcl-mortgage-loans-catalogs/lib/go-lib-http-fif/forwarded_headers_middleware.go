package http_fif

import (
	"net/http"
)

type forwardedHeadersMiddleware struct {
	contextKey string
	next       HTTPClient
}

func (m *forwardedHeadersMiddleware) Do(req *http.Request) (res *http.Response, err error) {
	headers, ok := req.Context().Value(m.contextKey).(map[string]string)
	if !ok {
		return m.next.Do(req)
	}
	for k, v := range headers {
		if old := req.Header.Get(k); old == "" {
			req.Header.Set(k, v)
		}
	}
	return m.next.Do(req)
}

func makeForwardedHeadersMiddleware(contextKey string) HTTPClientMiddleware {
	return func(next HTTPClient) HTTPClient {
		return &forwardedHeadersMiddleware{
			contextKey: contextKey,
			next:       next,
		}
	}
}
