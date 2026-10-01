package catalog

import (
	"errors"

	errorFIF "github.com/falabella-regulado/go-lib-error-fif"

	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

// NewCatalogErrorHandler registra el manejo de errores específicos del dominio de catálogos.
func NewCatalogErrorHandler(next errorFIF.ErrorHandler) errorFIF.ErrorHandler {
	return func(err error) (int, interface{}) {
		// Error de autenticación con el upstream (401)
		var authErr *domain.AuthenticationError
		if errors.As(err, &authErr) {
			d := "Given token not valid for any token type"
			code := "token_not_valid"
			type tokenMessage struct {
				TokenClass string `json:"token_class"`
				TokenType  string `json:"token_type"`
				Message    string `json:"message"`
			}
			return 401, map[string]interface{}{
				"detail":        &d,
				"code":          &code,
				"messages":      []interface{}{tokenMessage{TokenClass: "AccessToken", TokenType: "access", Message: "Token is invalid"}},
				"errors_detail": nil,
			}
		}

		// Catálogo no encontrado (402 — paridad con Java legado)
		if errors.Is(err, domain.ErrCatalogNotFound) {
			ed := "Catálogo no encontrado"
			return 402, map[string]interface{}{
				"detail":        nil,
				"code":          nil,
				"messages":      nil,
				"errors_detail": &ed,
			}
		}

		// Errores de infraestructura genéricos (400 — paridad con Java legado CatalogException)
		var svcErr *domain.ServiceError
		var unavailErr *domain.ServiceUnavailableError
		var clientErr *domain.ServiceClientError
		if errors.As(err, &svcErr) || errors.As(err, &unavailErr) || errors.As(err, &clientErr) {
			ed := "Error en la obtención de catálogos"
			return 400, map[string]interface{}{
				"detail":        nil,
				"code":          nil,
				"messages":      nil,
				"errors_detail": &ed,
			}
		}

		return next(err)
	}
}
