package streaming_connection

import (
	"context"
	"fmt"
	"io"
	"net"

	"go.uber.org/zap"
)

// TcpServer starts up a basic TCP server and supports forwarding those connections on
type TcpServer struct {
	slogger        zap.SugaredLogger
	listenPort     int
	connectionPool *StreamingConnectionPool
}

func NewTcpServer(logger *zap.Logger, listenPort int, connectionPool *StreamingConnectionPool) *TcpServer {
	return &TcpServer{
		slogger:        *logger.Sugar(),
		listenPort:     listenPort,
		connectionPool: connectionPool,
	}
}

// Listen is a blocking call that starts up the TCP server
//
// Will close on ctx.Done() being called
func (self *TcpServer) Listen(ctx context.Context) error {
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
			go self.TunnelRequest(conn)

		case err := <-errChan:
			// Handle accept error (usually indicates the server should shut down)
			self.slogger.Errorf("Error accepting connection: %v", err)
		}
	}
}

func (self *TcpServer) forwardConnection(conn net.Conn, forwardConn StreamingConnection) error {
	go func() {
		self.slogger.Info("Forwarding data from forwardConn to conn.")
		if _, err := io.Copy(conn, forwardConn); err != nil {
			self.slogger.Errorf("Error copying from forwardConn to conn: %v", err)
		}
		self.slogger.Info("Finished copying data from forwardConn to conn")
		conn.Close() // Close after copying
	}()
	go func() {
		self.slogger.Info("Forwarding data from conn to forwardConn.")
		if _, err := io.Copy(forwardConn, conn); err != nil {
			self.slogger.Errorf("Error copying from conn to forwardConn: %v", err)
		}
		self.slogger.Info("Finished copying data from conn to forwardConn")
		forwardConn.Close() // Close after copying
	}()

	return nil
}

func (self *TcpServer) TunnelRequest(conn net.Conn) error {
	// url := getUrlForRequest(c)
	registeredHostname := "seb-hostname"
	// TODO: This should come from the request
	connection, err := self.connectionPool.GetConnection(registeredHostname)
	self.slogger.Infof("Got connection for hostname %s", registeredHostname)
	if err != nil {
		self.slogger.Warnf("Received error when getting connection for hostname %s: %v", registeredHostname, err)
	}

	self.forwardConnection(conn, connection)

	self.slogger.Info("Streaming completed successfully.")
	return nil
}
