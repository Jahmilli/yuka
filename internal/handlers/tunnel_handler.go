package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"yuka/pkg/http_helper"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	connection := s.streamingHandler.getStreamForHostname(host)
	if connection == nil {
		c.JSON(http.StatusNotFound, &TunnelRequestResp{Host: host})
		return nil
	}
	url := getUrlForRequest(c)
	// We convert this to a known struct which we can serialise between the server and client

	var buf bytes.Buffer
	_, err := io.Copy(&buf, c.Request.Body)
	if err != nil {
		s.slogger.Errorf("Error reading request body: %v", err)
		c.JSON(http.StatusInternalServerError, &TunnelRequestResp{Host: host})
		return nil
	}
	reqBody := buf.Bytes()
	s.slogger.Debugf("Received request body n %v, %s", len(reqBody), string(reqBody))

	httpRequest := http_helper.NewHttpRequest(c.Request.Method, url, c.Request.Header, reqBody)
	b, err := httpRequest.ToJSON()
	if err != nil {
		s.slogger.Errorf("error occurred when serializing http request to json: %v", err)
		c.JSON(http.StatusInternalServerError, &TunnelRequestResp{Host: host})
	}
	if err = connection.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
		s.slogger.Errorf("write error %v", err)
		c.JSON(http.StatusInternalServerError, &TunnelRequestResp{Host: host})
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
