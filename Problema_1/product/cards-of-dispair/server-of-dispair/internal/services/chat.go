package services

import (
	"server-of-dispair/internal/domain"
	"server-of-dispair/internal/repositories"
)

type ChatService struct {
	UserRepo repositories.RepositoryInterface[domain.User]
	RoomRepo repositories.RepositoryInterface[domain.Room]
}

func (service *ChatService) SendMessage(roomID, userID, message string) error {
	room, err := service.RoomRepo.Read(roomID)
	if err != nil {
		return err
	}
	if !(&room).UsersID.Contains(userID) {
		return nil
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
	if !(&room).UsersID.Contains(userID) {
		return "", nil
	}
	messages, _ := room.Messages.Get(userID)
	message := <-messages
	return message, nil
}

func NewChatService(userRepo repositories.RepositoryInterface[domain.User], roomRepo repositories.RepositoryInterface[domain.Room]) *ChatService {
	return &ChatService{
		UserRepo: userRepo,
		RoomRepo: roomRepo,
	}
}
