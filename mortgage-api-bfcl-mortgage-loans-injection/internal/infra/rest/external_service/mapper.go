package external_service

import (
	"api-hello-world/internal/core/domain"
)

type RequestMapper = func(domain.Entity) (ExternalServiceRequest, error)
type ResponseMapper = func(ExternalServiceResponse) (domain.EntityOut, error)

func NewRequestMapper() RequestMapper {
	return func(entity domain.Entity) (ExternalServiceRequest, error) {
		return ExternalServiceRequest{
			FieldExternalService: entity.Field1,
		}, nil
	}
}

func NewResponseMapper() ResponseMapper {
	return func(res ExternalServiceResponse) (domain.EntityOut, error) {
		return domain.EntityOut{
			Field2: res.Field,
		}, nil
	}
}
