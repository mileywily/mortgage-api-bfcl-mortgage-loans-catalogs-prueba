package http_fif

import (
	"fmt"
	"net/http"
)

type AuthorizationProvider interface {
	GetAuthorization() (string, error)
}

type authorizationMiddlware struct {
	provider   AuthorizationProvider
	headerName string
	next       HTTPClient
}

func (t *authorizationMiddlware) Do(req *http.Request) (*http.Response, error) {
	if req.Header.Get(t.headerName) != "" {
		return t.next.Do(req)
	}
	token, err := t.provider.GetAuthorization()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}
	headerName := t.headerName
	req.Header.Set(headerName, token)

	return t.next.Do(req)
}

func makeTokenProviderMiddleware(provider AuthorizationProvider, headerName string) HTTPClientMiddleware {
	return func(next HTTPClient) HTTPClient {
		return &authorizationMiddlware{
			provider:   provider,
			headerName: headerName,
			next:       next,
		}
	}
}
