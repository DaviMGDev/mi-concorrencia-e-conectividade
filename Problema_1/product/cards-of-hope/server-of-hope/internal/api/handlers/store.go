package handlers

import (
	"math/rand"
	"server-of-hope/internal/api"
	"server-of-hope/internal/api/protocol"
	"server-of-hope/internal/domain"
	"server-of-hope/internal/state"
	"server-of-hope/internal/utils"
)

func HandleBuyPackage(server *api.Server, request protocol.Request) {
	responder := NewResponder(server, request)
	defer responder.Send()

	r, p, s := rand.Intn(5)+1, rand.Intn(5)+1, rand.Intn(5)+1
	pack := domain.CardPackage{
		domain.Card{Type: "rock", Stars: r},
		domain.Card{Type: "paper", Stars: p},
		domain.Card{Type: "scissors", Stars: s},
	}
	state.StoreService.AddPackage(pack)

	pack = state.StoreService.GetPackage()
	if len(pack) < 3 {
		responder.SetError("Package is empty", "Buy package failed", "from", request.From)
		return
	}

	rock, paper, scissors := pack[0], pack[1], pack[2]
	data := utils.Dict{
		"package": utils.Dict{
			"rock":     rock.Stars,
			"paper":    paper.Stars,
			"scissors": scissors.Stars,
		},
	}
	responder.SetSuccess(data, "Package bought successfully", "from", request.From)
}
