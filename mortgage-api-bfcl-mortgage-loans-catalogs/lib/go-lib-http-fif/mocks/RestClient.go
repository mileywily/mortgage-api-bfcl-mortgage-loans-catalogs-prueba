package mocks

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// RestClientMock es un mock de la interfaz RestClient
type RestClientMock struct {
	mock.Mock
}

func (m *RestClientMock) Get(ctx context.Context, path string, headers http.Header, params ...string) (*http.Response, error) {
	args := m.Called(ctx, path, headers, params)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *RestClientMock) Post(ctx context.Context, path string, body interface{}, headers http.Header) (*http.Response, error) {
	args := m.Called(ctx, path, body, headers)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *RestClientMock) Put(ctx context.Context, path string, body interface{}, headers http.Header) (*http.Response, error) {
	args := m.Called(ctx, path, body, headers)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *RestClientMock) Patch(ctx context.Context, path string, body interface{}, headers http.Header) (*http.Response, error) {
	args := m.Called(ctx, path, body, headers)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *RestClientMock) Delete(ctx context.Context, path string, headers http.Header) (*http.Response, error) {
	args := m.Called(ctx, path, headers)
	return args.Get(0).(*http.Response), args.Error(1)
}
