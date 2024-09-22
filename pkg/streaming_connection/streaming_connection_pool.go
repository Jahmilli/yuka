package streaming_connection

import (
	"errors"

	"go.uber.org/zap"
)

var (
	ErrConnectionNotFound = errors.New("No connection found")
)

type StreamingConnectionPool struct {
	slogger *zap.SugaredLogger
	// TODO: Make this interface that supports any type of streaming connection, not just websockets
	connections map[string]StreamingConnection
}

// NewStreamingConnectionPool provides an interface for adding/removing existing streaming connections.
func NewStreamingConnectionPool(logger *zap.Logger) *StreamingConnectionPool {
	return &StreamingConnectionPool{
		connections: make(map[string]StreamingConnection),
		slogger:     logger.Sugar(),
	}
}

func (c *StreamingConnectionPool) AddConnection(hostname string, conn StreamingConnection) {
	c.slogger.Debugf("Adding connection for hostname %s", hostname)
	c.connections[hostname] = conn
}

func (c *StreamingConnectionPool) GetConnection(hostname string) (StreamingConnection, error) {
	c.slogger.Debugf("Getting connection for hostname %s", hostname)
	conn, ok := c.connections[hostname]
	if !ok {
		return nil, ErrConnectionNotFound
	}
	return conn, nil
}

func (c *StreamingConnectionPool) RemoveConnection(hostname string) {
	c.slogger.Debugf("Removing connection for hostname %s", hostname)
	delete(c.connections, hostname)
}
