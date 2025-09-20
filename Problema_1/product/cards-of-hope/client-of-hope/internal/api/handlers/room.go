package handlers

import (
	"client-of-hope/internal/api"
	"client-of-hope/internal/api/protocol"
	"client-of-hope/internal/state"
	"client-of-hope/internal/ui"
	"client-of-hope/internal/utils"
	"fmt"
)

func HandleCreateRoom(client *api.Client, chat *ui.Chat, args []string) {
	if state.UserID == "" {
		chat.Outputs <- "You must be logged in to create a room."
		return
	}

	request := protocol.Request{
		Method: "create",
		Data:   utils.Dict{"user_id": state.UserID},
	}

	response, err := client.DoRequest(request)
	if err != nil {
		state.Log("Create room request failed: %v", err)
		chat.Outputs <- "Failed to create room."
		return
	}

	if response.Status != "ok" {
		message, _ := response.Data["message"].(string)
		chat.Outputs <- message
		return
	}

	roomID, ok := response.Data["room_id"].(string)
	if !ok {
		chat.Outputs <- "Invalid room ID from server."
		return
	}

	state.RoomID = roomID
	chat.Outputs <- fmt.Sprintf("Room created successfully! Room ID: %s", roomID)
}

func HandleJoinRoom(client *api.Client, chat *ui.Chat, args []string) {
	if state.UserID == "" {
		chat.Outputs <- "You must be logged in to join a room."
		return
	}
	if len(args) < 1 {
		chat.Outputs <- "Usage: /join <room_id>"
		return
	}

	roomID := args[0]
	request := protocol.Request{
		Method: "join",
		Data: utils.Dict{
			"user_id": state.UserID,
			"room_id": roomID,
		},
	}

	response, err := client.DoRequest(request)
	if err != nil {
		state.Log("Join room request failed: %v", err)
		chat.Outputs <- "Failed to join room."
		return
	}

	if response.Status != "ok" {
		message, _ := response.Data["message"].(string)
		chat.Outputs <- message
		return
	}

	state.RoomID = roomID
	chat.Outputs <- fmt.Sprintf("Successfully joined room %s", roomID)
}

func HandleLeaveRoom(client *api.Client, chat *ui.Chat, args []string) {
	if state.UserID == "" || state.RoomID == "" {
		chat.Outputs <- "You must be logged in and in a room to leave."
		return
	}

	request := protocol.Request{
		Method: "leave",
		Data: utils.Dict{
			"user_id": state.UserID,
			"room_id": state.RoomID,
		},
	}

	response, err := client.DoRequest(request)
	if err != nil {
		state.Log("Leave room request failed: %v", err)
		chat.Outputs <- "Failed to leave room."
		return
	}

	if response.Status != "ok" {
		message, _ := response.Data["message"].(string)
		chat.Outputs <- message
		return
	}

	chat.Outputs <- fmt.Sprintf("Successfully left room %s", state.RoomID)
	state.RoomID = ""
}