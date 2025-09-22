package handlers

import (
	"client-of-hope/internal/api"
	"client-of-hope/internal/api/protocol"
	"client-of-hope/internal/state"
	"client-of-hope/internal/ui"
	"client-of-hope/internal/utils"
	"fmt"
	"strings"
)

func validatePlay(chat *ui.Chat, args []string) (string, int, bool) {
	if state.UserID == "" || state.RoomID == "" {
		chat.Outputs <- "You must be logged in and in a room to play."
		return "", 0, false
	}
	if len(args) != 1 {
		chat.Outputs <- "Usage: /play <card>"
		return "", 0, false
	}

	cardToPlay := strings.ToLower(args[0])
	stars, exists := state.Cards.Get(cardToPlay)
	if !exists {
		chat.Outputs <- fmt.Sprintf("You don't have a '%s' card.", cardToPlay)
		return "", 0, false
	}

	return cardToPlay, stars, true
}

func playCard(client *api.Client, chat *ui.Chat, cardToPlay string, stars int) bool {
	playRequest := protocol.Request{
		Method: "play",
		Data:   utils.Dict{"user_id": state.UserID, "room_id": state.RoomID, "card": cardToPlay, "stars": stars},
	}

	playResponse, err := client.DoRequest(playRequest)
	if err != nil {
		state.Log("Play card request failed: %v", err)
		chat.Outputs <- "Failed to play card."
		return false
	}
	if playResponse.Status != "ok" {
		message, _ := playResponse.Data["message"].(string)
		chat.Outputs <- message
		return false
	}

	chat.Outputs <- fmt.Sprintf("You played a %s card with %d stars.", cardToPlay, stars)
	state.PlayedCard = cardToPlay
	state.PlayedCardStar = stars
	return true
}

func determineWinner(chat *ui.Chat) {
	wins := 0
	if state.CardWins[state.PlayedCard] == state.OpponentCard {
		wins = 1
	} else if state.CardWins[state.OpponentCard] == state.PlayedCard {
		wins = -1
	} else if state.PlayedCardStar > state.OpponentCardStar {
		wins = 1
	} else if state.PlayedCardStar < state.OpponentCardStar {
		wins = -1
	}

	if wins > 0 {
		chat.Outputs <- "You win this round!"
	} else if wins < 0 {
		chat.Outputs <- "You lose this round!"
	} else {
		chat.Outputs <- "This round is a tie!"
	}
}

func resetRound() {
	state.PlayedCard = ""
	state.PlayedCardStar = 0
	state.OpponentCard = ""
	state.OpponentCardStar = 0
}
