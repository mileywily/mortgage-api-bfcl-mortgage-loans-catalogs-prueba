package get

import (
	"context"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/in"
)

type Controller interface {
	GetCatalog(ctx context.Context, req RequestDto) (interface{}, error)
}

type RequestMapper func(RequestDto) (domain.GetCatalogRequest, error)
type ResponseMapper func(interface{}) (interface{}, error)

type controller struct {
	usecase        in.CatalogUsecase
	requestMapper  RequestMapper
	responseMapper ResponseMapper
}

func NewController(usecase in.CatalogUsecase, requestMapper RequestMapper, responseMapper ResponseMapper) Controller {
	return &controller{
		usecase:        usecase,
		requestMapper:  requestMapper,
		responseMapper: responseMapper,
	}
}

func (c *controller) GetCatalog(ctx context.Context, req RequestDto) (interface{}, error) {
	entity, err := c.requestMapper(req)
	if err != nil {
		return nil, err
	}

	result, err := c.usecase.GetCatalog(ctx, entity)
	if err != nil {
		return nil, err
	}

	return c.responseMapper(result)
}
