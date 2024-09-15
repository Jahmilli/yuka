package routers

import (
	"yuka/internal/handlers"

	"github.com/gin-gonic/gin"
)

func initializeStream(handler handlers.StreamingHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.InitialiseStream(c)
	}
}
