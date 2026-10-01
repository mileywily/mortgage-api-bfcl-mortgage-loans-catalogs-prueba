package catalog_service

import (
	"context"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/out"
)

// dummyRepository implementa CatalogRepository retornando datos estáticos
// que replican exactamente el comportamiento del legado Java.
type dummyRepository struct{}

// NewDummyRepository devuelve un repositorio mock para pruebas locales / entorno dummy.
func NewDummyRepository() out.CatalogRepository {
	return &dummyRepository{}
}

func (r *dummyRepository) GetCatalog(_ context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	// Seguros
	if req.CatalogName == "SegurosIncendio" || req.CatalogName == "SegurosDesgravamen" || req.CatalogName == "SegurosCesantia" {
		return []domain.InsuranceCatalogItem{
			{
				InsuranceIdentifier:  "1100438-01-2023-000",
				Description:          "INCENDIO - Everest compañía de seguros generales Chile(0.2255300)",
				Policy:               "100438-01-2023-000",
				CompanyName:          "Everest compañía de seguros generales Chile",
				CompanyCode:          39,
				InsuranceTypeCode:    0,
				PolicyCorrelative:    28,
				Rate:                 0.22553,
				Factor:               0.0002255,
				IndividualPolicyFlag: 0,
				PerQuotaValueFlag:    0,
				ExternalPolicyFlag:   0,
			},
		}, nil
	}

	// Tipos de documento
	if req.CatalogName == "TiposDocumentos" {
		return []domain.TipoDocumentoCatalogItem{
			{Code: 1, Description: "Escritura de Propiedad", GroupID: 2},
		}, nil
	}

	// Sin resultados / error de red simulado
	if req.CatalogName == "Vacio" || req.CatalogName == "ErrorRed" {
		return noResults(), nil
	}

	// Comunas
	if req.CatalogName == "Comunas" {
		return []domain.ComunaCatalogItem{
			{Code: "01101", Description: "Iquique", RegionID: 1},
		}, nil
	}

	// Catálogo estándar genérico
	return []domain.CatalogItem{
		{Code: "1", Description: "Vivienda Principal"},
	}, nil
}
