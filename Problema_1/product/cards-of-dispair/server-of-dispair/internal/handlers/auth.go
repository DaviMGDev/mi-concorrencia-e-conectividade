package handlers

import (
	"server-of-dispair/internal/protocol"

	"github.com/google/uuid"
)

func HandleRegister(server *protocol.Server, request *protocol.Request) {
	username, usernameOk := request.Data["username"].(string)
	password, passwordOk := request.Data["password"].(string)

	response := protocol.NewResponse(request.From, request.Method, "error", "invalid request data", nil)
	if !usernameOk || !passwordOk {
		server.Send(response)
		return
	}
	err := server.AuthService.Register(protocol.User{
		ID:       uuid.New().String(),
		Username: username,
		Password: password,
	})
	if err != nil {
		response.Message = "registration failed"
		server.Send(response)
		return
	}
	response.Status = "success"
	response.Message = "registration successful"
	server.Send(response)
}

func HandleLogin(server *protocol.Server, request *protocol.Request) {
	username, usernameOk := request.Data["username"].(string)
	password, passwordOk := request.Data["password"].(string)

	response := protocol.NewResponse(request.From, request.Method, "error", "invalid request data", nil)
	if !usernameOk || !passwordOk {
		server.Send(response)
		return
	}
	user, err := server.AuthService.Login(username, password)
	if err != nil {
		response.Message = "login failed"
		server.Send(response)
		return
	}
	if user == nil {
		response.Message = "invalid username or password"
		server.Send(response)
		return
	}
	response.Status = "success"
	response.Message = "login successful"
	response.Data = map[string]any{
		"user_id": user.ID,
	}
	server.Send(response)
}
