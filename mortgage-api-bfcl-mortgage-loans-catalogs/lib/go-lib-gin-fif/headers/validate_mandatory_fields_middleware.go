package headers

import (
	"github.com/gin-gonic/gin"
)

type VarCreator func() interface{}

func ValidateMandatoryFieldsMiddleware(f VarCreator) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := f()
		if err := c.ShouldBindHeader(h); err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		c.Next()
	}
}

func ValidateMandatoryFieldsMiddlewareWithDefaultHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		var h MandatoryHeaders
		if err := c.ShouldBindHeader(&h); err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		c.Next()
	}
}
