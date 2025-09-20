package application

import (
	"errors"
	"server-of-hope/internal/data"
	"server-of-hope/internal/domain"
	"server-of-hope/internal/utils"
)

// RoomServiceInterface descreve as operações para gerenciamento de salas.
//
// Métodos:
//   - CreateRoom: cria uma nova sala.
//   - JoinRoom: adiciona um usuário a uma sala.
//   - LeaveRoom: remove um usuário de uma sala.
type RoomServiceInterface interface {
	// CreateRoom cria uma nova sala e retorna seu ID.
	//
	// Retorno:
	//   - string: ID da sala criada.
	//   - erro caso não seja possível criar a sala.
	CreateRoom() (string, error)

	// JoinRoom adiciona um usuário a uma sala existente.
	//
	// Parâmetros:
	//   - roomID: identificador da sala.
	//   - userID: identificador do usuário.
	//
	// Retorno:
	//   - erro caso não seja possível adicionar o usuário.
	JoinRoom(roomID, userID string) error

	// LeaveRoom remove um usuário de uma sala existente.
	//
	// Parâmetros:
	//   - roomID: identificador da sala.
	//   - userID: identificador do usuário.
	//
	// Retorno:
	//   - erro caso não seja possível remover o usuário.
	LeaveRoom(roomID, userID string) error
}

// RoomService implementa a lógica de gerenciamento de salas.
//
// Campos:
//   - RoomRepo: repositório das salas.
type RoomService struct {
	RoomRepo data.RepositoryInterface[domain.Room]
}

// NewRoomService cria uma nova instância de RoomService.
//
// Parâmetros:
//   - roomRepo: repositório das salas.
//
// Retorno:
//   - ponteiro para RoomService.
func NewRoomService(roomRepo data.RepositoryInterface[domain.Room]) *RoomService {
	return &RoomService{RoomRepo: roomRepo}
}

// CreateRoom cria uma nova sala e retorna seu ID.
//
// Retorno:
//   - string: ID da sala criada.
//   - erro caso não seja possível criar a sala.
func (service *RoomService) CreateRoom() (string, error) {
	room := domain.NewRoom(utils.Count())
	err := service.RoomRepo.Create(room.ID, *room)
	if err != nil {
		return "", err
	}
	return room.ID, nil
}

// JoinRoom adiciona um usuário a uma sala existente e cria um canal de mensagens para ele.
//
// Parâmetros:
//   - roomID: identificador da sala.
//   - userID: identificador do usuário.
//
// Retorno:
//   - erro caso não seja possível adicionar o usuário.
func (service *RoomService) JoinRoom(roomID, userID string) error {
	room, err := service.RoomRepo.Read(roomID)
	if err != nil {
		return err // Sala não encontrada
	}

	if room.UserIDs.Contains(userID) {
		return nil // Usuário já está na sala
	}

	if room.UserIDs.Size() >= 2 {
		return errors.New("A sala está cheia")
	}

	room.UserIDs.Add(userID)
	room.Messages.Set(userID, make(chan string, 1))
	return service.RoomRepo.Update(roomID, room)
}

// LeaveRoom remove um usuário de uma sala e exclui seu canal de mensagens.
//
// Parâmetros:
//   - roomID: identificador da sala.
//   - userID: identificador do usuário.
//
// Retorno:
//   - erro caso não seja possível remover o usuário.
func (service *RoomService) LeaveRoom(roomID, userID string) error {
	room, err := service.RoomRepo.Read(roomID)
	if err != nil {
		return err
	}
	room.UserIDs.Remove(userID)
	room.Messages.Delete(userID)
	return service.RoomRepo.Update(roomID, room)
}
