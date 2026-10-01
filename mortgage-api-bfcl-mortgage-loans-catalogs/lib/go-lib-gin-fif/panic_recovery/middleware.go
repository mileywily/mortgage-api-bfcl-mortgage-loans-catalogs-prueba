package panic_recovery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	errorFif "github.com/falabella-regulado/go-lib-error-fif"
	"net/http"
)

func NewPanicRecoveryMiddleware(f errorFif.PanicRecoveryHandler) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.AbortWithStatusJSON(f(recovered))
	})
}

func NewDefaultPanicRecoveryMiddleware() gin.HandlerFunc {
	return NewPanicRecoveryMiddleware(defaultPanicRecoveryFunc())
}

func defaultPanicRecoveryFunc() errorFif.PanicRecoveryHandler {
	return func(recovered interface{}) (int, interface{}) {
		return http.StatusInternalServerError, fmt.Sprintf("%s", recovered)
	}
}

func NewFifErrorPanicRecoveryMiddleware() gin.HandlerFunc {
	return NewPanicRecoveryMiddleware(errorFif.FifPanicRecoveryHandler())
}
