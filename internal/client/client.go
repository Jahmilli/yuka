package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
	"yuka/internal/api/api_clients"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

type Client struct {
	apiServerAddress string
	ApiClient        *api_clients.Yuka
	Logger           *zap.Logger
	slogger          *zap.SugaredLogger
	Hostname         string
}

func NewClient(apiserverAddress string, logger *zap.Logger, hostname string) *Client {
	transport := httptransport.New(apiserverAddress, "", nil)
	transport.DefaultAuthentication = httptransport.BasicAuth(os.Getenv("HTTP_USERNAME"), os.Getenv("HTTP_PASSWORD"))
	return &Client{
		apiServerAddress: apiserverAddress,
		ApiClient:        api_clients.New(transport, strfmt.Default),
		Logger:           logger,
		slogger:          logger.Sugar(),
		Hostname:         hostname,
	}
}

func (c *Client) Start(ctx context.Context) error {
	tunnel := NewTunnel(c.Logger, 8082, "localhost:8080")
	if err := tunnel.Listen(ctx); err != nil {
		c.slogger.Errorf("Error occurred when listening on tunnel: %v", err)
		return err
	}

	return nil
}

func (c *Client) StartWs(ctx context.Context, wsHost string) error {
	if err := c.initializeWsConnection(ctx, wsHost); err != nil {
		return err
	}
	return nil
}

func (client *Client) initializeWsConnection(ctx context.Context, wsHost string) error {
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
				client.slogger.Errorf("Error when reading message: %v", err)
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
				client.slogger.Errorf("Error when writing message: %v", err)
				return nil
			}
		case <-ctx.Done():
			client.Logger.Info("Interrupt called")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				client.slogger.Errorf("Write close: %v", err)
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

func (c *Client) Cleanup(ctx context.Context) error {
	return nil
}
