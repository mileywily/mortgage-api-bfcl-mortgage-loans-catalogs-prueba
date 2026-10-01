package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"api-hello-world/internal/app/example/create"
)

type Controller struct {
	mock.Mock
}

func (m *Controller) Create(ctx context.Context, req create.RequestDto) (create.ResponseDto, error) {
	ret := m.Called(ctx, req)
	return ret.Get(0).(create.ResponseDto), ret.Error(1)
}
