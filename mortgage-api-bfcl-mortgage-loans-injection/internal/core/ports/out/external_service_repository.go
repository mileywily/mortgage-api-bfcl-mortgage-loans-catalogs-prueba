package out

import (
	"context"

	"api-hello-world/internal/core/domain"
)

// TODO: reemplazar por un nombre que represente al servicio externo real
type ExternalServiceRepository interface {
	Do(ctx context.Context, entity domain.Entity) (domain.EntityOut, error)
}
