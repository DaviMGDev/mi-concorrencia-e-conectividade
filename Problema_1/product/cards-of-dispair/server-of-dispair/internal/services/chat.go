package services

import (
	"server-of-dispair/internal/domain"
	"server-of-dispair/internal/repositories"
)

type ChatService struct {
	UserRepo *repositories.InMemoryRepository[domain.User]
	RoomRepo *repositories.InMemoryRepository[domain.Room]
}

func (service *ChatService) SendMessage(roomID, userID, message string) error {
	room, err := service.RoomRepo.Read(roomID)
	if err != nil {
		return err
	}
	if !(&room).UsersID.Has(userID) {
		return nil // User not in room
	}
	for _, key := range room.Messages.Keys() {
		if key == userID {
			continue
		}
		messages, _ := room.Messages.Get(key)
		messages <- message
	}
	return nil
}

func (service *ChatService) GetMessage(roomID, userID string) (string, error) {
	room, err := service.RoomRepo.Read(roomID)
	if err != nil {
		return "", err
	}
	if !(&room).UsersID.Has(userID) {
		return "", nil // User not in room
	}
	messages, _ := room.Messages.Get(userID)
	message := <-messages
	return message, nil
}

func NewChatService(userRepo *repositories.InMemoryRepository[domain.User], roomRepo *repositories.InMemoryRepository[domain.Room]) *ChatService {
	return &ChatService{
		UserRepo: userRepo,
		RoomRepo: roomRepo,
	}
}
