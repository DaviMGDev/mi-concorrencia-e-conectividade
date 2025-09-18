package protocol

import (
	"encoding/json"
	"net"
)

type Client struct {
	Address    string
	Connection net.Conn
	Encoder    *json.Encoder
	Decoder    *json.Decoder
}

func NewClient(address string, conn net.Conn) *Client {
	return &Client{
		Address:    address,
		Connection: conn,
		Encoder:    json.NewEncoder(conn),
		Decoder:    json.NewDecoder(conn),
	}
}

func (c *Client) Write(response *Response) error {
	return c.Encoder.Encode(response)
}

func (c *Client) Read() (*Request, error) {
	var request Request
	err := c.Decoder.Decode(&request)
	if err != nil {
		return nil, err
	}
	request.From = c.Address
	return &request, nil
}

func (c *Client) Close() error {
	return c.Connection.Close()
}