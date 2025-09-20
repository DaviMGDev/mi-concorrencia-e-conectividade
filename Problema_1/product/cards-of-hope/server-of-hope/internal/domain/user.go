package domain

// User representa um usuário do sistema.
//
// Campos:
//   - ID: identificador único do usuário.
//   - Username: nome de usuário.
//   - Password: senha do usuário (omitida na serialização, se vazia).
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}
