package repository_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	catalog_repository "mortgage-api-bfcl-mortgage-loans-catalogs/internal/infra/rest/catalog_repository"

	"github.com/stretchr/testify/assert"
)

func TestRestRepository_GetCatalog_Insurance_Success(t *testing.T) {
	// Arrange: Create a mock upstream server
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/catalogo_detail/SegurosIncendio", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Verificar Basic Auth (user/pass)
		user, pass, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "test_user", user)
		assert.Equal(t, "test_pass", pass)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"IdentificadorSeguro":"PARITY-000","Descripcion":"SEGURO PARITY","Poliza":"P-000","NombreCompania":"PARITY CIA","CodigoCompania":99,"CodigoTipoSeguro":1,"CorrelativoPoliza":1,"Tasa":0.1,"Factor":2.255E-4,"IndicadorPolizaIndividual":0,"PorValorCuota":0,"IndicadorPolizaExterna":0}]`))
	}))
	defer mockUpstream.Close()

	repo := catalog_repository.NewRestRepository(mockUpstream.URL, "/api/catalogo_detail/", true, "test_user", "test_pass")

	req := domain.GetCatalogRequest{CatalogName: "SegurosIncendio"}

	// Act
	result, err := repo.GetCatalog(context.Background(), req)

	// Assert
	assert.NoError(t, err)

	items, ok := result.([]domain.InsuranceCatalogItem)
	assert.True(t, ok, "Expected result to be []domain.InsuranceCatalogItem")
	assert.Len(t, items, 1)
	assert.Equal(t, 0.0002255, items[0].Factor) // 2.255E-4
}

func TestRestRepository_GetCatalog_401(t *testing.T) {
	// Arrange
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer mockUpstream.Close()

	repo := catalog_repository.NewRestRepository(mockUpstream.URL, "/api/catalogo_detail/", true, "u", "p")

	// Act
	result, err := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "Destino"})

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AuthenticationError")
	assert.Nil(t, result)
}

func TestRestRepository_GetCatalog_402_NotFound(t *testing.T) {
	// Arrange
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockUpstream.Close()

	repo := catalog_repository.NewRestRepository(mockUpstream.URL, "/api/catalogo_detail/", true, "u", "p")

	// Act
	result, err := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "Destino"})

	// Assert
	assert.ErrorIs(t, err, domain.ErrCatalogNotFound)
	assert.Nil(t, result)
}

func TestRestRepository_GetCatalog_SinResultados(t *testing.T) {
	// Arrange
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"codRespuesta":3,"Mensaje":"Sin resultados.","Excepcion":"Ninguna"}]`))
	}))
	defer mockUpstream.Close()

	repo := catalog_repository.NewRestRepository(mockUpstream.URL, "/api/catalogo_detail/", true, "u", "p")

	// Act
	result, err := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{CatalogName: "Destino"})

	// Assert
	assert.NoError(t, err)

	items, ok := result.([]domain.NoResultsCatalogItem)
	assert.True(t, ok)
	assert.Len(t, items, 1)
	assert.Equal(t, 3, items[0].CodRespuesta)
}

func TestRestRepository_JavaPostForwardsLegacyHeaders(t *testing.T) {
	mockUpstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/bfcl/mortgage-loan/catalogs/Destino", r.URL.Path)
		assert.Equal(t, "WEB", r.Header.Get("X-Channel"))
		assert.Equal(t, "FALABELLA", r.Header.Get("X-Commerce"))
		assert.Equal(t, "123", r.Header.Get("X-Transaction-ID"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"codigo_adm":1,"descripcion":"Vivienda Principal"}]`))
	}))
	defer mockUpstream.Close()

	repo := catalog_repository.NewRestRepository(mockUpstream.URL, "/v1/bfcl/mortgage-loan/catalogs/", false, "", "")
	result, err := repo.GetCatalog(context.Background(), domain.GetCatalogRequest{
		CatalogName: "Destino", Channel: "WEB", Commerce: "FALABELLA", TransactionID: "123",
	})
	if !assert.NoError(t, err) {
		return
	}
	items, ok := result.([]domain.CatalogItem)
	if !assert.True(t, ok) || !assert.Len(t, items, 1) {
		return
	}
	assert.Equal(t, "1", items[0].Code)
}
