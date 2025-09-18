package services

import (
	"server-of-dispair/internal/domain"
	"server-of-dispair/internal/repositories"
)

type RoomService struct {
	RoomRepo *repositories.InMemoryRepository[domain.Room]
}

func (service *RoomService) CreateRoom() (*domain.Room, error) {
	room := domain.NewRoom()
	err := service.RoomRepo.Create(room.ID, *room)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (service *RoomService) JoinRoom(roomID, userID string) error {
	room, err := service.RoomRepo.Read(roomID)
	if err != nil {
		return err
	}
	room.UsersID.Set(userID, true)
	room.Messages.Set(userID, make(chan string, 100))
	return service.RoomRepo.Update(roomID, room)
}

func (service *RoomService) LeaveRoom(roomID, userID string) error {
	room, err := service.RoomRepo.Read(roomID)
	if err != nil {
		return err
	}
	room.UsersID.Delete(userID)
	room.Messages.Delete(userID)
	room.Cards.Delete(userID)
	return service.RoomRepo.Update(roomID, room)
}

func NewRoomService(roomRepo *repositories.InMemoryRepository[domain.Room]) *RoomService {
	return &RoomService{
		RoomRepo: roomRepo,
	}
}
