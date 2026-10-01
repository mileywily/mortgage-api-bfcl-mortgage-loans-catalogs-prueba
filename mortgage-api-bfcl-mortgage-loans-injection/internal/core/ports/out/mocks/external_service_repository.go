package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"api-hello-world/internal/core/domain"
)

type ExternalServiceRepository struct {
	mock.Mock
}

func (m *ExternalServiceRepository) Do(ctx context.Context, entity domain.Entity) (domain.EntityOut, error) {
	ret := m.Called(ctx, entity)
	return ret.Get(0).(domain.EntityOut), ret.Error(1)
}
