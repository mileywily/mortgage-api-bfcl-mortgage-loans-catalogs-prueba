package domain_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

func TestServiceError_ErrorAndUnwrap(t *testing.T) {
	cause := errors.New("connection refused")
	tests := []struct {
		name            string
		err             error
		expectedMessage string
		expectedCause   error
	}{
		{
			name:            "message and cause",
			err:             &domain.ServiceError{Message: "request failed", Cause: cause},
			expectedMessage: "request failed",
			expectedCause:   cause,
		},
		{
			name:            "without cause",
			err:             &domain.ServiceError{Message: "can't decode response"},
			expectedMessage: "can't decode response",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var serviceErr *domain.ServiceError
			assert.ErrorAs(t, tt.err, &serviceErr)
			assert.Equal(t, tt.expectedMessage, tt.err.Error())
			assert.Equal(t, tt.expectedCause, errors.Unwrap(tt.err))
		})
	}
}

func TestServiceUnavailableError_Error(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		expectedMessage string
	}{
		{
			name:            "message with body",
			err:             &domain.ServiceUnavailableError{Message: "external service returned status 503", Body: `{"error":"down"}`},
			expectedMessage: "external service returned status 503",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unavailableErr *domain.ServiceUnavailableError
			assert.ErrorAs(t, tt.err, &unavailableErr)
			assert.Equal(t, tt.expectedMessage, tt.err.Error())
		})
	}
}

func TestServiceClientError_Error(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		expectedMessage string
	}{
		{
			name:            "status in message",
			err:             &domain.ServiceClientError{StatusCode: 404, Body: `{"error":"not found"}`},
			expectedMessage: "external service rejected the request with status 404",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var clientErr *domain.ServiceClientError
			assert.ErrorAs(t, tt.err, &clientErr)
			assert.Equal(t, tt.expectedMessage, tt.err.Error())
		})
	}
}
