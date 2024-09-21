package handlers

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TunnelRequestResp struct {
	Host string `json:"host"`
}

type TunnelHandler struct {
	db               *gorm.DB
	slogger          *zap.SugaredLogger
	streamingHandler *StreamingHandler
}

func NewTunnelHandler(logger *zap.Logger, db *gorm.DB, streamingHandler *StreamingHandler) TunnelHandler {
	return TunnelHandler{
		db:               db,
		slogger:          logger.Sugar(),
		streamingHandler: streamingHandler,
	}
}

// func (s *TunnelHandler) TunnelRequest(c *gin.Context) error {
// 	host := c.Request.Host
// 	connection := s.streamingHandler.getStreamForHostname(host)
// 	if connection == nil {
// 		c.JSON(http.StatusNotFound, &TunnelRequestResp{Host: host})
// 		return nil
// 	}
// 	url := getUrlForRequest(c)
// 	// We convert this to a known struct which we can serialise between the server and client

// 	var buf bytes.Buffer
// 	_, err := io.Copy(&buf, c.Request.Body)
// 	if err != nil {
// 		s.slogger.Errorf("Error reading request body: %v", err)
// 		c.JSON(http.StatusInternalServerError, &TunnelRequestResp{Host: host})
// 		return nil
// 	}
// 	reqBody := buf.Bytes()
// 	s.slogger.Debugf("Received request body n %v, %s", len(reqBody), string(reqBody))

// 	httpRequest := http_helper.NewHttpRequest(c.Request.Method, url, c.Request.Header, reqBody)
// 	b, err := httpRequest.ToJSON()
// 	if err != nil {
// 		s.slogger.Errorf("error occurred when serializing http request to json: %v", err)
// 		c.JSON(http.StatusInternalServerError, &TunnelRequestResp{Host: host})
// 	}
// 	if err = connection.Conn.WriteMessage(websocket.TextMessage, b); err != nil {
// 		s.slogger.Errorf("write error %v", err)
// 		c.JSON(http.StatusInternalServerError, &TunnelRequestResp{Host: host})
// 	}

// 	c.JSON(http.StatusOK, &TunnelRequestResp{Host: host})
// 	return nil
// }

func (self *TunnelHandler) TunnelRequest(c *gin.Context) error {
	// host := c.Request.Host
	conn, err := net.Dial("tcp", "localhost:8082")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to connect to TCP server"})
		return err
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(30 * time.Second)) // Set a 30-second read timeout
	// Required so the response doesn't stall
	c.Request.Header.Set("Connection", "Close")
	sendRequestMetadata(conn, c)
	println("Copying data")
	// Stream the request body to the TCP server
	_, err = io.Copy(conn, c.Request.Body)
	println("Finished copying data")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream body to TCP server"})
		return err
	}

	// Stream data from the TCP connection directly to the response writer
	if _, err := io.Copy(c.Writer, conn); err != nil {
		self.slogger.Errorf("Error streaming response: %v", err)
		c.JSON(500, gin.H{"error": "Error streaming response"})
		return nil
	}

	// // Optionally, you can log that the streaming has completed
	self.slogger.Info("Streaming completed successfully.")
	c.Status(200) // Send a 200 OK status after completion

	// TODO: Return response back in gin request

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

func sendRequestMetadata(conn net.Conn, c *gin.Context) {
	// Write the request method and URL
	fmt.Fprintf(conn, "%s %s %s\r\n", c.Request.Method, c.Request.URL.Path, c.Request.Proto)

	// Write the headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			fmt.Fprintf(conn, "%s: %s\r\n", key, value)
		}
	}

	// End of headers
	fmt.Fprint(conn, "\r\n")
}
