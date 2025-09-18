package services

import (
	"server-of-dispair/internal/domain"
	"server-of-dispair/internal/repositories"
)

type GameService struct {
	UserRepo *repositories.InMemoryRepository[domain.User]
	RoomRepo *repositories.InMemoryRepository[domain.Room]
}

func (service *GameService) AddCardToRoom(roomID, userID string, card domain.Card) error {
	room, err := service.RoomRepo.Read(roomID)
	if err != nil {
		return err
	}
	if !(&room).UsersID.Has(userID) {
		return nil // User not in room
	}
	room.Cards.Set(userID, card)
	return service.RoomRepo.Update(roomID, room)
}

func (service *GameService) GetCardsInRoom(roomID string) (map[string]domain.Card, error) {
	room, err := service.RoomRepo.Read(roomID)
	if err != nil {
		return nil, err
	}
	cards := make(map[string]domain.Card)
	for _, key := range room.Cards.Keys() {
		card, _ := room.Cards.Get(key)
		cards[key] = card
	}
	return cards, nil
}

func NewGameService(userRepo *repositories.InMemoryRepository[domain.User], roomRepo *repositories.InMemoryRepository[domain.Room]) *GameService {
	return &GameService{
		UserRepo: userRepo,
		RoomRepo: roomRepo,
	}
}
