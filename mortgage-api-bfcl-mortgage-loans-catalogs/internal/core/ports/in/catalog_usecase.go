package in

import (
	"context"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

type CatalogUsecase interface {
	GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error)
}
