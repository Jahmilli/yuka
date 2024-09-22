package streaming_connection

import "io"

// StreamingConnection is an interface for handling any type of streaming connection.
type StreamingConnection interface {
	// Write sends data over the connection.
	Write([]byte) (int, error)

	// Read reads data from the connection.
	Read([]byte) (int, error)

	// Close closes the connection.
	Close() error

	// IsOpen checks if the connection is still open.
	IsOpen() bool

	// GetReader returns an io.Reader for the connection (for streaming use cases).
	GetReader() io.Reader

	// GetWriter returns an io.Writer for the connection (for streaming use cases).
	GetWriter() io.Writer
}
