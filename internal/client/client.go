package client

import (
	"context"
	"os"
	"yuka/internal/api/api_clients"

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
	tunnel := NewTunnel(c.Logger, "localhost:8085", "localhost:5432")
	if err := tunnel.Connect(ctx); err != nil {
		c.slogger.Errorf("Error occurred when listening on tunnel: %v", err)
		return err
	}

	return nil
}

func (c *Client) Cleanup(ctx context.Context) error {
	return nil
}
