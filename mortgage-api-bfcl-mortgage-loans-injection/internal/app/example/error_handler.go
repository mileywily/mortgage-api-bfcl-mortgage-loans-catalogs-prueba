package example

import (
	errorFIF "github.com/falabella-regulado/go-lib-error-fif"

	"api-hello-world/internal/core/domain"
)

func NewExampleErrorHandler(next errorFIF.ErrorHandler) errorFIF.ErrorHandler {
	return func(err error) (int, interface{}) {
		switch err.(type) {
		case *domain.ServiceUnavailableError:
			return errorFIF.CreateBadGateway(err.Error())
		case *domain.ServiceClientError:
			return errorFIF.CreateInternalServerError(err.Error())
		case *domain.ServiceError:
			return errorFIF.CreateInternalServerError(err.Error())
		default:
			return next(err)
		}
	}
}
