package handlers

import (
	"fmt"
	"io"
	"net/http"
	"yuka/pkg/streaming_connection"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TunnelRequestResp struct {
	Host string `json:"host"`
}

type TunnelHandler struct {
	db             *gorm.DB
	slogger        *zap.SugaredLogger
	connectionPool streaming_connection.StreamingConnectionPool
}

func NewTunnelHandler(logger *zap.Logger, db *gorm.DB, connectionPool streaming_connection.StreamingConnectionPool) TunnelHandler {
	return TunnelHandler{
		db:             db,
		slogger:        logger.Sugar(),
		connectionPool: connectionPool,
	}
}

// func (self *TunnelHandler) TunnelRequest(c *gin.Context) error {
// 	// host := c.Request.Host
// 	conn, err := net.Dial("tcp", "localhost:8082")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to TCP server"})
// 		return err
// 	}
// 	defer conn.Close()
// 	conn.SetReadDeadline(time.Now().Add(30 * time.Second)) // Set a 30-second read timeout
// 	// Required so the response doesn't stall
// 	c.Request.Header.Set("Connection", "Close")
// 	sendRequestMetadata(c, conn)
// 	// Stream the request body to the TCP server
// 	_, err = io.Copy(conn, c.Request.Body)
// 	println("Finished copying data")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream body to TCP server"})
// 		return err
// 	}

// 	// Stream data from the TCP connection directly to the response writer
// 	if _, err := io.Copy(c.Writer, conn); err != nil {
// 		self.slogger.Errorf("Error streaming response: %v", err)
// 		c.JSON(500, gin.H{"error": "Error streaming response"})
// 		return nil
// 	}

// 	// // Optionally, you can log that the streaming has completed
// 	self.slogger.Info("Streaming completed successfully.")
// 	c.Status(200) // Send a 200 OK status after completion

// 	// TODO: Return response back in gin request

// 	return nil
// }

func (self *TunnelHandler) TunnelRequest(c *gin.Context) error {
	// url := getUrlForRequest(c)
	fmt.Println("Host is ", c.Request.Host)
	registeredHostname := "seb-hostname"
	// TODO: This should come from the request
	connection, err := self.connectionPool.GetConnection(registeredHostname)
	if err != nil {
		self.slogger.Warnf("Received error when getting connection for hostname %s: %v", registeredHostname, err)
	}
	c.Request.Header.Set("Connection", "Close")
	writeRequestMetadata(c, connection.GetWriter())

	// Stream the request body to the TCP server
	_, err = io.Copy(connection.GetWriter(), c.Request.Body)
	println("Finished copying data")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream body to TCP server"})
		return err
	}

	// Stream data from the connection directly to the response writer
	if _, err := io.Copy(c.Writer, connection.GetReader()); err != nil {
		self.slogger.Errorf("Error streaming response: %v", err)
		c.JSON(500, gin.H{"error": "Error streaming response"})
		return nil
	}

	self.slogger.Info("Streaming completed successfully.")
	c.Status(200) // Send a 200 OK status after completion

	return nil
}

// getUrlForRequest returns the URL in the format <schema>://<host><uri>
//
// example: http://localhost:8081/healthz
func getUrlForRequest(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// Retrieve host (e.g., example.com)
	host := c.Request.Host

	// Retrieve request URI (e.g., /path?query=value)
	uri := c.Request.RequestURI

	// Construct the full URL
	return fmt.Sprintf("%s://%s%s", scheme, host, uri)

}

func writeRequestMetadata(c *gin.Context, writer io.Writer) {
	// Write the request method and URL
	fmt.Fprintf(writer, "%s %s %s\r\n", c.Request.Method, c.Request.URL.Path, c.Request.Proto)

	// Write the headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			fmt.Fprintf(writer, "%s: %s\r\n", key, value)
		}
	}

	// End of headers
	fmt.Fprint(writer, "\r\n")
}
