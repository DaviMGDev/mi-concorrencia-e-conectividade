package services

import (
	"server-of-dispair/internal/entities"
	"server-of-dispair/internal/repositories"
)

type PlayerServiceInterface interface {
	GetHealth(playerID string) (uint, error)
	TakeDamage(playerID string, damage uint) error
	ChangeHand(playerID string, newHand [3]entities.Card) error
	Read(playerID string) (*entities.Player, error)
}

type PlayerService struct {
	PlayerRepo *repositories.InMemoryRepository[*entities.Player]
}

func NewPlayerService(playerRepo *repositories.InMemoryRepository[*entities.Player]) *PlayerService {
	return &PlayerService{
		PlayerRepo: playerRepo,
	}
}

func (p *PlayerService) Read(playerID string) (*entities.Player, error) {
	return p.PlayerRepo.Read(playerID)
}

func (p *PlayerService) GetHealth(playerID string) (uint, error) {
	player, err := p.PlayerRepo.Read(playerID)
	if err != nil {
		return 0, err
	}
	return player.GetHP(), nil
}

func (p *PlayerService) TakeDamage(playerID string, damage uint) error {
	player, err := p.PlayerRepo.Read(playerID)
	if err != nil {
		return err
	}
	player.TakeDamage(damage)
	return p.PlayerRepo.Update(playerID, player)
}

func (p *PlayerService) ChangeHand(playerID string, newHand [3]entities.Card) error {
	player, err := p.PlayerRepo.Read(playerID)
	if err != nil {
		return err
	}
	player.ChangeHand(newHand)
	return p.PlayerRepo.Update(playerID, player)
}
