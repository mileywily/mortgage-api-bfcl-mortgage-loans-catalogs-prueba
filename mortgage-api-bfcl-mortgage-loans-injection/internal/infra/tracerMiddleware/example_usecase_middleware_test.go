package tracerMiddleware_test

import (
	"context"
	"testing"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/mocktracer"
	loggerFifMocks "github.com/falabella-regulado/go-lib-logger-fif/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-hello-world/internal/core/domain"
	inMocks "api-hello-world/internal/core/ports/in/mocks"
	"api-hello-world/internal/infra/tracerMiddleware"
)

func TestExampleUsecaseMiddleware_Execute(t *testing.T) {
	tests := []struct {
		name           string
		usecaseOut     domain.EntityOut
		usecaseError   error
		expectLogError bool
	}{
		{
			name:       "success finishes span as passed",
			usecaseOut: domain.EntityOut{Field2: "result"},
		},
		{
			name:           "error finishes span as failed and logs",
			usecaseError:   &domain.ServiceError{Message: "request failed"},
			expectLogError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := mocktracer.Start()
			defer mt.Stop()

			usecaseMock := new(inMocks.ExampleUsecase)
			usecaseMock.On("Execute", mock.Anything, mock.Anything).Return(tt.usecaseOut, tt.usecaseError)
			loggerMock := &loggerFifMocks.Logger{}
			loggerMock.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

			decorated := tracerMiddleware.NewExampleUsecaseMiddleware(usecaseMock, loggerMock)
			res, err := decorated.Execute(context.Background(), domain.Entity{Field1: "value"})

			assert.Equal(t, tt.usecaseOut, res)
			assert.Equal(t, tt.usecaseError, err)
			assert.NotEmpty(t, mt.FinishedSpans())
			if tt.expectLogError {
				loggerMock.AssertCalled(t, "Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			} else {
				loggerMock.AssertNotCalled(t, "Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}
