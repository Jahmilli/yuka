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
	"yuka/internal/consts"

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
	if err := c.initializeConnection(ctx); err != nil {
		return err
	}

	return nil
}

func (client *Client) initializeConnection(ctx context.Context) error {
	u := url.URL{Scheme: "ws", Host: client.apiServerAddress, Path: "/ws"}
	headers := http.Header{}
	headers.Add(consts.YukaHeaderHostname, "localhost:8081")

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
				client.slogger.Errorf("error when reading message: %v", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second * 10)
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
				client.slogger.Errorf("error when writing message: %v", err)
				return nil
			}
		case <-ctx.Done():
			client.Logger.Info("interrupt called")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				client.slogger.Errorf("write close: %v", err)
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
