package handlers

import (
	"math/rand"
	"server-of-dispair/internal/config"
	"server-of-dispair/internal/domain"
	"server-of-dispair/internal/protocol"
)

func HandleBuyPackage(server *protocol.Server, request *protocol.Request) {
	pack := config.StoreService.GetPackage()

	response := &protocol.Response{}
	pack = domain.CardPackage{
		Cards: [3]domain.Card{
			{Type: "rock", Stars: rand.Intn(5) + 1},
			{Type: "paper", Stars: rand.Intn(5) + 1},
			{Type: "scissors", Stars: rand.Intn(5) + 1},
		},
	}
	server.Responses <- response
}
