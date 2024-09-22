package streaming_connection

import (
	"io"
	"net"
)

// TcpStreamingConnection abstracts a general TCP connection implementing StreamingConnection
type TcpStreamingConnection struct {
	tcpConn  net.Conn
	isOpen   bool
	metadata *ConnectionMetadata
}

// NewTcpStreamingConnection builds a TcpStreamingConnection and validates the correct metadata is sent before returning.
// If invalid metadata is sent or no metadata is sent before a given timeout, this will error
func NewTcpStreamingConnection(tcpConn net.Conn) (*TcpStreamingConnection, error) {
	conn := TcpStreamingConnection{
		tcpConn: tcpConn,
		isOpen:  true,
	}
	// We read in the metadata here before returning so we can at least ensure all connections are valid.
	// Also as the metadata is sent in the first chunk of bytes sent over the wire this prevents any race conditions
	// of accidentally reading it from elsewhere.
	metadata, err := ReadMetadataFromNewConnection(&conn)
	if err != nil {
		return nil, err
	}
	conn.metadata = metadata
	return &conn, nil
}

// Write writes data to the TCP connection.
func (self *TcpStreamingConnection) Write(data []byte) (int, error) {
	if !self.IsOpen() {
		return 0, io.EOF // Return EOF if the connection is closed.
	}
	return self.tcpConn.Write(data)
}

// Read reads data from the TCP connection.
func (self *TcpStreamingConnection) Read(buffer []byte) (int, error) {
	if !self.IsOpen() {
		return 0, io.EOF // Return EOF if the connection is closed.
	}
	return self.tcpConn.Read(buffer)
}

// Close closes the TCP connection.
func (self *TcpStreamingConnection) Close() error {
	if self.IsOpen() {
		self.isOpen = false
		return self.tcpConn.Close()
	}
	return nil
}

// IsOpen checks if the connection is still open.
func (self *TcpStreamingConnection) IsOpen() bool {
	return self.isOpen
}

// GetReader returns an io.Reader for the connection (for streaming use cases).
func (self *TcpStreamingConnection) GetReader() io.Reader {
	return self.tcpConn
}

// GetWriter returns an io.Writer for the connection (for streaming use cases).
func (self *TcpStreamingConnection) GetWriter() io.Writer {
	return self.tcpConn
}
