package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"go.uber.org/zap"
)

type ClientServer struct {
	slogger zap.SugaredLogger
	port    int
}

func NewClientServer(logger *zap.Logger, port int) *ClientServer {
	return &ClientServer{
		slogger: *logger.Sugar(),
		port:    port,
	}

}

// Listen is a blocking call that starts up the TCP server
//
// Will close on ctx.Done() being called
func (self *ClientServer) Listen(ctx context.Context) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", self.port))

	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	self.slogger.Infof("Server is listening on port %v", self.port)

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
			self.slogger.Info("Shutting down server...")
			return

		case conn := <-connChan:
			// Handle the new connection
			go self.forwardConnection(conn)
			// go self.handleConnection(conn)

		case err := <-errChan:
			// Handle accept error (usually indicates the server should shut down)
			self.slogger.Warnf("Error accepting connection: %v", err)
			return
		}
	}
}

func (self *ClientServer) forwardConnection(conn net.Conn) error {
	forwardConn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		self.slogger.Errorf("Error occurred when dialing connection: %v", err)
		conn.Close()
		return err
	}

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

// Function to handle each connection
func (self *ClientServer) logRequest(conn net.Conn) error {
	// Ensure the connection is closed after we're done
	defer conn.Close()
	// Use a buffered reader to read data from the connection
	reader := bufio.NewReader(conn)
	for {
		// Read incoming data until a newline (or connection closed)
		message, err := reader.ReadString('\n')
		if err != nil {
			self.slogger.Info("Connection closed")
			break
		}

		// Print the received message
		self.slogger.Infof("Received message: %s", message)
	}
	return nil
}

func (self *ClientServer) Close() {

}
