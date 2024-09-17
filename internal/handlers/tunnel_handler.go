package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TunnelRequestResp struct {
	Host string `json:"host"`
}

type TunnelHandler struct {
	db               *gorm.DB
	slogger          *zap.SugaredLogger
	streamingHandler *StreamingHandler
}

func NewTunnelHandler(logger *zap.Logger, db *gorm.DB, streamingHandler *StreamingHandler) TunnelHandler {
	return TunnelHandler{
		db:               db,
		slogger:          logger.Sugar(),
		streamingHandler: streamingHandler,
	}
}

func (s *TunnelHandler) TunnelRequest(c *gin.Context) error {
	host := c.Request.Host
	if conn := s.streamingHandler.getStreamForHostname(host); conn == nil {
		c.JSON(http.StatusNotFound, &TunnelRequestResp{Host: host})
		return nil
	}

	c.JSON(http.StatusOK, &TunnelRequestResp{Host: host})
	return nil
}

// getUrlForRequest returns the URL in the format <schema>://<host><uri>
//
// example: http://localhost:8081/healthz
func getUrlForRequest(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// Retrieve host (e.g., example.com)
	host := c.Request.Host

	// Retrieve request URI (e.g., /path?query=value)
	uri := c.Request.RequestURI

	// Construct the full URL
	return fmt.Sprintf("%s://%s%s", scheme, host, uri)

}
