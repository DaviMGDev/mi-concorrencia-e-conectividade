package handlers

import (
	"client-of-hope/internal/api"
	"client-of-hope/internal/api/protocol"
	"client-of-hope/internal/state"
	"client-of-hope/internal/ui"
	"client-of-hope/internal/utils"
	"fmt"
	"strings"
)

func HandleSendMessage(client *api.Client, chat *ui.Chat, args []string) {
	if state.UserID == "" || state.RoomID == "" {
		chat.Outputs <- "You must be logged in and in a room to send messages."
		return
	}
	if len(args) < 1 {
		chat.Outputs <- "Usage: Just type your message and press Enter to send."
		return
	}

	message := strings.Join(args, " ")
	request := protocol.Request{
		Method: "send",
		Data: utils.Dict{
			"user_id": state.UserID,
			"room_id": state.RoomID,
			"message": fmt.Sprintf("%s: %s", state.Username, message),
		},
	}

	_, err := client.DoRequest(request)
	if err != nil {
		state.Log("Send message request failed: %v", err)
		chat.Outputs <- "Failed to send message."
	}
}

func HandleFetchMessage(client *api.Client, chat *ui.Chat, args []string) {
	if state.UserID == "" || state.RoomID == "" {
		return // Don't show any error, just fail silently
	}

	request := protocol.Request{
		Method: "fetch",
		Data: utils.Dict{
			"user_id": state.UserID,
			"room_id": state.RoomID,
		},
	}

	response, err := client.DoRequest(request)
	if err != nil {
		state.Log("Fetch messages failed: %v", err)
		return
	}

	if response.Status != "ok" {
		return
	}

	textMessage, ok := response.Data["message"].(string)
	if !ok || textMessage == "" {
		return
	}

	chat.Outputs <- textMessage
}