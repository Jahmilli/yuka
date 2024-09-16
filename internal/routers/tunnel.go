package routers

import (
	"yuka/internal/handlers"

	"github.com/gin-gonic/gin"
)

func tunnelRequest(handler handlers.TunnelHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.TunnelRequest(c)
	}
}
