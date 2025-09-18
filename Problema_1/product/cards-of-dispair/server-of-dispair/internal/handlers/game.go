package handlers

import (
	"server-of-dispair/internal/config"
	"server-of-dispair/internal/protocol"
)

func HandlePlay(server *protocol.Server, request *protocol.Request) {
	response := protocol.NewResponse(request.From, request.Method, "", "", nil)
	userID, userExists := request.Data["user_id"].(string)
	cardType, cardExists := request.Data["card_type"].(string)
	cardStars, starsExists := request.Data["card_stars"].(string)
}
func HandleOpponentPlay(server *protocol.Server, request *protocol.Request) {}
