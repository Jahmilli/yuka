package handlers

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Connection struct {
	conn *websocket.Conn
}

type StreamingConnectionPool struct {
	Slogger *zap.SugaredLogger
	// TODO: Make this interface that supports any type of streaming connection, not just websockets
	connections map[string]*Connection
}

// NewStreamingConnectionPool provides an interface for adding/removing existing streaming connections.
func NewStreamingConnectionPool() *StreamingConnectionPool {
	return &StreamingConnectionPool{
		connections: map[string]*Connection{},
	}
}

func (c *StreamingConnectionPool) AddConnection(domain string, conn *websocket.Conn) {
	c.Slogger.Debugf("Adding connection for domain %s", domain)
	c.connections[domain] = &Connection{
		conn,
	}
}
func (c *StreamingConnectionPool) GetConnection(domain string) *Connection {
	c.Slogger.Debugf("Getting connection for domain %s", domain)
	return c.connections[domain]
}

func (c *StreamingConnectionPool) RemoveConnection(domain string) {
	c.Slogger.Debugf("Removing connection for domain %s", domain)
	delete(c.connections, domain)
}
