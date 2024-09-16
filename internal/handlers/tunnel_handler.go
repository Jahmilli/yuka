package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TunnelRequestResp struct {
	Path string `json:"path"`
}

type TunnelHandler struct {
	Db             *gorm.DB
	Slogger        *zap.SugaredLogger
	ConnectionPool *StreamingConnectionPool
}

func NewTunnelHandler(logger *zap.Logger, db *gorm.DB) TunnelHandler {
	return TunnelHandler{
		Db:      db,
		Slogger: logger.Sugar(),
	}
}

func (s *TunnelHandler) TunnelRequest(c *gin.Context) error {
	param := c.Param("tunnelPath")
	c.JSON(http.StatusOK, &TunnelRequestResp{Path: param})
	return nil
}
