package handlers

import (
	"errors"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var (
	ErrConnectionNotFound = errors.New("No connection found")
)

type Connection struct {
	Conn *websocket.Conn
}

type StreamingConnectionPool struct {
	slogger *zap.SugaredLogger
	// TODO: Make this interface that supports any type of streaming connection, not just websockets
	connections map[string]*Connection
}

// NewStreamingConnectionPool provides an interface for adding/removing existing streaming connections.
func NewStreamingConnectionPool(logger *zap.Logger) *StreamingConnectionPool {
	return &StreamingConnectionPool{
		connections: make(map[string]*Connection),
		slogger:     logger.Sugar(),
	}
}

func (c *StreamingConnectionPool) AddConnection(hostname string, conn *websocket.Conn) {
	c.slogger.Debugf("Adding connection for hostname %s", hostname)
	c.connections[hostname] = &Connection{
		conn,
	}
}
func (c *StreamingConnectionPool) GetConnection(hostname string) (*Connection, error) {
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
