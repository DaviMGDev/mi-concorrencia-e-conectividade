package services

import (
	"server-of-dispair/internal/entities"
	"server-of-dispair/internal/repositories"

	"github.com/google/uuid"
)

type RoomServiceInterface interface {
	CreateRoom() (string, error)
	JoinRoom(roomID, playerID string) error
	LeaveRoom(roomID, playerID string) error
}

type RoomService struct {
	RoomRepo *repositories.InMemoryRepository[*entities.Room]
	UserRepo *repositories.InMemoryRepository[*entities.User]
}

func NewRoomService(roomRepo *repositories.InMemoryRepository[*entities.Room], userRepo *repositories.InMemoryRepository[*entities.User]) *RoomService {
	return &RoomService{
		RoomRepo: roomRepo,
		UserRepo: userRepo,
	}
}

func (r *RoomService) CreateRoom() (string, error) {
	roomID := uuid.New().String()
	room := entities.NewRoom(roomID)
	err := r.RoomRepo.Create(roomID, room)
	if err != nil {
		return "", err
	}
	return roomID, nil
}

func (r *RoomService) JoinRoom(roomID, playerID string) error {
	room, err := r.RoomRepo.Read(roomID)
	if err != nil {
		return err
	}

	user, err := r.UserRepo.Read(playerID)
	if err != nil {
		return err
	}

	err = room.AddMember(*user)
	if err != nil {
		return err
	}
	return r.RoomRepo.Update(roomID, room)
}

func (r *RoomService) LeaveRoom(roomID, playerID string) error {
	room, err := r.RoomRepo.Read(roomID)
	if err != nil {
		return err
	}
	err = room.RemoveMember(playerID)
	if err != nil {
		return err
	}
	return r.RoomRepo.Update(roomID, room)
}