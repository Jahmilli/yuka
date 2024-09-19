package http_helper

import "encoding/json"

// Serialize the HttpRequest struct to JSON
func (req *HttpRequest) ToJSON() ([]byte, error) {
	// Marshal the struct into a JSON byte slice
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// FromJSON deserializes a JSON byte slice into an HttpRequest struct
func FromJSON(jsonData []byte) (*HttpRequest, error) {
	var req HttpRequest
	// Unmarshal the JSON data into the struct
	err := json.Unmarshal(jsonData, &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}
