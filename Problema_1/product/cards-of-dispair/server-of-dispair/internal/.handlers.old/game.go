package handlers

import (
	"server-of-dispair/internal/config"
	"server-of-dispair/internal/protocol"
)

func HandleReady(server *protocol.Server, request *protocol.Request) {
	playerID, playerIDOk := request.Data["player_id"].(string)
	gameID, gameIDOk := request.Data["game_id"].(string)

	if !playerIDOk || !gameIDOk {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data", nil)
		server.Responses <- response
		return
	}

	err := config.GameService.AddPlayer(gameID, playerID)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	response := protocol.NewResponse(request.From, request.Method, "success", "Player is ready", nil)
	server.Responses <- response
}

func HandleUnready(server *protocol.Server, request *protocol.Request) {
	playerID, playerIDOk := request.Data["player_id"].(string)
	gameID, gameIDOk := request.Data["game_id"].(string)

	if !playerIDOk || !gameIDOk {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data", nil)
		server.Responses <- response
		return
	}

	err := config.GameService.RemovePlayer(gameID, playerID)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	response := protocol.NewResponse(request.From, request.Method, "success", "Player is unready", nil)
	server.Responses <- response
}

func HandlePlay(server *protocol.Server, request *protocol.Request) {
	playerID, playerIDOk := request.Data["player_id"].(string)
	gameID, gameIDOk := request.Data["game_id"].(string)
	cardType, cardTypeOk := request.Data["card_type"].(string)
	stars, starsOk := request.Data["stars"].(float64) // JSON numbers are float64 by default

	if !playerIDOk || !gameIDOk || !cardTypeOk || !starsOk {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data", nil)
		server.Responses <- response
		return
	}

	err := config.GameService.Play(gameID, playerID, cardType, uint(stars))
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	response := protocol.NewResponse(request.From, request.Method, "success", "Card played", nil)
	server.Responses <- response
}

func HandleGetOpponentChoice(server *protocol.Server, request *protocol.Request) {
	playerID, playerIDOk := request.Data["player_id"].(string)
	gameID, gameIDOk := request.Data["game_id"].(string)

	if !playerIDOk || !gameIDOk {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data", nil)
		server.Responses <- response
		return
	}

	card, err := config.GameService.GetOpponentChoice(gameID, playerID)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	responseData := map[string]any{
		"type":  card.GetType(),
		"stars": card.GetStars(),
	}
	response := protocol.NewResponse(request.From, request.Method, "success", "Opponent choice retrieved", responseData)
	server.Responses <- response
}
