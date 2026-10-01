package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

type CatalogUsecase struct {
	mock.Mock
}

func (m *CatalogUsecase) GetCatalog(ctx context.Context, req domain.GetCatalogRequest) (interface{}, error) {
	ret := m.Called(ctx, req)
	return ret.Get(0), ret.Error(1)
}
