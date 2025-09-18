package services

import (
	"errors"
	"server-of-dispair/internal/entities"
	"server-of-dispair/internal/repositories"
	
	"sync"
)

type GameServiceInterface interface {
	AddPlayer(gameID, playerID string) error
	RemovePlayer(gameID, playerID string) error
	Play(gameID, playerID, cardType string, stars uint) error
	GetOpponentChoice(gameID, playerID string) (entities.Card, error)
}

type GameService struct {
	PlayerRepo *repositories.InMemoryRepository[*entities.Player]
	GameRepo   *repositories.InMemoryRepository[*entities.Game]

	// A map to notify waiting players about opponent moves.
	// Key: gameID + "_" + playerID
	// Value: a channel to send the opponent's card to.
	waiters map[string]chan entities.Card
	mutex   sync.Mutex
}

func NewGameService(gameRepo *repositories.InMemoryRepository[*entities.Game], playerRepo *repositories.InMemoryRepository[*entities.Player]) *GameService {
	return &GameService{
		GameRepo:   gameRepo,
		PlayerRepo: playerRepo,
		waiters:    make(map[string]chan entities.Card),
	}
}

func (s *GameService) AddPlayer(gameID, playerID string) error {
	game, err := s.GameRepo.Read(gameID)
	if err != nil {
		return err
	}
	return game.AddPlayer(playerID)
}

func (s *GameService) RemovePlayer(gameID, playerID string) error {
	// Implementation for removing a player from a game would go here.
	// For now, we can leave it empty as it is not critical for the game flow.
	return nil
}

func (s *GameService) Play(gameID, playerID, cardType string, stars uint) error {
	game, err := s.GameRepo.Read(gameID)
	if err != nil {
		return err
	}

	err = game.MakeChoice(playerID, cardType, stars)
	if err != nil {
		return err
	}

	// Notify the opponent if they are waiting.
	opponentID := ""
	for id := range game.Choice {
		if id != playerID {
			opponentID = id
			break
		}
	}

	if opponentID != "" {
		waiterKey := gameID + "_" + opponentID
		s.mutex.Lock()
		waiterChan, ok := s.waiters[waiterKey]
		if ok {
			// Send the card to the waiting opponent
			waiterChan <- game.Choice[playerID]
			// Remove the waiter from the map
			delete(s.waiters, waiterKey)
		}
		s.mutex.Unlock()
	}

	return nil
}

func (s *GameService) GetOpponentChoice(gameID, playerID string) (entities.Card, error) {
	game, err := s.GameRepo.Read(gameID)
	if err != nil {
		return entities.Card{}, err
	}

	opponentID := ""
	for id := range game.Choice {
		if id != playerID {
			opponentID = id
			break
		}
	}

	if opponentID == "" {
		return entities.Card{}, errors.New("no opponent found in the game")
	}

	// Check if the opponent has already made a choice
	if opponentCard, ok := game.Choice[opponentID]; ok && opponentCard.Type != "" {
		return opponentCard, nil
	}

	// If not, wait for the opponent's choice
	waiterKey := gameID + "_" + playerID
	waiterChan := make(chan entities.Card, 1)

	s.mutex.Lock()
	s.waiters[waiterKey] = waiterChan
	s.mutex.Unlock()

	// Wait for the card to be sent on the channel
	opponentCard := <-waiterChan

	return opponentCard, nil
}

// This is a helper function that might be useful for the game logic.
func (s *GameService) getOpponentID(game *entities.Game, playerID string) (string, error) {
	for id := range game.Choice {
		if id != playerID {
			return id, nil
		}
	}
	return "", errors.New("opponent not found")
}