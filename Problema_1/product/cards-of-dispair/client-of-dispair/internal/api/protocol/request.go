package protocol

type Request struct {
	Method    string         `json:"method"`
	Timestamp string         `json:"timestamp"`
	Data      map[string]any `json:"data"`
}
