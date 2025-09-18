package handlers

import "server-of-dispair/internal/protocol"

func HandlePing(server *protocol.Server, request *protocol.Request) {
	response := protocol.NewResponse(request.From, request.Method, "ok", "pong", nil)
	server.Responses <- response
}
