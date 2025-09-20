package domain

import "server-of-hope/internal/utils"

// Room representa uma sala de jogo.
//
// Campos:
//   - ID: identificador único da sala.
//   - UserIDs: IDs dos usuários presentes na sala.
//   - Messages: canais de mensagens para cada usuário.
type Room struct {
	ID       string                          `json:"id"`
	UserIDs  *utils.Set[string]              `json:"user_ids"`
	Messages *utils.Map[string, chan string] `json:"-"`
}

// NewRoom cria uma nova sala com o ID informado.
//
// Parâmetros:
//   - id: identificador da sala.
//
// Retorno:
//   - ponteiro para Room.
func NewRoom(id string) *Room {
	return &Room{
		ID:       id,
		UserIDs:  utils.NewSet[string](),
		Messages: utils.NewMap[string, chan string](),
	}
}
