package get_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	app "mortgage-api-bfcl-mortgage-loans-catalogs/internal/app/catalog/get"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	handler "mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/http/gin/catalog/get"
)

type controllerStub struct {
	response interface{}
	err      error
	request  app.RequestDto
}

func (s *controllerStub) GetCatalog(_ context.Context, request app.RequestDto) (interface{}, error) {
	s.request = request
	return s.response, s.err
}

func setupRouter(controller *controllerStub) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/v1/bfcl/mortgage-loan/catalogs/:catalog", handler.NewHandler(controller, handler.NewDecoder()))
	return router
}

func TestHandlerKeepsLegacyRouteAndAllowsMissingHeaders(t *testing.T) {
	controller := &controllerStub{response: []map[string]string{{"codigo_adm": "1"}}}
	router := setupRouter(controller)
	request := httptest.NewRequest(http.MethodPost, "/v1/bfcl/mortgage-loan/catalogs/Destino", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "Destino", controller.request.Uri.CatalogName)
	assert.Empty(t, controller.request.Header.Channel)
	assert.JSONEq(t, `[{"codigo_adm":"1"}]`, response.Body.String())
}

func TestHandlerMapsLegacyErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code int
	}{
		{name: "unauthorized", err: domain.ErrAuthentication, code: http.StatusUnauthorized},
		{name: "catalog not found", err: domain.ErrCatalogNotFound, code: http.StatusPaymentRequired},
		{name: "other error", err: errors.New("upstream failed"), code: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(&controllerStub{err: tt.err})
			request := httptest.NewRequest(http.MethodPost, "/v1/bfcl/mortgage-loan/catalogs/Destino", nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			assert.Equal(t, tt.code, response.Code)
		})
	}
}
