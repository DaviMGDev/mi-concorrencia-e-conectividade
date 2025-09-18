package handlers

import (
	"client-of-dispair/internal/api"
	"client-of-dispair/internal/api/protocol"
	"client-of-dispair/internal/state"
	"client-of-dispair/internal/ui"
	"time"
)

func SendMessageHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	if len(arguments) < 1 {
		chat.Outputs <- "Did you forget to type a message?"
		return
	}
	message := arguments[0]
	request := protocol.Request{
		Method:    "message",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]any{
			"user_id": state.UserID,
			"room_id": state.RoomID,
			"message": message,
		},
	}
	err := client.SendRequest(&request)
	if err != nil {
		chat.Outputs <- "Error sending message request: " + err.Error()
		return
	}
	response, err := client.ReceiveResponse()
	if err != nil {
		chat.Outputs <- "Error receiving message response: " + err.Error()
		return
	}
	if response.Status != "ok" {
		chat.Outputs <- "Message sending failed: " + response.Message
		return
	}
	chat.Outputs <- "Message sent successfully."
}

func GetMessageHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	request := protocol.Request{
		Method:    "get",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]any{
			"user_id": state.UserID,
			"room_id": state.RoomID,
		},
	}
	err := client.SendRequest(&request)
	if err != nil {
		chat.Outputs <- "Error sending get message request: " + err.Error()
		return
	}
	response, err := client.ReceiveResponse()
	if err != nil {
		chat.Outputs <- "Error receiving get message response: " + err.Error()
		return
	}
	if response.Status != "ok" {
		chat.Outputs <- "Get message failed: " + response.Data["message"].(string)
		return
	}
	message, ok := response.Data["message"].(string)
	if !ok {
		chat.Outputs <- "Invalid message format in response"
		return
	}
	chat.Outputs <- message
}
