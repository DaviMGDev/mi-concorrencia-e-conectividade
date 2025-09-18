package handlers

import "server-of-dispair/internal/protocol"

func HandleRegister(server *protocol.Server, request *protocol.Request) {
	username, usernameOk := request.Data["username"].(string)
	password, passwordOk := request.Data["password"].(string)

	response := &protocol.Response{}
	if !usernameOk || !passwordOk {
	}
}
func HandleLogin(server *protocol.Server, request *protocol.Request) {}
