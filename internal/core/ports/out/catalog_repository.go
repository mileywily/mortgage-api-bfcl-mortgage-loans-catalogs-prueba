package out

import (
	"context"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

// CatalogRepository es el puerto de salida (outbound) para obtener catálogos desde un sistema externo.
type CatalogRepository interface {
	GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error)
}
