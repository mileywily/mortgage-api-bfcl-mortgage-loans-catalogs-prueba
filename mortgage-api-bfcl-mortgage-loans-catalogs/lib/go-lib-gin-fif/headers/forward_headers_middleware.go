package headers

import (
	"context"
	"github.com/gin-gonic/gin"
)

const mandatoryHeaders = "mandatory-headers"

// ForwardHeadersMiddleware captura ciertos headers y los inyecta en el contexto
func ForwardHeadersMiddleware(headerKeys []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capturar los headers requeridos a través del servicio de dominio
		headers := readHeaders(c, headerKeys)
		ctx := c.Request.Context()

		// Inyectar headers en el contexto
		ctx = context.WithValue(ctx, mandatoryHeaders, headers)
		c.Request = c.Request.WithContext(ctx)

		// Continuar con la ejecución del request original
		c.Next()
	}
}

func readHeaders(c *gin.Context, headerKeys []string) map[string]string {
	headers := make(map[string]string)
	for _, key := range headerKeys {
		if value := c.GetHeader(key); value != "" {
			headers[key] = value
		}
	}
	return headers
}
