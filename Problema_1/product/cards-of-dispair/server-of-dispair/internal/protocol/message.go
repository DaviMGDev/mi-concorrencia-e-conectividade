package protocol

import "time"

// Request represents an incoming message from a client.
type Request struct {
	Method    string         `json:"method"`
	Timestamp string         `json:"timestamp"`
	Data      map[string]any `json:"data"`

	// From is the client's address, not part of the JSON payload.
	From string `json:"-"`
}

// Response represents an outgoing message to a client.
type Response struct {
	Method    string         `json:"method"`
	Status    string         `json:"status"`
	Message   string         `json:"message"`
	Timestamp string         `json:"timestamp"`
	Data      map[string]any `json:"data"`

	// To is the client's address, not part of the JSON payload.
	To string `json:"-"`
}

func NewResponse(to, method, status, message string, data map[string]any) *Response {
	return &Response{
		To:        to,
		Method:    method,
		Status:    status,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      data,
	}
}