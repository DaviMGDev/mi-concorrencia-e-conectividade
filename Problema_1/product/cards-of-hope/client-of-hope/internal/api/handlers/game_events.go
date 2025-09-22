package handlers

import (
	"client-of-hope/internal/api"
	"client-of-hope/internal/api/protocol"
	"client-of-hope/internal/state"
	"client-of-hope/internal/ui"
	"fmt"
)

func HandleOpponentPlayed(client *api.Client, chat *ui.Chat, response protocol.Response) {
	if response.Status != "ok" {
		// O servidor pode enviar um erro se algo der errado do lado dele
		message, _ := response.Data["message"].(string)
		chat.Outputs <- fmt.Sprintf("Server error: %s", message)
		return
	}

	opponentCard, cardOk := response.Data["opponent_card"].(string)
	opponentCardStar, starOk := response.Data["opponent_card_star"].(float64)

	if !cardOk || !starOk {
		chat.Outputs <- "Invalid opponent card data from server."
		return
	}

	state.OpponentCard = opponentCard
	state.OpponentCardStar = int(opponentCardStar)
	chat.Outputs <- fmt.Sprintf("Opponent played a %s card with %d stars.", state.OpponentCard, state.OpponentCardStar)

	// Agora que temos as duas cartas, podemos determinar o vencedor
	determineWinner(chat)
	resetRound()
}
