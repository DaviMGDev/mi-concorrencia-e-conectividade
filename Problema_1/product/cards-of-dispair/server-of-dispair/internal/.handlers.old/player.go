package handlers

import (
	"server-of-dispair/internal/config"
	"server-of-dispair/internal/protocol"
)

func HandleGetHealth(server *protocol.Server, request *protocol.Request) {
	playerID, playerIDOk := request.Data["player_id"].(string)

	if !playerIDOk {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data: missing player_id", nil)
		server.Responses <- response
		return
	}

	health, err := config.PlayerService.GetHealth(playerID)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	responseData := map[string]any{"health": health}
	response := protocol.NewResponse(request.From, request.Method, "success", "Player health retrieved successfully", responseData)
	server.Responses <- response
}