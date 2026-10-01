package usecases

import (
	"context"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/in"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/out"
)

type catalogUsecase struct {
	defaultBackend domain.CatalogBackend
	repositories   map[domain.CatalogBackend]out.CatalogRepository
}

func NewCatalogUsecase(
	defaultBackend domain.CatalogBackend,
	dummyRepository out.CatalogRepository,
	finnflowRepository out.CatalogRepository,
	javaRepository out.CatalogRepository,
) in.CatalogUsecase {
	if defaultBackend != domain.CatalogBackendDummy &&
		defaultBackend != domain.CatalogBackendFinnflow &&
		defaultBackend != domain.CatalogBackendJava {
		defaultBackend = domain.CatalogBackendFinnflow
	}

	return &catalogUsecase{
		defaultBackend: defaultBackend,
		repositories: map[domain.CatalogBackend]out.CatalogRepository{
			domain.CatalogBackendDummy:    dummyRepository,
			domain.CatalogBackendFinnflow: finnflowRepository,
			domain.CatalogBackendJava:     javaRepository,
		},
	}
}

func (uc *catalogUsecase) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	backend := req.Backend
	if backend == domain.CatalogBackendDefault {
		backend = uc.defaultBackend
	}

	repository, ok := uc.repositories[backend]
	if !ok {
		repository = uc.repositories[uc.defaultBackend]
	}
	return repository.GetCatalog(ctx, req)
}
