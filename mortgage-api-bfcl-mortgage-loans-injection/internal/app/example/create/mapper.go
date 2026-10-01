package create

import (
	"api-hello-world/internal/core/domain"
)

type RequestMapper = func(RequestDto) (domain.Entity, error)
type ResponseMapper = func(domain.EntityOut) (ResponseDto, error)

func NewRequestMapper() RequestMapper {
	return func(req RequestDto) (domain.Entity, error) {
		return domain.Entity{
			Field1: req.Body.Body.Field2,
		}, nil
	}
}

func NewResponseMapper() ResponseMapper {
	return func(out domain.EntityOut) (ResponseDto, error) {
		return ResponseDto{
			Code:    "222",
			Message: out.Field2,
		}, nil
	}
}
