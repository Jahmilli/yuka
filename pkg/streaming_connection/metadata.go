package streaming_connection

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

const metadataSize = 1024

// Expect Metadata is always in JSON format
type ConnectionMetadata struct {
	RegisteredHostname string `json:"registeredHostname"`
}

func NewConnectionMetadata(registeredHostname string) *ConnectionMetadata {
	return &ConnectionMetadata{
		RegisteredHostname: registeredHostname,
	}
}

func ReadMetadataFromNewConnection(conn StreamingConnection) (*ConnectionMetadata, error) {
	// First, read the size of the incoming metadata
	var metadataSize int32
	if err := binary.Read(conn, binary.BigEndian, &metadataSize); err != nil {
		return nil, fmt.Errorf("error reading metadata size: %v", err)
	}

	// Allocate a buffer of the appropriate size
	metadataBuffer := make([]byte, metadataSize)

	// Read the metadata into the buffer
	n, err := conn.Read(metadataBuffer)
	if err != nil {
		return nil, err
	}

	// Ensure we only use the valid part of the buffer
	if n != int(metadataSize) {
		return nil, fmt.Errorf("expected to read %d bytes, but read %d", metadataSize, n)
	}

	// Parse the metadata
	var metadata map[string]string
	if err := json.Unmarshal(metadataBuffer, &metadata); err != nil {
		return nil, fmt.Errorf("error parsing metadata: %v", err)
	}

	fmt.Println("Metadata is ", metadata)

	// Extract hostname and other metadata information
	registeredHostname := metadata["registeredHostname"]

	return &ConnectionMetadata{
		RegisteredHostname: registeredHostname,
	}, nil
}

// Serialize converts the metadata into JSON and writes it into a buffer of a fixed size
func (self *ConnectionMetadata) Serialize() ([]byte, error) {
	metadataBytes, err := json.Marshal(self)
	if err != nil {
		return nil, err
	}
	return metadataBytes, nil
}
