package get_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/app/catalog/get"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	inMocks "mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/in/mocks"
)

func TestController_GetCatalog(t *testing.T) {
	expectedItems := []domain.CatalogItem{{Code: "1", Description: "Vivienda Principal"}}

	tests := []struct {
		name          string
		backendEnv    string
		usecaseResult interface{}
		usecaseError  error
		expectedError bool
	}{
		{
			name:          "default backend returns catalog",
			backendEnv:    "",
			usecaseResult: expectedItems,
		},
		{
			name:          "dummy backend selected",
			backendEnv:    "dummy",
			usecaseResult: expectedItems,
		},
		{
			name:          "real backend selected",
			backendEnv:    "real",
			usecaseResult: expectedItems,
		},
		{
			name:          "java backend selected",
			backendEnv:    "java",
			usecaseResult: expectedItems,
		},
		{
			name:          "usecase error is propagated",
			backendEnv:    "",
			usecaseError:  domain.ErrCatalogNotFound,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaultMock := new(inMocks.CatalogUsecase)
			dummyMock := new(inMocks.CatalogUsecase)
			realMock := new(inMocks.CatalogUsecase)
			javaMock := new(inMocks.CatalogUsecase)

			req := get.RequestDto{
				Header: get.HeaderDto{BackendEnv: tt.backendEnv},
				Uri:    get.UriDto{Catalog: "Destino"},
			}

			switch tt.backendEnv {
			case "dummy":
				dummyMock.On("GetCatalog", mock.Anything, mock.Anything).Return(tt.usecaseResult, tt.usecaseError)
			case "real":
				realMock.On("GetCatalog", mock.Anything, mock.Anything).Return(tt.usecaseResult, tt.usecaseError)
			case "java":
				javaMock.On("GetCatalog", mock.Anything, mock.Anything).Return(tt.usecaseResult, tt.usecaseError)
			default:
				defaultMock.On("GetCatalog", mock.Anything, mock.Anything).Return(tt.usecaseResult, tt.usecaseError)
			}

			ctrl := get.NewController(defaultMock, dummyMock, realMock, javaMock, get.NewResponseMapper())
			res, err := ctrl.GetCatalog(context.Background(), req)

			if tt.expectedError {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, domain.ErrCatalogNotFound))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}
