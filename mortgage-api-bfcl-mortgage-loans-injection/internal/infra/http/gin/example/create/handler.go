package create

import (
	"net/http"

	"github.com/gin-gonic/gin"

	createApp "api-hello-world/internal/app/example/create"
)

type Decoder = func(c *gin.Context) (createApp.RequestDto, error)

func NewDecoder() Decoder {
	return func(c *gin.Context) (createApp.RequestDto, error) {
		var req createApp.RequestDto

		if err := c.ShouldBindHeader(&req.Header); err != nil {
			return createApp.RequestDto{}, err
		}

		if err := c.ShouldBindUri(&req.Uri); err != nil {
			return createApp.RequestDto{}, err
		}

		if err := c.ShouldBindJSON(&req.Body); err != nil {
			return createApp.RequestDto{}, err
		}

		return req, nil
	}
}

func NewHandler(ctrl createApp.Controller, decoder Decoder) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := decoder(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res, err := ctrl.Create(c.Request.Context(), req)
		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, res)
	}
}
