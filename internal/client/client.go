package client

import (
	"context"
	"time"
	"yuka/internal/api/api_clients"

	"go.uber.org/zap"
)

type Client struct {
	ApiClient      *api_clients.Yuka
	Logger         *zap.Logger
	UserId         string
	OrganizationId string
	TunnelIp       string
	Hostname       string
	ExposedPorts   []string
}

func NewClient(apiClient *api_clients.Yuka, logger *zap.Logger, userId string, organizationId string, tunnelIp string, hostname string, exposedPorts []string) *Client {
	logger.Debug("Client Initialised")
	return &Client{
		ApiClient:      apiClient,
		Logger:         logger,
		UserId:         userId,
		OrganizationId: organizationId,
		TunnelIp:       tunnelIp,
		Hostname:       hostname,
		ExposedPorts:   exposedPorts,
	}
}

func (c *Client) Start(ctx context.Context) error {
	slog := c.Logger.Sugar()
	slog.Infof("Started")

	if err := c.getPeersOnInterval(ctx, 20); err != nil {
		return err
	}

	return nil
}

func (c *Client) getPeersOnInterval(ctx context.Context, intervalSeconds int) error {
	stunTicker := time.NewTicker(time.Second * 20)

	for {
		select {
		case <-ctx.Done():
			c.Logger.Info("Stopping reconciliation of peers for client")
			return nil
		case <-stunTicker.C:
			c.Logger.Info("Stun ticker")

		}
	}
}

func (c *Client) Cleanup(ctx context.Context) error {
	return nil
}
