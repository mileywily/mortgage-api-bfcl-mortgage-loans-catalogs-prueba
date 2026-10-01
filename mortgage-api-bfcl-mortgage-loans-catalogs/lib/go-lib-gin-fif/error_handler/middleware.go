package error_handler

import (
	"github.com/gin-gonic/gin"
	errorFif "github.com/falabella-regulado/go-lib-error-fif"
	"net/http"
)

func NewErrorHandlerMiddleware(h errorFif.ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			c.AbortWithStatusJSON(h(c.Errors[0].Err))
		}
	}
}

func NewDefaultErrorHandlerMiddleware() gin.HandlerFunc {
	return NewErrorHandlerMiddleware(defaultErrorHandler())
}

func defaultErrorHandler() errorFif.ErrorHandler {
	return func(err error) (int, interface{}) {
		return http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": err.Error()}
	}
}

func NewFifErrorHandlerMiddleware(handlers ...errorFif.ErrorHandlerMiddleware) gin.HandlerFunc {
	h := errorFif.DefaultErrorHandler()
	h = errorFif.FifErrorHandler(h)
	for _, mdw := range handlers {
		h = mdw(h)
	}
	return NewErrorHandlerMiddleware(h)
}
