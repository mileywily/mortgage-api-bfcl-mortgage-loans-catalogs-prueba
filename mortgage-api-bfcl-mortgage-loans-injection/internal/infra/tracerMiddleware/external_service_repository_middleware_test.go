package tracerMiddleware_test

import (
	"context"
	"testing"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/mocktracer"
	loggerFifMocks "github.com/falabella-regulado/go-lib-logger-fif/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"api-hello-world/internal/core/domain"
	outMocks "api-hello-world/internal/core/ports/out/mocks"
	"api-hello-world/internal/infra/tracerMiddleware"
)

func TestExternalServiceRepositoryMiddleware_Do(t *testing.T) {
	tests := []struct {
		name            string
		repositoryOut   domain.EntityOut
		repositoryError error
		expectLogError  bool
	}{
		{
			name:          "success finishes span as passed",
			repositoryOut: domain.EntityOut{Field2: "result"},
		},
		{
			name:            "error finishes span as failed and logs",
			repositoryError: &domain.ServiceError{Message: "request failed"},
			expectLogError:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := mocktracer.Start()
			defer mt.Stop()

			repositoryMock := new(outMocks.ExternalServiceRepository)
			repositoryMock.On("Do", mock.Anything, mock.Anything).Return(tt.repositoryOut, tt.repositoryError)
			loggerMock := &loggerFifMocks.Logger{}
			loggerMock.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

			decorated := tracerMiddleware.NewExternalServiceRepositoryMiddleware(repositoryMock, loggerMock)
			res, err := decorated.Do(context.Background(), domain.Entity{Field1: "value"})

			assert.Equal(t, tt.repositoryOut, res)
			assert.Equal(t, tt.repositoryError, err)
			assert.NotEmpty(t, mt.FinishedSpans())
			if tt.expectLogError {
				loggerMock.AssertCalled(t, "Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			} else {
				loggerMock.AssertNotCalled(t, "Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}
