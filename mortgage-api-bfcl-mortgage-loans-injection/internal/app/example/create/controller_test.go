package create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-hello-world/internal/app/example/create"
	"api-hello-world/internal/core/domain"
	inMocks "api-hello-world/internal/core/ports/in/mocks"
)

func TestController_Create(t *testing.T) {
	request := create.RequestDto{Body: create.BodyRequestDto{Body: create.BodyDto{Field2: "value"}}}
	mapperError := errors.New("mapper failed")

	tests := []struct {
		name           string
		requestMapper  create.RequestMapper
		responseMapper create.ResponseMapper
		usecaseOut     domain.EntityOut
		usecaseError   error
		expected       create.ResponseDto
		expectedError  error
	}{
		{
			name:           "success",
			requestMapper:  create.NewRequestMapper(),
			responseMapper: create.NewResponseMapper(),
			usecaseOut:     domain.EntityOut{Field2: "result"},
			expected:       create.ResponseDto{Code: "222", Message: "result"},
		},
		{
			name: "request mapper error is propagated",
			requestMapper: func(create.RequestDto) (domain.Entity, error) {
				return domain.Entity{}, mapperError
			},
			responseMapper: create.NewResponseMapper(),
			expectedError:  mapperError,
		},
		{
			name:           "usecase error is propagated",
			requestMapper:  create.NewRequestMapper(),
			responseMapper: create.NewResponseMapper(),
			usecaseError:   &domain.ServiceError{Message: "request failed"},
			expectedError:  &domain.ServiceError{Message: "request failed"},
		},
		{
			name:          "response mapper error is propagated",
			requestMapper: create.NewRequestMapper(),
			responseMapper: func(domain.EntityOut) (create.ResponseDto, error) {
				return create.ResponseDto{}, mapperError
			},
			usecaseOut:    domain.EntityOut{Field2: "result"},
			expectedError: mapperError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecaseMock := new(inMocks.ExampleUsecase)
			usecaseMock.On("Execute", mock.Anything, mock.Anything).Return(tt.usecaseOut, tt.usecaseError).Maybe()

			ctrl := create.NewController(usecaseMock, tt.requestMapper, tt.responseMapper)
			res, err := ctrl.Create(context.Background(), request)

			assert.Equal(t, tt.expected, res)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
