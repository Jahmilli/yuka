package handlers

import (
	"strings"
	"time"
	"yuka/internal/consts"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type StreamingHandler struct {
	db             *gorm.DB
	slogger        *zap.SugaredLogger
	connectionPool *StreamingConnectionPool
}

func NewStreamingHandler(logger *zap.Logger, db *gorm.DB) StreamingHandler {
	return StreamingHandler{
		db:             db,
		slogger:        logger.Sugar(),
		connectionPool: NewStreamingConnectionPool(logger),
	}
}

func (s *StreamingHandler) InitializeStream(c *gin.Context) error {
	s.slogger.Info("stream initialized")
	hostnameHeader := c.GetHeader(consts.YukaHeaderHostname)

	// Print the header value or handle if it's missing
	if hostnameHeader == "" {
		s.slogger.Warnf("hostname header is missing for new connection")
		c.JSON(400, gin.H{"error": "Hostname header is missing"})
		return nil
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		s.slogger.Errorf("error when upgrading connection request: %v", err)
		return nil
	}
	s.connectionPool.AddConnection(hostnameHeader, conn)

	defer func() {
		if err := conn.Close(); err != nil {
			s.slogger.Errorf("error when closing connection: %v", err)
		}
		s.connectionPool.RemoveConnection(hostnameHeader)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			// Occurs when a connection is closed (Can possibly look at a better way to handle this later...)
			if strings.Contains(err.Error(), "close 1000 (normal)") {
				s.slogger.Info("client websocket connection closed")
			} else {
				s.slogger.Errorf("error occurred when reading message: %v", err)
			}
			return err
		}

		s.slogger.Infof("received message %s", message)

		err = conn.WriteMessage(websocket.TextMessage, []byte("Pong"))
		if err != nil {
			s.slogger.Errorf("write error %v", err)
			break
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (s *StreamingHandler) getStreamForHostname(hostname string) *Connection {
	conn, err := s.connectionPool.GetConnection(hostname)
	if err != nil {
		s.slogger.Errorf("error occurred when getting connection for hostname %s: %v", hostname, err)
		return nil
	}
	return conn
}
