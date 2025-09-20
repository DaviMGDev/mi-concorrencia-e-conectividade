// Pacote state armazena informações globais do usuário logado na aplicação.
package state

// Username armazena o nome do usuário atualmente logado.
// UserID armazena o identificador único do usuário.
// RoomID armazena o identificador da sala em que o usuário está.
var (
	Username string
	UserID   string
	RoomID   string
)
