package protocol

import "server-of-hope/internal/utils"

type Response struct {
	Method string     `json:"method"`
	Status string     `json:"status"`
	Data   utils.Dict `json:"data,omitempty"`

	To string `json:"-"`
}
