package create

import (
	"context"

	"api-hello-world/internal/core/ports/in"
)

type Controller interface {
	Create(ctx context.Context, req RequestDto) (ResponseDto, error)
}

type controller struct {
	useCase        in.ExampleUsecase
	requestMapper  RequestMapper
	responseMapper ResponseMapper
}

func NewController(useCase in.ExampleUsecase, requestMapper RequestMapper, responseMapper ResponseMapper) Controller {
	return &controller{
		useCase:        useCase,
		requestMapper:  requestMapper,
		responseMapper: responseMapper,
	}
}

func (c *controller) Create(ctx context.Context, req RequestDto) (ResponseDto, error) {
	entity, err := c.requestMapper(req)
	if err != nil {
		return ResponseDto{}, err
	}

	out, err := c.useCase.Execute(ctx, entity)
	if err != nil {
		return ResponseDto{}, err
	}

	return c.responseMapper(out)
}
