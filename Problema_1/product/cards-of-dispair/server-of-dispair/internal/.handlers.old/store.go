package handlers

import (
	"server-of-dispair/internal/config"
	"server-of-dispair/internal/entities"
	"server-of-dispair/internal/protocol"
)

func HandleBuyPackage(server *protocol.Server, request *protocol.Request) {
	userID, ok := request.Data["user_id"].(string)
	if !ok {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data: missing user_id", nil)
		server.Responses <- response
		return
	}

	packageBought, err := config.StoreService.BuyPackage()
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", "No packages available in store", nil)
		server.Responses <- response
		return
	}
	// Restock a new package for the next player
	go config.StoreService.RestockPackage()

	// Convert from [3]CardInterface to [3]Card
	var newHand [3]entities.Card
	cards := packageBought.GetCards()
	for i, c := range cards {
		// This is a bit of a hack, assuming the interface holds a Card struct.
		// A more robust solution would be a type switch or a more specific interface.
		newHand[i] = c.(entities.Card)
	}

	err = config.PlayerService.ChangeHand(userID, newHand)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", "Could not add cards to player's hand", nil)
		server.Responses <- response
		return
	}

	responseData := map[string]any{"package": newHand}
	response := protocol.NewResponse(request.From, request.Method, "success", "Package bought successfully", responseData)
	server.Responses <- response
}
