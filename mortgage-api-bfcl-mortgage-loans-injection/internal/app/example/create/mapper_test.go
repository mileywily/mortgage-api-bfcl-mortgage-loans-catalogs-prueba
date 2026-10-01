package create_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"api-hello-world/internal/app/example/create"
	"api-hello-world/internal/core/domain"
)

func TestNewRequestMapper_Map(t *testing.T) {
	tests := []struct {
		name     string
		input    create.RequestDto
		expected domain.Entity
	}{
		{
			name: "all fields",
			input: create.RequestDto{
				Header: create.HeaderDto{XField: "headerValue"},
				Uri:    create.UriDto{Something: "anything"},
				Body:   create.BodyRequestDto{Body: create.BodyDto{Field2: "bodyValue"}},
			},
			expected: domain.Entity{Field1: "bodyValue"},
		},
		{
			name:     "empty optional data",
			input:    create.RequestDto{},
			expected: domain.Entity{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := create.NewRequestMapper()(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}

func TestNewResponseMapper_Map(t *testing.T) {
	tests := []struct {
		name     string
		input    domain.EntityOut
		expected create.ResponseDto
	}{
		{
			name:     "all fields",
			input:    domain.EntityOut{Field2: "result"},
			expected: create.ResponseDto{Code: "222", Message: "result"},
		},
		{
			name:     "empty data",
			input:    domain.EntityOut{},
			expected: create.ResponseDto{Code: "222", Message: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := create.NewResponseMapper()(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}
