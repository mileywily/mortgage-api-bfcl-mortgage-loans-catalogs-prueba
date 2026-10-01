package get

import (
	"net/http"

	"github.com/gin-gonic/gin"

	catalogGet "mortgage-api-bfcl-mortgage-loans-catalogs/internal/app/catalog/get"
)

// Decoder lee los parámetros de entrada del contexto Gin y los convierte en RequestDto.
type Decoder = func(c *gin.Context) (catalogGet.RequestDto, error)

// NewDecoder crea el decodificador de solicitudes para el endpoint de catálogos.
func NewDecoder() Decoder {
	return func(c *gin.Context) (catalogGet.RequestDto, error) {
		var req catalogGet.RequestDto

		// Leer headers (X-Channel, X-Commerce, X-Transaction-ID, X-Backend-Env)
		if err := c.ShouldBindHeader(&req.Header); err != nil {
			return catalogGet.RequestDto{}, err
		}

		// Leer parámetro de ruta /:catalog
		if err := c.ShouldBindUri(&req.Uri); err != nil {
			return catalogGet.RequestDto{}, err
		}

		return req, nil
	}
}

// NewHandler crea el gin.HandlerFunc que decodifica la petición, delega al controller
// y escribe la respuesta JSON.
func NewHandler(ctrl catalogGet.Controller, decoder Decoder) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := decoder(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res, err := ctrl.GetCatalog(c.Request.Context(), req)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
