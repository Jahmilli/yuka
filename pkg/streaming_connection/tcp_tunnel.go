package streaming_connection

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
)

// TcpTunnel is responsible for listening to TCP requests from yukactl clients and then
// adding those connections to a connection pool. This is required for any TCP streaming
type TcpTunnel struct {
	slogger        zap.SugaredLogger
	listenPort     int
	connectionPool *StreamingConnectionPool
}

func NewTcpTunnel(logger *zap.Logger, listenPort int, connectionPool *StreamingConnectionPool) *TcpTunnel {
	return &TcpTunnel{
		slogger:        *logger.Sugar(),
		listenPort:     listenPort,
		connectionPool: connectionPool,
	}
}

// Listen is a blocking call that starts up the TCP server
//
// Will close on ctx.Done() being called
func (self *TcpTunnel) Listen(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", self.listenPort))

	if err != nil {
		return err
	}
	defer listener.Close()

	self.slogger.Infof("TcpTunnel is listening on port %v", self.listenPort)

	// Channel to signal new connections
	connChan := make(chan net.Conn)
	errChan := make(chan error)

	// Start a goroutine to accept connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				errChan <- err
				return
			}
			connChan <- conn
		}
	}()

	for {
		select {
		case <-ctx.Done():
			// The context has been canceled, stop accepting new connections
			self.slogger.Info("Shutting down tunnel...")
			return nil

		case conn := <-connChan:
			// Handle the new connection
			go self.handleNewConnection(conn)

		case err := <-errChan:
			// Handle accept error (usually indicates the server should shut down)
			self.slogger.Errorf("Error accepting connection: %v", err)
		}
	}
}

func (self *TcpTunnel) handleNewConnection(conn net.Conn) error {
	self.slogger.Infof("Handling new connection")
	tcpConn, err := NewTcpStreamingConnection(conn)
	if err != nil {
		self.slogger.Errorf("Error occurred when creating new TcpStreamingConnection: %v", err)
		return err
	}
	self.slogger.Infof("Created new TcpStreamingConnection")

	fmt.Println("metadata is ", tcpConn.metadata)
	registeredHostname := tcpConn.metadata.RegisteredHostname
	self.connectionPool.AddConnection(registeredHostname, tcpConn)

	// return self.forwardConnection(conn)
	return nil
}
