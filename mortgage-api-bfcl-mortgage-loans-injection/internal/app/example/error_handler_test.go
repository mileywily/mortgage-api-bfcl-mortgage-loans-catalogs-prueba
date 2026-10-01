package example_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"api-hello-world/internal/app/example"
	"api-hello-world/internal/core/domain"
)

func TestNewExampleErrorHandler_Handle(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectDelegate bool
	}{
		{
			name:           "service unavailable error maps to bad gateway",
			err:            &domain.ServiceUnavailableError{Message: "external service returned status 503"},
			expectedStatus: http.StatusBadGateway,
		},
		{
			name:           "service client error maps to internal server error",
			err:            &domain.ServiceClientError{StatusCode: 404},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "service error maps to internal server error",
			err:            &domain.ServiceError{Message: "can't decode response"},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "unknown error delegates to next handler",
			err:            errors.New("unknown"),
			expectedStatus: http.StatusTeapot,
			expectDelegate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegated := false
			next := func(err error) (int, interface{}) {
				delegated = true
				return http.StatusTeapot, nil
			}

			status, _ := example.NewExampleErrorHandler(next)(tt.err)

			assert.Equal(t, tt.expectedStatus, status)
			assert.Equal(t, tt.expectDelegate, delegated)
		})
	}
}
