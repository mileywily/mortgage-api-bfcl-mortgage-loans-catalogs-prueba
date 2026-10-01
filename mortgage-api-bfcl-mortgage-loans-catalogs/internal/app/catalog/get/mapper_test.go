package get_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	app "mortgage-api-bfcl-mortgage-loans-catalogs/internal/app/catalog/get"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

func TestRequestMapper(t *testing.T) {
	mapper := app.NewRequestMapper()
	got, err := mapper(app.RequestDto{
		Header: app.HeaderDto{Channel: "WEB", Commerce: "FALABELLA", TransactionID: "123", Backend: "JAVA"},
		Uri:    app.UriDto{CatalogName: "Destino"},
	})
	require.NoError(t, err)
	assert.Equal(t, domain.GetCatalogRequest{
		CatalogName: "Destino", Channel: "WEB", Commerce: "FALABELLA", TransactionID: "123", Backend: domain.CatalogBackendJava,
	}, got)
}

func TestResponseMapperPreservesLegacyJSON(t *testing.T) {
	mapper := app.NewResponseMapper()
	got, err := mapper([]domain.InsuranceCatalogItem{{InsuranceIdentifier: "POL-1", Rate: 0.1, Factor: 0.0002255}})
	require.NoError(t, err)
	body, err := json.Marshal(got)
	require.NoError(t, err)
	assert.JSONEq(t, `[{"IdentificadorSeguro":"POL-1","Descripcion":"","Poliza":"","NombreCompania":"","CodigoCompania":0,"CodigoTipoSeguro":0,"CorrelativoPoliza":0,"Tasa":0.1,"Factor":2.255E-4,"IndicadorPolizaIndividual":0,"PorValorCuota":0,"IndicadorPolizaExterna":0}]`, string(body))
}
