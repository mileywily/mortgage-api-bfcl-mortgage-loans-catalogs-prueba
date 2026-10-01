package external_service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"api-hello-world/internal/core/domain"
	"api-hello-world/internal/infra/rest/external_service"
)

func TestNewRequestMapper_Map(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.Entity
		expected external_service.ExternalServiceRequest
	}{
		{
			name:     "all fields",
			input:    domain.Entity{Field1: "value"},
			expected: external_service.ExternalServiceRequest{FieldExternalService: "value"},
		},
		{
			name:     "empty data",
			input:    domain.Entity{},
			expected: external_service.ExternalServiceRequest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := external_service.NewRequestMapper()(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestNewResponseMapper_Map(t *testing.T) {
	tests := []struct {
		name     string
		input    external_service.ExternalServiceResponse
		expected domain.EntityOut
	}{
		{
			name:     "all fields",
			input:    external_service.ExternalServiceResponse{Field: "value"},
			expected: domain.EntityOut{Field2: "value"},
		},
		{
			name:     "empty data",
			input:    external_service.ExternalServiceResponse{},
			expected: domain.EntityOut{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := external_service.NewResponseMapper()(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}
