package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WsWrapper struct {
	slogger *zap.SugaredLogger
}

func NewWsWrapper(logger *zap.Logger) *WsWrapper {
	return &WsWrapper{
		slogger: logger.Sugar(),
	}
}

func (self *WsWrapper) StartWs(ctx context.Context, wsHost string) error {
	if err := self.initializeWsConnection(ctx, wsHost); err != nil {
		return err
	}
	return nil
}

func (self *WsWrapper) initializeWsConnection(ctx context.Context, wsHost string) error {
	u := url.URL{Scheme: "ws", Host: wsHost, Path: "/ws"}
	headers := http.Header{}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return fmt.Errorf("failed to initialize websocket connection: %v", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				self.slogger.Errorf("Error when reading message: %v", err)
				return
			}

			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	counter := 0
	for {
		select {
		case <-done:
			return nil
		case _ = <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%v", counter)))
			counter++
			if err != nil {
				self.slogger.Errorf("Error when writing message: %v", err)
				return nil
			}
		case <-ctx.Done():
			self.slogger.Info("Interrupt called")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				self.slogger.Errorf("Write close: %v", err)
				return nil
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

func (self *WsWrapper) Cleanup(ctx context.Context) error {
	return nil
}
