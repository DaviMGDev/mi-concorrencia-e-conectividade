package handlers

import (
	"client-of-dispair/internal/api"
	"client-of-dispair/internal/api/protocol"
	"client-of-dispair/internal/state"
	"client-of-dispair/internal/ui"
	"time"
)

func ChoiceHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	if len(arguments) < 1 {
		chat.Outputs <- "Usage: /choice <option>"
		return
	}
	request := protocol.Request{
		Method:    "choice",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]any{
			"user_id": state.UserID,
			"room_id": state.RoomID,
			"choice":  arguments[0],
		},
	}
	err := client.SendRequest(&request)
	if err != nil {
		chat.Outputs <- "Error sending choice request: " + err.Error()
		return
	}
	response, err := client.ReceiveResponse()
	if err != nil {
		chat.Outputs <- "Error receiving choice response: " + err.Error()
		return
	}
	if response.Status != "ok" {
		chat.Outputs <- "Choice failed: " + response.Message
		return
	}
	chat.Outputs <- "Choice submitted successfully."
}

func PlayHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	if len(arguments) < 1 {
		chat.Outputs <- "Usage: /play <move>"
		return
	}
	ChoiceHandler(client, chat, arguments)
	request := protocol.Request{
		Method:    "get_opponent_choice",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]any{
			"user_id": state.UserID,
			"room_id": state.RoomID,
		},
	}
	err := client.SendRequest(&request)
	if err != nil {
		chat.Outputs <- "Error sending get_opponent_choice request: " + err.Error()
		return
	}
	opponentChoice := ""
	ok := false
	for opponentChoice == "" {
		response, err := client.ReceiveResponse()
		if err != nil {
			chat.Outputs <- "Error receiving get_opponent_choice response: " + err.Error()
			return
		}
		if response.Status != "ok" {
			chat.Outputs <- "Get opponent choice failed: " + response.Message
			return
		}
		opponentChoice, ok = response.Data["opponent_choice"].(string)
		if !ok {
			chat.Outputs <- "Invalid opponent choice data"
			return
		}
	}
	chat.Outputs <- "You played: " + arguments[0] + ", Opponent played: " + opponentChoice + "."
}
