package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-hello-world/internal/core/domain"
	outMocks "api-hello-world/internal/core/ports/out/mocks"
	"api-hello-world/internal/core/usecases"
)

func TestExampleUsecase_Execute(t *testing.T) {
	tests := []struct {
		name            string
		input           domain.Entity
		repositoryOut   domain.EntityOut
		repositoryError error
		expected        domain.EntityOut
		expectedError   error
	}{
		{
			name:          "success",
			input:         domain.Entity{Field1: "value"},
			repositoryOut: domain.EntityOut{Field2: "result"},
			expected:      domain.EntityOut{Field2: "result"},
		},
		{
			name:            "repository error is propagated",
			input:           domain.Entity{Field1: "value"},
			repositoryError: &domain.ServiceError{Message: "request failed"},
			expectedError:   &domain.ServiceError{Message: "request failed"},
		},
		{
			name:          "empty data",
			input:         domain.Entity{},
			repositoryOut: domain.EntityOut{},
			expected:      domain.EntityOut{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := new(outMocks.ExternalServiceRepository)
			repositoryMock.On("Do", mock.Anything, tt.input).Return(tt.repositoryOut, tt.repositoryError)

			usecase := usecases.NewExampleUsecase(repositoryMock)
			res, err := usecase.Execute(context.Background(), tt.input)

			assert.Equal(t, tt.expected, res)
			assert.Equal(t, tt.expectedError, err)
			repositoryMock.AssertExpectations(t)
		})
	}
}
