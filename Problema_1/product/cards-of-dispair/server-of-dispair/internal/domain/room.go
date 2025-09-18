package domain

import "server-of-dispair/internal/utils"
import "github.com/google/uuid"

type Room struct {
	ID       string
	UsersID  *utils.Set[string]
	Messages *utils.Map[string, chan string]
	Cards    *utils.Map[string, Card]
}

func NewRoom() *Room {
	return &Room{
		ID:       uuid.New().String(),
		UsersID:  utils.NewSet[string](),
		Messages: utils.NewMap[string, chan string](),
		Cards:    utils.NewMap[string, Card](),
	}
}
