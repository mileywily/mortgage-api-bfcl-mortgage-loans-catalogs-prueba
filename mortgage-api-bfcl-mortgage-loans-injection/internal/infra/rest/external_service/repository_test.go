package external_service_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	httpFifMocks "github.com/falabella-regulado/go-lib-http-fif/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-hello-world/internal/core/domain"
	"api-hello-world/internal/infra/rest/external_service"
)

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, errors.New("read failed") }

func httpResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestExternalServiceRepository_Do(t *testing.T) {
	mapperError := errors.New("mapper failed")
	transportError := errors.New("connection refused")

	tests := []struct {
		name           string
		requestMapper  external_service.RequestMapper
		responseMapper external_service.ResponseMapper
		httpResponse   *http.Response
		httpError      error
		expected       domain.EntityOut
		expectedError  interface{}
		expectedInMsg  string
	}{
		{
			name:         "success",
			httpResponse: httpResponse(http.StatusOK, `{"field":"value"}`),
			expected:     domain.EntityOut{Field2: "value"},
		},
		{
			name: "request mapper error",
			requestMapper: func(domain.Entity) (external_service.ExternalServiceRequest, error) {
				return external_service.ExternalServiceRequest{}, mapperError
			},
			expectedError: &domain.ServiceError{},
			expectedInMsg: "mapping request",
		},
		{
			name:          "transport error",
			httpError:     transportError,
			expectedError: &domain.ServiceError{},
			expectedInMsg: "request to",
		},
		{
			name:          "read body error",
			httpResponse:  &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(errorReader{})},
			expectedError: &domain.ServiceError{},
			expectedInMsg: "can't read",
		},
		{
			name:          "5xx maps to service unavailable error",
			httpResponse:  httpResponse(http.StatusServiceUnavailable, `{"error":"down"}`),
			expectedError: &domain.ServiceUnavailableError{},
			expectedInMsg: "returned status 503",
		},
		{
			name:          "4xx maps to service client error",
			httpResponse:  httpResponse(http.StatusNotFound, `{"error":"not found"}`),
			expectedError: &domain.ServiceClientError{},
			expectedInMsg: "status 404",
		},
		{
			name:          "unexpected 3xx status",
			httpResponse:  httpResponse(http.StatusMovedPermanently, ``),
			expectedError: &domain.ServiceError{},
			expectedInMsg: "unexpected status 301",
		},
		{
			name:          "invalid JSON",
			httpResponse:  httpResponse(http.StatusOK, `not-json`),
			expectedError: &domain.ServiceError{},
			expectedInMsg: "can't decode",
		},
		{
			name: "response mapper error",
			responseMapper: func(external_service.ExternalServiceResponse) (domain.EntityOut, error) {
				return domain.EntityOut{}, mapperError
			},
			httpResponse:  httpResponse(http.StatusOK, `{"field":"value"}`),
			expectedError: mapperError,
			expectedInMsg: "mapper failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := new(httpFifMocks.RestClientMock)
			clientMock.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tt.httpResponse, tt.httpError).Maybe()

			requestMapper := tt.requestMapper
			if requestMapper == nil {
				requestMapper = external_service.NewRequestMapper()
			}
			responseMapper := tt.responseMapper
			if responseMapper == nil {
				responseMapper = external_service.NewResponseMapper()
			}

			repository := external_service.NewExternalServiceRepository(clientMock, requestMapper, responseMapper)
			res, err := repository.Do(context.Background(), domain.Entity{Field1: "value"})

			assert.Equal(t, tt.expected, res)
			if tt.expectedError == nil {
				assert.NoError(t, err)
				return
			}
			if !assert.Error(t, err) {
				return
			}
			assert.Contains(t, err.Error(), tt.expectedInMsg)
			switch expected := tt.expectedError.(type) {
			case *domain.ServiceError:
				var serviceErr *domain.ServiceError
				assert.ErrorAs(t, err, &serviceErr)
				_ = expected
			case *domain.ServiceUnavailableError:
				var unavailableErr *domain.ServiceUnavailableError
				assert.ErrorAs(t, err, &unavailableErr)
				assert.Equal(t, `{"error":"down"}`, unavailableErr.Body)
			case *domain.ServiceClientError:
				var clientErr *domain.ServiceClientError
				assert.ErrorAs(t, err, &clientErr)
				assert.Equal(t, http.StatusNotFound, clientErr.StatusCode)
			default:
				assert.Equal(t, tt.expectedError, err)
			}
		})
	}
}
