package http_helper

import (
	"net/http"
)

/*
 * Right now this is a fairly basic struct used just for encoding.
 * Will have to see how this behaves with streaming, multi-part uploads etc.
 * For now, just testing with basic GET, PUT, UPDATE, DELETE requests
 * It's also worth recognizing that this will just be limited to HTTP requests,
 * however, it'll be worth considering how this would work across other transport
 * layers such as WebSockets, gRPC, etc.
 **/

// HttpRequest provides a custom wrapper around the http request
// that we can serialize/deserialize using different encodings
// between the server and client
type HttpRequest struct {
	Method  string      `json:"method"`
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
	// TODO: This is going to be limited if we want any sort of streaming but it should be okay for initial approach
	Body []byte `json:"body"`
}

func NewHttpRequest(method string, url string, headers http.Header, body []byte) *HttpRequest {
	return &HttpRequest{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
	}
}
