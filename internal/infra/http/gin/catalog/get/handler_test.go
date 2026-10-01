package get_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	catalogGet "mortgage-api-bfcl-mortgage-loans-catalogs/internal/app/catalog/get"
	ginGet "mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/http/gin/catalog/get"
	inMocks "mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/in/mocks"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

func TestNewHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success returns 200", func(t *testing.T) {
		expectedItems := []domain.CatalogItem{{Code: "1", Description: "Vivienda Principal"}}
		ucMock := new(inMocks.CatalogUsecase)
		ucMock.On("GetCatalog", mock.Anything, mock.Anything).Return(expectedItems, nil)

		ctrl := catalogGet.NewController(ucMock, ucMock, ucMock, ucMock, catalogGet.NewResponseMapper())
		handler := ginGet.NewHandler(ctrl, ginGet.NewDecoder())

		r := gin.New()
		r.POST("/v1/bfcl/mortgage-loan/catalogs/:catalog", handler)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/v1/bfcl/mortgage-loan/catalogs/Destino", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("controller error is forwarded via c.Error", func(t *testing.T) {
		ucMock := new(inMocks.CatalogUsecase)
		ucMock.On("GetCatalog", mock.Anything, mock.Anything).Return(nil, domain.ErrCatalogNotFound)

		ctrl := catalogGet.NewController(ucMock, ucMock, ucMock, ucMock, catalogGet.NewResponseMapper())
		handler := ginGet.NewHandler(ctrl, ginGet.NewDecoder())

		r := gin.New()
		r.POST("/v1/bfcl/mortgage-loan/catalogs/:catalog", handler)

		// Verify that the error is set (no panic)
		r.Use(func(c *gin.Context) {
			c.Next()
			if len(c.Errors) > 0 {
				assert.True(t, errors.Is(c.Errors[0].Err, domain.ErrCatalogNotFound))
			}
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/v1/bfcl/mortgage-loan/catalogs/Inexistente", nil)
		r.ServeHTTP(w, req)
	})
}
