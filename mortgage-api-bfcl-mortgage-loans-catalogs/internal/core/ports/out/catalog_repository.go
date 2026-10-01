package out

import (
	"context"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

type CatalogRepository interface {
	GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error)
}
