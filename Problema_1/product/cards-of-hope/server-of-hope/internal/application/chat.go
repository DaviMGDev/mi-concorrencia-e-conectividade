package application

import (
	"server-of-hope/internal/data"
	"server-of-hope/internal/domain"
)

// ChatServiceInterface descreve as operações para envio e recebimento de mensagens em salas de chat.
//
// Métodos:
//   - SendMessage: envia mensagem para todos da sala, exceto o remetente.
//   - ReceiveMessage: recebe mensagem para um usuário específico.
type ChatServiceInterface interface {
	// SendMessage envia uma mensagem para todos os usuários da sala, exceto o remetente.
	//
	// Parâmetros:
	//   - roomID: identificador da sala.
	//   - userID: identificador do usuário remetente.
	//   - message: mensagem a ser enviada.
	//
	// Retorno:
	//   - erro caso não seja possível enviar a mensagem.
	SendMessage(roomID, userID, message string) error

	// ReceiveMessage recebe uma mensagem para um usuário específico na sala.
	//
	// Parâmetros:
	//   - roomID: identificador da sala.
	//   - userID: identificador do usuário destinatário.
	//
	// Retorno:
	//   - string: mensagem recebida.
	//   - erro caso não seja possível receber a mensagem.
	ReceiveMessage(roomID, userID string) (string, error)
}

// ChatService implementa a lógica de chat entre usuários em salas.
//
// Campos:
//   - RoomRepo: repositório das salas.
//   - UserRepo: repositório dos usuários.
type ChatService struct {
	RoomRepo data.RepositoryInterface[domain.Room]
	UserRepo data.RepositoryInterface[domain.User]
}

// NewChatService cria uma nova instância de ChatService.
//
// Parâmetros:
//   - roomRepo: repositório das salas.
//   - userRepo: repositório dos usuários.
//
// Retorno:
//   - ponteiro para ChatService.
func NewChatService(roomRepo data.RepositoryInterface[domain.Room], userRepo data.RepositoryInterface[domain.User]) *ChatService {
	return &ChatService{RoomRepo: roomRepo, UserRepo: userRepo}
}

// SendMessage envia uma mensagem para todos os usuários da sala, exceto o remetente.
//
// Parâmetros:
//   - roomID: identificador da sala.
//   - userID: identificador do usuário remetente.
//   - message: mensagem a ser enviada.
//
// Retorno:
//   - erro caso não seja possível enviar a mensagem.
func (service *ChatService) SendMessage(roomID, userID, message string) error {
	room, err := service.RoomRepo.Read(roomID)
	if err != nil {
		return err // Sala não encontrada
	}
	if !room.UserIDs.Contains(userID) {
		return err // Usuário não está na sala
	}
	room.Messages.ForEach(func(id string, messages chan string) {
		if id != userID {
			messages <- message
		}
	})
	return nil
}

// ReceiveMessage recebe uma mensagem para um usuário específico na sala, se disponível.
//
// Parâmetros:
//   - roomID: identificador da sala.
//   - userID: identificador do usuário destinatário.
//
// Retorno:
//   - string: mensagem recebida.
//   - erro caso não seja possível receber a mensagem.
func (service *ChatService) ReceiveMessage(roomID, userID string) (string, error) {
	room, err := service.RoomRepo.Read(roomID)
	if err != nil {
		return "", err // Sala não encontrada
	}
	if !room.UserIDs.Contains(userID) {
		return "", err // Usuário não está na sala
	}
	messages, exists := room.Messages.Get(userID)
	if !exists {
		return "", err // Canal de mensagens não existe para o usuário
	}
	select {
	case message := <-messages:
		return message, nil
	default:
		return "", nil // Nenhuma mensagem disponível
	}
}
