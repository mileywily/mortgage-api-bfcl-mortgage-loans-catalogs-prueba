package usecases_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/usecases"
)

type catalogRepositoryStub struct {
	called bool
	result interface{}
}

func (r *catalogRepositoryStub) GetCatalog(_ context.Context, _ domain.GetCatalogRequest) (interface{}, error) {
	r.called = true
	return r.result, nil
}

func TestCatalogUsecaseSelectsBackend(t *testing.T) {
	dummy := &catalogRepositoryStub{result: "dummy"}
	finnflow := &catalogRepositoryStub{result: "finnflow"}
	java := &catalogRepositoryStub{result: "java"}
	usecase := usecases.NewCatalogUsecase(domain.CatalogBackendFinnflow, dummy, finnflow, java)

	tests := []struct {
		name    string
		backend domain.CatalogBackend
		want    string
	}{
		{name: "configured default", want: "finnflow"},
		{name: "dummy override", backend: domain.CatalogBackendDummy, want: "dummy"},
		{name: "java override", backend: domain.CatalogBackendJava, want: "java"},
		{name: "unknown override uses default", backend: "unknown", want: "finnflow"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := domain.GetCatalogRequest{Backend: tt.backend}
			got, err := usecase.GetCatalog(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
