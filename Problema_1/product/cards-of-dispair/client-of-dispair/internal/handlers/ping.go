package handlers

import (
	"client-of-dispair/internal/api"
	"client-of-dispair/internal/api/protocol"
	"client-of-dispair/internal/ui"
	"time"
)

func PingHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	request := protocol.Request{
		Method:    "ping",
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      nil,
	}
	err := client.SendRequest(&request)

	if err != nil {
		chat.Outputs <- "Error sending ping request: " + err.Error()
		return
	}
	response, err := client.ReceiveResponse()
	if err != nil {
		chat.Outputs <- "Error receiving ping response: " + err.Error()
		return
	}
	responseTime, err := time.Parse(time.RFC3339, response.Timestamp)
	if err != nil {
		chat.Outputs <- "Error parsing response timestamp: " + err.Error()
		return
	}
	latency := time.Since(responseTime)
	chat.Outputs <- "Ping response received. Latency: " + latency.String()
}
