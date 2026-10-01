package usecases

import (
	"context"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/in"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/ports/out"
)

type catalogUsecase struct {
	repo out.CatalogRepository
}

// NewCatalogUsecase crea una nueva instancia del caso de uso de catálogos.
func NewCatalogUsecase(repo out.CatalogRepository) in.CatalogUsecase {
	return &catalogUsecase{repo: repo}
}

// GetCatalog ejecuta la lógica de negocio del catálogo.
// Delega al repositorio de infraestructura; cualquier validación de dominio futura se agrega aquí.
func (uc *catalogUsecase) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	return uc.repo.GetCatalog(ctx, req)
}
