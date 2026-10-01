package get

import (
	"context"
	"strings"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/in"
)

// Controller define el contrato del controlador para obtener catálogos.
type Controller interface {
	GetCatalog(ctx context.Context, req RequestDto) (interface{}, error)
}

type controller struct {
	defaultUsecase in.CatalogUsecase
	dummyUsecase   in.CatalogUsecase
	realUsecase    in.CatalogUsecase
	javaUsecase    in.CatalogUsecase
	responseMapper ResponseMapper
}

// NewController crea un controlador con los 4 backends y el mapper de respuesta.
func NewController(
	defaultUc, dummyUc, realUc, javaUc in.CatalogUsecase,
	responseMapper ResponseMapper,
) Controller {
	return &controller{
		defaultUsecase: defaultUc,
		dummyUsecase:   dummyUc,
		realUsecase:    realUc,
		javaUsecase:    javaUc,
		responseMapper: responseMapper,
	}
}

// GetCatalog selecciona el backend según el header X-Backend-Env, invoca el caso de uso
// y convierte el resultado de dominio a DTO.
func (c *controller) GetCatalog(ctx context.Context, req RequestDto) (interface{}, error) {
	activeUsecase := c.defaultUsecase

	switch strings.ToLower(req.Header.BackendEnv) {
	case "real":
		activeUsecase = c.realUsecase
	case "java":
		activeUsecase = c.javaUsecase
	case "dummy":
		activeUsecase = c.dummyUsecase
	}

	result, err := activeUsecase.GetCatalog(ctx, req.toDomainRequest())
	if err != nil {
		return nil, err
	}

	return c.responseMapper(result)
}
