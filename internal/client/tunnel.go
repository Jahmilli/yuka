package client

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"yuka/pkg/streaming_connection"

	"go.uber.org/zap"
)

type Tunnel struct {
	slogger         zap.SugaredLogger
	serverHostname  string
	forwardHostname string
}

func NewTunnel(logger *zap.Logger, serverHostname string, forwardHostname string) *Tunnel {
	return &Tunnel{
		slogger:         *logger.Sugar(),
		serverHostname:  serverHostname,
		forwardHostname: forwardHostname,
	}

}

// Listen is a blocking call that starts up the TCP server
//
// Will close on ctx.Done() being called
func (self *Tunnel) Connect(ctx context.Context) error {
	// Setup connection to server
	conn, err := net.Dial("tcp", self.serverHostname)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Register client with metadata
	metadata := streaming_connection.NewConnectionMetadata("seb-hostname")
	b, err := metadata.Serialize()
	if err != nil {
		return err
	}
	// Send the size of the metadata first
	metadataSize := int32(len(b))
	if err := binary.Write(conn, binary.BigEndian, metadataSize); err != nil {
		return fmt.Errorf("error writing metadata size: %v", err)
	}

	// Send the actual metadata
	if _, err := conn.Write(b); err != nil {
		return fmt.Errorf("error writing metadata: %v", err)
	}
	// conn.Write(b)

	// self.slogger.Infof("Tunnel is listening on port %v", self.listenPort)

	// Channel to signal new connections
	// connChan := make(chan net.Conn)
	errChan := make(chan error)
	go self.forwardConnection(conn)
	for {
		select {
		case <-ctx.Done():
			// The context has been canceled, stop accepting new connections
			self.slogger.Info("Shutting down connection...")
			return nil

		case err := <-errChan:
			// Handle accept error (usually indicates the server should shut down)
			self.slogger.Warnf("Error accepting connection: %v", err)
		}
	}
}

func (self *Tunnel) forwardConnection(conn net.Conn) error {
	forwardConn, err := net.Dial("tcp", self.forwardHostname)
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
