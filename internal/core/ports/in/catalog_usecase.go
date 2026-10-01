package in

import (
	"context"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

// CatalogUsecase es el puerto de entrada (inbound) para obtener catálogos hipotecarios.
type CatalogUsecase interface {
	GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error)
}
