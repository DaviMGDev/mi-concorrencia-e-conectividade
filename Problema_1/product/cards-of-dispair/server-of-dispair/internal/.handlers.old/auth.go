package handlers

import (
	"server-of-dispair/internal/config"
	"server-of-dispair/internal/protocol"
)

func HandleRegister(server *protocol.Server, request *protocol.Request) {
	username, usernameOk := request.Data["username"].(string)
	password, passwordOk := request.Data["password"].(string)

	if !usernameOk || !passwordOk {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data: missing username or password", nil)
		server.Responses <- response
		return
	}

	err := config.AuthService.Register(username, password)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	response := protocol.NewResponse(request.From, request.Method, "success", "User registered successfully", nil)
	server.Responses <- response
}

func HandleLogin(server *protocol.Server, request *protocol.Request) {
	username, usernameOk := request.Data["username"].(string)
	password, passwordOk := request.Data["password"].(string)

	if !usernameOk || !passwordOk {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data: missing username or password", nil)
		server.Responses <- response
		return
	}

	userID, err := config.AuthService.Login(username, password)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	responseData := map[string]any{"user_id": userID}
	response := protocol.NewResponse(request.From, request.Method, "success", "User logged in successfully", responseData)
	server.Responses <- response
}