package application

import (
	"errors"
	"server-of-hope/internal/data"
	"server-of-hope/internal/domain"
	"server-of-hope/internal/utils"
)

// GameServiceInterface descreve as operações para manipulação da lógica do jogo.
type GameServiceInterface interface {
	PlayCard(gameID string, playerID string, card domain.Card) error
	GetGame(gameID string) (domain.Game, error)
	ResetRound(gameID string) error
}

// GameService implementa a lógica do jogo, incluindo jogadas e controle de estado.
type GameService struct {
	gameRepo data.RepositoryInterface[domain.Game]
	userRepo data.RepositoryInterface[domain.User]
	roomRepo data.RepositoryInterface[domain.Room]
}

// NewGameService cria uma nova instância de GameService.
func NewGameService(
	gameRepo data.RepositoryInterface[domain.Game],
	userRepo data.RepositoryInterface[domain.User],
	roomRepo data.RepositoryInterface[domain.Room],
) *GameService {
	return &GameService{
		gameRepo: gameRepo,
		userRepo: userRepo,
		roomRepo: roomRepo,
	}
}

// getOrCreateGame busca uma partida existente ou cria uma nova se não existir.
func (s *GameService) getOrCreateGame(gameID string) (domain.Game, error) {
	game, err := s.gameRepo.Read(gameID)
	if err != nil {
		if err.Error() == "item not found" {
			game = domain.Game{
				ID:             gameID,
				Plays:          utils.NewMap[string, domain.Card](),
				ResultsSeenBy:  utils.NewSet[string](),
				FailedAttempts: utils.NewMap[string, int](),
			}
			if err := s.gameRepo.Create(gameID, game); err != nil {
				return domain.Game{}, err
			}
		} else {
			return domain.Game{}, err
		}
	} else {
		if game.Plays == nil {
			game.Plays = utils.NewMap[string, domain.Card]()
		}
		if game.ResultsSeenBy == nil {
			game.ResultsSeenBy = utils.NewSet[string]()
		}
		if game.FailedAttempts == nil {
			game.FailedAttempts = utils.NewMap[string, int]()
		}
	}
	return game, nil
}

func (s *GameService) GetGame(gameID string) (domain.Game, error) {
	return s.gameRepo.Read(gameID)
}

// PlayCard registra a jogada de um jogador em uma partida, validando a carta e o estado do jogo.
func (s *GameService) PlayCard(gameID string, playerID string, card domain.Card) error {
	if card.Type != "rock" && card.Type != "paper" && card.Type != "scissors" {
		return errors.New("tipo de carta inválido")
	}
	if card.Stars < 1 || card.Stars > 5 {
		return errors.New("quantidade de estrelas inválida")
	}

	game, err := s.getOrCreateGame(gameID)
	if err != nil {
		return err
	}

	room, err := s.roomRepo.Read(gameID)
	if err != nil {
		return err
	}

	if !room.UserIDs.Contains(playerID) {
		return errors.New("jogador não está na sala")
	}

	if _, exists := game.Plays.Get(playerID); exists {
		return errors.New("jogador já jogou neste turno")
	}

	if game.Plays.Size() >= 2 {
		return errors.New("o jogo já está cheio")
	}

	game.Plays.Set(playerID, card)

	if game.Plays.Size() == 2 {
		game.ResultsSeenBy.Clear()
	}

	return s.gameRepo.Update(gameID, game)
}

// ResetRound redefine o estado de uma partida para o próximo turno.
func (s *GameService) ResetRound(gameID string) error {
	game, err := s.gameRepo.Read(gameID)
	if err != nil {
		return err
	}

	game.Plays.Clear()
	game.ResultsSeenBy.Clear()
	game.FailedAttempts = utils.NewMap[string, int]()

	return s.gameRepo.Update(gameID, game)
}