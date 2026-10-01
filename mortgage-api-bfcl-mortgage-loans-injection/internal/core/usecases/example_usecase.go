package usecases

import (
	"context"

	"api-hello-world/internal/core/domain"
	"api-hello-world/internal/core/ports/in"
	"api-hello-world/internal/core/ports/out"
)

type exampleUsecase struct {
	externalServiceRepository out.ExternalServiceRepository
}

func NewExampleUsecase(externalServiceRepository out.ExternalServiceRepository) in.ExampleUsecase {
	return &exampleUsecase{
		externalServiceRepository: externalServiceRepository,
	}
}

func (uc *exampleUsecase) Execute(ctx context.Context, entity domain.Entity) (domain.EntityOut, error) {
	return uc.externalServiceRepository.Do(ctx, entity)
}
