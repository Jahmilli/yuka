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
	Db             *gorm.DB
	Slogger        *zap.SugaredLogger
	ConnectionPool *StreamingConnectionPool
}

func NewStreamingHandler(logger *zap.Logger, db *gorm.DB) StreamingHandler {
	return StreamingHandler{
		Db:             db,
		Slogger:        logger.Sugar(),
		ConnectionPool: NewStreamingConnectionPool(),
	}
}

func (s *StreamingHandler) InitialiseStream(c *gin.Context) error {
	s.Slogger.Info("Stream initialized")
	domainHeader := c.GetHeader(consts.YukaHeaderDomain)

	// Print the header value or handle if it's missing
	if domainHeader == "" {
		s.Slogger.Warnf("Domain header is missing for new connection")
		c.JSON(400, gin.H{"error": "domainHeader header is missing"})
		return nil
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		s.Slogger.Errorf("Error when upgrading connection request %s", err)
		return nil
	}
	s.ConnectionPool.AddConnection(domainHeader, conn)

	defer func() {
		if err := conn.Close(); err != nil {
			s.Slogger.Errorf("Error when closing connection %s", err)
		}
		s.ConnectionPool.RemoveConnection(domainHeader)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			// Occurs when a connection is closed (Can possibly look at a better way to handle this later...)
			if strings.Contains(err.Error(), "close 1000 (normal)") {
				s.Slogger.Info("Client websocket connection closed")
			} else {
				s.Slogger.Errorf("Error occurred when reading message %s", err)
			}
			return err
		}

		s.Slogger.Infof("Received message %s", message)

		err = conn.WriteMessage(websocket.TextMessage, []byte("Pong"))
		if err != nil {
			s.Slogger.Errorf("write:", err)
			break
		}
		time.Sleep(time.Second)
	}
	return nil
}
