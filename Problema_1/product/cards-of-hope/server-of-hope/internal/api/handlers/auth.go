package handlers

import (
	"server-of-hope/internal/api"
	"server-of-hope/internal/api/protocol"
	"server-of-hope/internal/state"
	"server-of-hope/internal/utils"
)

func HandleRegisterUser(server *api.Server, request protocol.Request) {
	responder := NewResponder(server, request)
	defer responder.Send()

	username, usernameOk := request.Data["username"].(string)
	password, passwordOk := request.Data["password"].(string)

	if !usernameOk || !passwordOk {
		responder.SetError("Invalid username or password", "User registration failed", "from", request.From)
		return
	}

	err := state.AuthService.Register(username, password)
	if err != nil {
		responder.SetError("User already exists", "User registration failed", "username", username, "error", err)
		return
	}

	data := utils.Dict{"message": "User registered successfully"}
	responder.SetSuccess(data, "User registered successfully", "username", username)
}

func HandleLoginUser(server *api.Server, request protocol.Request) {
	responder := NewResponder(server, request)
	defer responder.Send()

	username, usernameOk := request.Data["username"].(string)
	password, passwordOk := request.Data["password"].(string)

	if !usernameOk || !passwordOk {
		responder.SetError("Invalid username or password", "User login failed", "from", request.From)
		return
	}

	userId, err := state.AuthService.Login(username, password)
	if err != nil {
		responder.SetError("Invalid username or password", "User login failed", "username", username, "error", err)
		return
	}

	client, exists := server.Clients.Get(request.From)
	if exists {
		client.UserID = userId
		state.UserConnections.Set(userId, request.From)
	}

	data := utils.Dict{
		"message": "User logged in successfully",
		"user_id": userId,
	}
	responder.SetSuccess(data, "User logged in successfully", "username", username, "userId", userId)
}
