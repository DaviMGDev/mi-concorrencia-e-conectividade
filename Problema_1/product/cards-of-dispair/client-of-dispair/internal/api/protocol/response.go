package protocol

type Response struct {
	Method    string         `json:"method"`
	Status    string         `json:"status"`
	Message   string         `json:"message"`
	Timestamp string         `json:"timestamp"`
	Data      map[string]any `json:"data"`
}
