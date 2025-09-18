package handlers

import (
	"client-of-dispair/internal/api"
	"client-of-dispair/internal/api/protocol"
	"client-of-dispair/internal/state"
	"client-of-dispair/internal/ui"
	"time"
)

func CreateRoomHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	request := protocol.Request{
		Method:    "create",
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      nil,
	}
	err := client.SendRequest(&request)
	if err != nil {
		chat.Outputs <- "Error sending create room request: " + err.Error()
		return
	}
	response, err := client.ReceiveResponse()
	if err != nil {
		chat.Outputs <- "Error receiving create room response: " + err.Error()
		return
	}
	if response.Status != "ok" {
		chat.Outputs <- "Create room failed: " + response.Message
		return
	}
	roomID, ok := response.Data["room_id"].(string)
	if !ok {
		chat.Outputs <- "Invalid room ID in response"
		return
	}
	chat.Outputs <- "Room created successfully. Room ID: " + roomID
}

func JoinRoomHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	if len(arguments) < 1 {
		chat.Outputs <- "Usage: /join <room_id>"
		return
	}
	roomID := arguments[0]

	request := protocol.Request{
		Method:    "join",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]any{
			"user_id": state.UserID,
			"room_id": roomID,
		},
	}
	err := client.SendRequest(&request)
	if err != nil {
		chat.Outputs <- "Error sending join room request: " + err.Error()
		return
	}
	response, err := client.ReceiveResponse()
	if err != nil {
		chat.Outputs <- "Error receiving join room response: " + err.Error()
		return
	}
	if response.Status != "ok" {
		chat.Outputs <- "Join room failed: " + response.Message
		return
	}
	state.RoomID = roomID
	chat.Outputs <- "Joined room successfully. Room ID: " + roomID
}

func LeaveRoomHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	if state.RoomID == "" {
		chat.Outputs <- "You are not in a room."
		return
	}

	request := protocol.Request{
		Method:    "leave",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]any{
			"user_id": state.UserID,
			"room_id": state.RoomID,
		},
	}
	err := client.SendRequest(&request)
	if err != nil {
		chat.Outputs <- "Error sending leave room request: " + err.Error()
		return
	}
	response, err := client.ReceiveResponse()
	if err != nil {
		chat.Outputs <- "Error receiving leave room response: " + err.Error()
		return
	}
	if response.Status != "ok" {
		chat.Outputs <- "Leave room failed: " + response.Message
		return
	}
	chat.Outputs <- "Left room successfully. Room ID: " + state.RoomID
	state.RoomID = ""
}
