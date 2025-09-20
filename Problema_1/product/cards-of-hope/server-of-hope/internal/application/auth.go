package application

import (
	"errors"
	"server-of-hope/internal/data"
	"server-of-hope/internal/domain"
)

// AuthServiceInterface descreve as operações de autenticação de usuários.
//
// Métodos:
//   - Register: registra um novo usuário.
//   - Login: autentica um usuário e retorna seu ID.
type AuthServiceInterface interface {
	// Register registra um novo usuário com nome de usuário e senha.
	//
	// Parâmetros:
	//   - username: nome de usuário.
	//   - password: senha do usuário.
	//
	// Retorno:
	//   - erro caso o usuário já exista ou haja falha no cadastro.
	Register(username, password string) error

	// Login autentica um usuário e retorna seu ID se as credenciais forem válidas.
	//
	// Parâmetros:
	//   - username: nome de usuário.
	//   - password: senha do usuário.
	//
	// Retorno:
	//   - string: ID do usuário autenticado.
	//   - erro caso as credenciais estejam incorretas ou usuário não exista.
	Login(username, password string) (string, error)
}

// AuthService implementa a lógica de autenticação utilizando um repositório de usuários.
//
// Campos:
//   - UserRepo: repositório responsável pelo armazenamento dos usuários.
type AuthService struct {
	UserRepo data.RepositoryInterface[domain.User]
}

// NewAuthService cria uma nova instância de AuthService.
//
// Parâmetros:
//   - userRepo: repositório de usuários.
//
// Retorno:
//   - ponteiro para AuthService.
func NewAuthService(userRepo data.RepositoryInterface[domain.User]) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

// Register registra um novo usuário se o nome de usuário ainda não existir.
//
// Parâmetros:
//   - username: nome de usuário.
//   - password: senha do usuário.
//
// Retorno:
//   - erro caso o usuário já exista ou haja falha no cadastro.
func (service *AuthService) Register(username, password string) error {
	_, err := service.UserRepo.Read(username)
	if err == nil {
		return err // Usuário já existe
	}
	user := domain.User{Username: username, Password: password}
	user.ID = username
	return service.UserRepo.Create(username, user)
}

// Login autentica um usuário e retorna seu ID se as credenciais estiverem corretas.
//
// Parâmetros:
//   - username: nome de usuário.
//   - password: senha do usuário.
//
// Retorno:
//   - string: ID do usuário autenticado.
//   - erro caso as credenciais estejam incorretas ou usuário não exista.
func (service *AuthService) Login(username, password string) (string, error) {
	user, err := service.UserRepo.Read(username)
	if err != nil {
		return "", err // Usuário não encontrado
	}
	if user.Password != password {
		return "", errors.New("senha inválida")
	}
	return user.ID, nil
}
