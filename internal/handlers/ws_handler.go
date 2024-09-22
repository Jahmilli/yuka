package handlers

import (
	"strings"
	"time"
	"yuka/pkg/streaming_connection"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// WsHandler just exposes a basic websocket gateway that returns responses
type WsHandler struct {
	db             *gorm.DB
	slogger        *zap.SugaredLogger
	connectionPool *streaming_connection.StreamingConnectionPool
}

func NewWsHandler(logger *zap.Logger, db *gorm.DB) WsHandler {
	return WsHandler{
		db:      db,
		slogger: logger.Sugar(),
	}
}

func (s *WsHandler) HandleWsConnection(c *gin.Context) error {
	s.slogger.Info("Received new WS connection")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		s.slogger.Errorf("Error when upgrading connection request: %v", err)
		return nil
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			// Occurs when a connection is closed (Can possibly look at a better way to handle this later...)
			if strings.Contains(err.Error(), "close 1000 (normal)") {
				s.slogger.Info("Client websocket connection closed")
			} else {
				s.slogger.Errorf("Error occurred when reading message: %v", err)
			}
			return err
		}

		s.slogger.Infof("Received message %s", message)

		err = conn.WriteMessage(websocket.TextMessage, []byte("Pong"))
		if err != nil {
			s.slogger.Errorf("Write error %v", err)
			break
		}
		time.Sleep(time.Second)
	}
	return nil
}
