package services

import (
	"errors"
	"server-of-dispair/internal/entities"
	"server-of-dispair/internal/repositories"

	"github.com/google/uuid"
)

type AuthServiceInterface interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}

type AuthService struct {
	UserRepo   *repositories.InMemoryRepository[*entities.User]
	PlayerRepo *repositories.InMemoryRepository[*entities.Player]
}

func NewAuthService(userRepo *repositories.InMemoryRepository[*entities.User], playerRepo *repositories.InMemoryRepository[*entities.Player]) *AuthService {
	return &AuthService{
		UserRepo:   userRepo,
		PlayerRepo: playerRepo,
	}
}

func (a *AuthService) Register(username, password string) error {
	list, err := a.UserRepo.List()
	if err != nil {
		return err
	}
	for _, u := range list {
		if u.Username == username {
			return errors.New("user with this username already exists")
		}
	}
	user, err := entities.NewUser(uuid.New().String(), username, password)
	if err != nil {
		return err
	}
	player := entities.NewPlayer(user.ID)

	if err := a.UserRepo.Create(user.ID, user); err != nil {
		return err
	}
	if err := a.PlayerRepo.Create(user.ID, player); err != nil {
		// Rollback user creation if player creation fails
		a.UserRepo.Delete(user.ID)
		return err
	}
	return nil
}

func (a *AuthService) Login(username, password string) (string, error) {
	users, err := a.UserRepo.List()
	if err != nil {
		return "", err
	}

	for _, user := range users {
		if user.Username == username {
			if user.CheckPassword(password) {
				return user.ID, nil
			}
			return "", errors.New("invalid password")
		}
	}

	return "", errors.New("user not found")
}
