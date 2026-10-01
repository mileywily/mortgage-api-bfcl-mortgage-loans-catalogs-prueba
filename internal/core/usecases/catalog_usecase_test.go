package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	outMocks "mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/out/mocks"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/usecases"
)

func TestCatalogUsecase_GetCatalog(t *testing.T) {
	req := domain.GetCatalogRequest{CatalogName: "Destino", Channel: "WEB"}

	expectedResult := []domain.CatalogItem{
		{Code: "1", Description: "Vivienda Principal"},
	}

	tests := []struct {
		name            string
		input           domain.GetCatalogRequest
		repositoryOut   interface{}
		repositoryError error
		expected        interface{}
		expectedError   error
	}{
		{
			name:          "success",
			input:         req,
			repositoryOut: expectedResult,
			expected:      expectedResult,
		},
		{
			name:            "repository error is propagated",
			input:           req,
			repositoryError: domain.ErrCatalogNotFound,
			expectedError:   domain.ErrCatalogNotFound,
		},
		{
			name:          "empty result",
			input:         domain.GetCatalogRequest{CatalogName: "Vacio"},
			repositoryOut: []domain.NoResultsCatalogItem{{CodRespuesta: 3, Mensaje: "Sin resultados.", Excepcion: "Ninguna"}},
			expected:      []domain.NoResultsCatalogItem{{CodRespuesta: 3, Mensaje: "Sin resultados.", Excepcion: "Ninguna"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(outMocks.CatalogRepository)
			repoMock.On("GetCatalog", mock.Anything, tt.input).Return(tt.repositoryOut, tt.repositoryError)

			uc := usecases.NewCatalogUsecase(repoMock)
			res, err := uc.GetCatalog(context.Background(), tt.input)

			assert.Equal(t, tt.expected, res)
			assert.Equal(t, tt.expectedError, err)
			repoMock.AssertExpectations(t)
		})
	}
}
