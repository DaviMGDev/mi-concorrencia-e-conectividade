package api

import (
	"client-of-dispair/internal/api/protocol"
	"encoding/json"
	"net"
	"sync"
	"time"
)

type Client struct {
	Address     string
	Connection  net.Conn
	Mutex       sync.Mutex
	Encoder     *json.Encoder
	Decoder     *json.Decoder
	lastRequest *protocol.Request
}

func NewClient(address string) *Client {
	return &Client{
		Address: address,
	}
}

func (c *Client) Connect() error {
	return nil
}

func (c *Client) Close() error {
	return nil
}

func (c *Client) SendRequest(req *protocol.Request) error {
	c.lastRequest = req
	return nil
}

func (c *Client) ReceiveResponse() (*protocol.Response, error) {
	if c.lastRequest == nil {
		return &protocol.Response{
			Method:    "mock",
			Status:    "ok",
			Message:   "This is a mocked response",
			Timestamp: time.Now().Format(time.RFC3339),
		}, nil
	}

	switch c.lastRequest.Method {
	case "ping":
		return &protocol.Response{
			Method:    "ping",
			Status:    "ok",
			Timestamp: c.lastRequest.Timestamp,
		}, nil
	case "register", "login", "create", "join", "leave":
		return &protocol.Response{
			Method:    c.lastRequest.Method,
			Status:    "ok",
			Timestamp: time.Now().Format(time.RFC3339),
			Data:      map[string]any{"user_id": "mock_user_id", "room_id": "mock_room_id"},
		}, nil
	case "message", "choice":
		return &protocol.Response{
			Method:    c.lastRequest.Method,
			Status:    "ok",
			Timestamp: time.Now().Format(time.RFC3339),
		}, nil
	case "get":
		return &protocol.Response{
			Method:    "get",
			Status:    "ok",
			Timestamp: time.Now().Format(time.RFC3339),
			Data:      map[string]any{"message": "hello from mock server"},
		}, nil
	case "get_opponent_choice":
		return &protocol.Response{
			Method:    "get_opponent_choice",
			Status:    "ok",
			Timestamp: time.Now().Format(time.RFC3339),
			Data:      map[string]any{"opponent_choice": "rock"},
		}, nil
	default:
		return &protocol.Response{
			Method:    c.lastRequest.Method,
			Status:    "error",
			Message:   "Unknown method",
			Timestamp: time.Now().Format(time.RFC3339),
		}, nil
	}
}

