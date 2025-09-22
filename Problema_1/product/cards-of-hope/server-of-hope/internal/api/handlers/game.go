package handlers

import (
	"server-of-hope/internal/api"
	"server-of-hope/internal/api/protocol"
	"server-of-hope/internal/domain"
	"server-of-hope/internal/state"
	"server-of-hope/internal/utils"
)

func HandlePlayCard(server *api.Server, request protocol.Request) {
	responder := NewResponder(server, request)

	userID, _ := request.Data["user_id"].(string)
	gameID, _ := request.Data["room_id"].(string) // In client, it's room_id
	cardType, _ := request.Data["card"].(string)
	cardStars, _ := request.Data["stars"].(float64)

	if userID == "" || gameID == "" || cardType == "" {
		responder.SetError("Invalid parameters", "Card play failed", "from", request.From)
		responder.Send()
		return
	}

	card := domain.Card{
		Type:  cardType,
		Stars: int(cardStars),
	}

	err := state.GameService.PlayCard(gameID, userID, card)
	if err != nil {
		responder.SetError(err.Error(), "Card play failed", "user_id", userID, "game_id", gameID, "error", err)
		responder.Send()
		return
	}

	data := utils.Dict{"message": "Card played successfully"}
	responder.SetSuccess(data, "Card played successfully", "user_id", userID, "game_id", gameID, "card", cardType, "stars", cardStars)
	responder.Send()

	// Check if game is ready
	game, err := state.GameService.GetGame(gameID)
	if err != nil {
		state.Logger.Error("Failed to get game state after play", "game_id", gameID, "error", err)
		return
	}

	if game.Plays.Size() == 2 {
		// Notify both players
		playerIDs := game.Plays.Keys()
		card1, _ := game.Plays.Get(playerIDs[0])
		card2, _ := game.Plays.Get(playerIDs[1])

		notifyPlayer(server, playerIDs[0], card2)
		notifyPlayer(server, playerIDs[1], card1)

		// Reset round
		state.GameService.ResetRound(gameID)
	}
}

func notifyPlayer(server *api.Server, playerID string, opponentCard domain.Card) {
	playerAddress, ok := state.UserConnections.Get(playerID)
	if !ok {
		state.Logger.Warn("Could not find connection for user to notify", "user_id", playerID)
		return
	}

	response := protocol.Response{
		Method: "opponent_played",
		Status: "ok",
		Data: utils.Dict{
			"opponent_card":      opponentCard.Type,
			"opponent_card_star": opponentCard.Stars,
		},
		To: playerAddress,
	}
	server.Responses <- response
}
