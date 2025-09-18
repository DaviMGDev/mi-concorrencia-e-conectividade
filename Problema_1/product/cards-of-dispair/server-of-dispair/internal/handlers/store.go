package handlers

import (
	"math/rand"
	"server-of-dispair/internal/config"
	"server-of-dispair/internal/domain"
	"server-of-dispair/internal/protocol"
)

func HandleBuyPackage(server *protocol.Server, request *protocol.Request) {
	response := &protocol.Response{}
	pack := config.StoreService.GetPackage()
	response = protocol.NewResponse(request.From, request.Method, "ok", "purchase successfully", map[string]any{
		"rock":     pack.Cards[0].Stars,
		"paper":    pack.Cards[1].Stars,
		"scissors": pack.Cards[2].Stars,
	})

	pack = domain.CardPackage{
		Cards: [3]domain.Card{
			{Type: "rock", Stars: rand.Intn(5) + 1},
			{Type: "paper", Stars: rand.Intn(5) + 1},
			{Type: "scissors", Stars: rand.Intn(5) + 1},
		},
	}
	config.StoreService.AddPackage(pack)
	server.Responses <- response
}
