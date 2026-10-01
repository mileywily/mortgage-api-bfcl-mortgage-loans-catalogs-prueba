package in

import (
	"context"

	"api-hello-world/internal/core/domain"
)

// TODO: cambiar nombre por algo representativo del caso de uso real
type ExampleUsecase interface {
	Execute(ctx context.Context, entity domain.Entity) (domain.EntityOut, error)
}
