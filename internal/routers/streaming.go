package routers

import (
	"yuka/internal/handlers"

	"github.com/gin-gonic/gin"
)

func handleWsConnection(handler handlers.WsHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.HandleWsConnection(c)
	}
}
