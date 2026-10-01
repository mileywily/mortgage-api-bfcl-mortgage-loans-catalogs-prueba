package get

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	app "mortgage-api-bfcl-mortgage-loans-catalogs/internal/app/catalog/get"
	"mortgage-api-bfcl-mortgage-loans-catalogs/internal/core/domain"
)

type Decoder func(*gin.Context) (app.RequestDto, error)

func NewDecoder() Decoder {
	return func(c *gin.Context) (app.RequestDto, error) {
		return app.RequestDto{
			Header: app.HeaderDto{
				Channel:       c.GetHeader("X-Channel"),
				Commerce:      c.GetHeader("X-Commerce"),
				TransactionID: c.GetHeader("X-Transaction-ID"),
				Backend:       c.GetHeader("X-Backend-Env"),
			},
			Uri: app.UriDto{CatalogName: c.Param("catalog")},
		}, nil
	}
}

func NewHandler(controller app.Controller, decoder Decoder) gin.HandlerFunc {
	return func(c *gin.Context) {
		request, err := decoder(c)
		if err != nil {
			writeError(c, err)
			return
		}

		response, err := controller.GetCatalog(c.Request.Context(), request)
		if err != nil {
			writeError(c, err)
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

type errorResponseDTO struct {
	Detail       *string       `json:"detail"`
	Code         *string       `json:"code"`
	Messages     []interface{} `json:"messages"`
	ErrorsDetail *string       `json:"errors_detail"`
}

type tokenMessageDTO struct {
	TokenClass string `json:"token_class"`
	TokenType  string `json:"token_type"`
	Message    string `json:"message"`
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrAuthentication):
		detail := "Given token not valid for any token type"
		code := "token_not_valid"
		c.JSON(http.StatusUnauthorized, errorResponseDTO{
			Detail: &detail,
			Code:   &code,
			Messages: []interface{}{tokenMessageDTO{
				TokenClass: "AccessToken",
				TokenType:  "access",
				Message:    "Token is invalid",
			}},
		})
	case errors.Is(err, domain.ErrCatalogNotFound):
		detail := "Catálogo no encontrado"
		c.JSON(http.StatusPaymentRequired, errorResponseDTO{ErrorsDetail: &detail})
	default:
		detail := "Error en la obtención de catálogos"
		c.JSON(http.StatusBadRequest, errorResponseDTO{ErrorsDetail: &detail})
	}
}
