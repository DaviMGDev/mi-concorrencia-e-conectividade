package handlers

import (
	"client-of-hope/internal/api"
	"client-of-hope/internal/api/protocol"
	"client-of-hope/internal/state"
	"client-of-hope/internal/ui"
	"fmt"
	"strings"
)

func HandleCards(client *api.Client, chat *ui.Chat, args []string) {
	var cardList []string
	state.Cards.ForEach(func(card string, stars int) {
		cardList = append(cardList, fmt.Sprintf("%s (%d stars)", card, stars))
	})

	chat.Outputs <- "Your cards: " + strings.Join(cardList, ", ")
}

func HandlePlay(client *api.Client, chat *ui.Chat, args []string) {
	cardToPlay, stars, ok := validatePlay(chat, args)
	if !ok {
		return
	}

	if !playCard(client, chat, cardToPlay, stars) {
		return
	}
}

func HandleBuy(client *api.Client, chat *ui.Chat, args []string) {
	request := protocol.Request{
		Method: "buy",
		Data:   nil,
	}
	response, err := client.DoRequest(request)
	if err != nil {
		state.Log("Buy card request failed: %v", err)
		chat.Outputs <- "Failed to buy card."
		return
	}
	if response.Status != "ok" {
		message, _ := response.Data["message"].(string)
		chat.Outputs <- message
		return
	}
	cardPackage, _ := response.Data["package"].(map[string]any)
	rockStarsF, _ := cardPackage["rock"].(float64)
	paperStarsF, _ := cardPackage["paper"].(float64)
	scissorsStarsF, _ := cardPackage["scissors"].(float64)
	rockStars := int(rockStarsF)
	paperStars := int(paperStarsF)
	scissorsStars := int(scissorsStarsF)

	if rockStars > 0 {
		state.Cards.Set("rock", rockStars)
	}
	if paperStars > 0 {
		state.Cards.Set("paper", paperStars)
	}
	if scissorsStars > 0 {
		state.Cards.Set("scissors", scissorsStars)
	}

	chat.Outputs <- fmt.Sprintf("You bought a card package: rock (%d stars), paper (%d stars), scissors (%d stars).", rockStars, paperStars, scissorsStars)
}
