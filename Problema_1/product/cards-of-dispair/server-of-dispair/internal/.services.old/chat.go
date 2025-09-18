package services

import (
	"server-of-dispair/internal/entities"
	"server-of-dispair/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type ChatServiceInterface interface {
	SendMessage(msg *entities.Message) error
	GetMessages(roomID string, since time.Time) ([]*entities.Message, error)
}

type ChatService struct {
	MessageRepo *repositories.InMemoryRepository[*entities.Message]
}

func NewChatService(messageRepo *repositories.InMemoryRepository[*entities.Message]) *ChatService {
	return &ChatService{
		MessageRepo: messageRepo,
	}
}

func (s *ChatService) SendMessage(msg *entities.Message) error {
	// In a real application, you would generate a unique ID for the message.
	// For simplicity, we can use a UUID here.
	return s.MessageRepo.Create(uuid.New().String(), msg)
}

func (s *ChatService) GetMessages(roomID string, since time.Time) ([]*entities.Message, error) {
	messages, err := s.MessageRepo.List()
	if err != nil {
		return nil, err
	}

	var filteredMessages []*entities.Message
	for _, msg := range messages {
		if msg.RoomID == roomID {
			msgTimestamp, err := time.Parse(time.RFC3339, msg.TimeStamp)
			if err != nil {
				// Ignore messages with invalid timestamp format
				continue
			}
			if msgTimestamp.After(since) {
				filteredMessages = append(filteredMessages, msg)
			}
		}
	}

	return filteredMessages, nil
}