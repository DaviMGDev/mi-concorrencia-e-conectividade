package handlers

import (
	"server-of-hope/internal/api"
	"server-of-hope/internal/api/protocol"
	"server-of-hope/internal/utils"
)

func HandlePing(server *api.Server, request protocol.Request) {
	responder := NewResponder(server, request)
	defer responder.Send()

	data := utils.Dict{"message": "pong"}
	responder.SetSuccess(data, "Ping received", "from", request.From)
}
