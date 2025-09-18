package handlers

import (
	"client-of-dispair/internal/api"
	"client-of-dispair/internal/api/protocol"
	"client-of-dispair/internal/state"
	"client-of-dispair/internal/ui"
	"time"
)

func RegisterHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	if len(arguments) < 2 {
		chat.Outputs <- "Usage: /register <username> <password>"
		return
	}
	username := arguments[0]
	password := arguments[1]

	request := protocol.Request{
		Method:    "register",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]any{
			"username": username,
			"password": password,
		},
	}
	err := client.SendRequest(&request)
	if err != nil {
		chat.Outputs <- "Error sending register request: " + err.Error()
		return
	}
	response, err := client.ReceiveResponse()
	if err != nil {
		chat.Outputs <- "Error receiving register response: " + err.Error()
		return
	}
	if response.Status != "ok" {
		chat.Outputs <- "Register failed: " + response.Message
		return
	}
	chat.Outputs <- "Register successful. You can now log in."
}

func LoginHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	if len(arguments) < 2 {
		chat.Outputs <- "Usage: /login <username> <password>"
		return
	}
	username := arguments[0]
	password := arguments[1]

	request := protocol.Request{
		Method:    "login",
		Timestamp: time.Now().Format(time.RFC3339),
		Data: map[string]any{
			"username": username,
			"password": password,
		},
	}
	err := client.SendRequest(&request)
	if err != nil {
		chat.Outputs <- "Error sending login request: " + err.Error()
		return
	}
	response, err := client.ReceiveResponse()
	if err != nil {
		chat.Outputs <- "Error receiving login response: " + err.Error()
		return
	}
	if response.Status != "ok" {
		chat.Outputs <- "Login failed: " + response.Message
		return
	}
	chat.Outputs <- "Login successful. Welcome, " + username + "!"
	state.UserID = response.Data["user_id"].(string)
}

func LogoutHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	state.UserID = ""
	chat.Outputs <- "Logged out successfully."
}
