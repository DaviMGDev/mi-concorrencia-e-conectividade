package protocol

import "server-of-hope/internal/utils"

type Request struct {
	Method string     `json:"method"`
	Data   utils.Dict `json:"data,omitempty"`

	From string `json:"-"`
}
