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
	defer responder.Send()

	userID, _ := request.Data["user_id"].(string)
	gameID, _ := request.Data["room_id"].(string) // In client, it's room_id
	cardType, _ := request.Data["card"].(string)
	cardStars, _ := request.Data["stars"].(float64)

	if userID == "" || gameID == "" || cardType == "" {
		responder.SetError("Invalid parameters", "Card play failed", "from", request.From)
		return
	}

	card := domain.Card{
		Type:  cardType,
		Stars: int(cardStars),
	}

	err := state.GameService.PlayCard(gameID, userID, card)
	if err != nil {
		responder.SetError(err.Error(), "Card play failed", "user_id", userID, "game_id", gameID, "error", err)
		return
	}

	data := utils.Dict{"message": "Card played successfully"}
	responder.SetSuccess(data, "Card played successfully", "user_id", userID, "game_id", gameID, "card", cardType, "stars", cardStars)
}

func HandleGetOpponentCard(server *api.Server, request protocol.Request) {
	responder := NewResponder(server, request)
	defer responder.Send()

	userID, _ := request.Data["user_id"].(string)
	gameID, _ := request.Data["room_id"].(string) // In client, it's room_id

	if userID == "" || gameID == "" {
		responder.SetError("Invalid parameters", "Get opponent card failed", "from", request.From)
		return
	}

	card, err := state.GameService.GetOpponentCard(gameID, userID)
	if err != nil {
		responder.SetError(err.Error(), "Get opponent card failed", "user_id", userID, "game_id", gameID, "error", err)
		return
	}

	data := utils.Dict{
		"opponent_card":      card.Type,
		"opponent_card_star": card.Stars,
	}
	responder.SetSuccess(data, "Card fetched successfully", "user_id", userID, "game_id", gameID)
}
