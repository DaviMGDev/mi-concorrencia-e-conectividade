package services

import (
	"server-of-dispair/internal/domain"
	"server-of-dispair/internal/repositories"
)

type AuthService struct {
	UserRepo *repositories.InMemoryRepository[domain.User]
}

func (service *AuthService) Register(user domain.User) error {
	return service.UserRepo.Create(user.ID, user)
}

func (service *AuthService) Login(username, password string) (*domain.User, error) {
	users, err := service.UserRepo.List()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if user.Username == username && user.Password == password {
			return &user, nil
		}
	}
	return nil, nil
}

func NewAuthService(userRepo *repositories.InMemoryRepository[domain.User]) *AuthService {
	return &AuthService{
		UserRepo: userRepo,
	}
}
