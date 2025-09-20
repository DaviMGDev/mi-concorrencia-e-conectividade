package application

import (
	"errors"
	"server-of-hope/internal/data"
	"server-of-hope/internal/domain"
	"server-of-hope/internal/utils"
)

// GameServiceInterface descreve as operações para manipulação da lógica do jogo.
//
// Métodos:
//   - PlayCard: registra a jogada de um jogador.
//   - GetOpponentCard: retorna a carta do oponente.
type GameServiceInterface interface {
	// PlayCard registra a jogada de um jogador em uma partida.
	//
	// Parâmetros:
	//   - gameID: identificador da partida.
	//   - playerID: identificador do jogador.
	//   - card: carta jogada pelo jogador.
	//
	// Retorno:
	//   - erro caso a jogada seja inválida ou não possa ser registrada.
	PlayCard(gameID string, playerID string, card domain.Card) error

	// GetOpponentCard retorna a carta jogada pelo oponente.
	//
	// Parâmetros:
	//   - gameID: identificador da partida.
	//   - playerID: identificador do jogador.
	//
	// Retorno:
	//   - Card: carta do oponente.
	//   - erro caso o oponente não tenha jogado ou haja falha na busca.
	GetOpponentCard(gameID string, playerID string) (domain.Card, error)
}

// GameService implementa a lógica do jogo, incluindo jogadas e controle de estado.
//
// Campos:
//   - gameRepo: repositório das partidas.
//   - userRepo: repositório dos usuários.
//   - roomRepo: repositório das salas.
type GameService struct {
	gameRepo data.RepositoryInterface[domain.Game]
	userRepo data.RepositoryInterface[domain.User]
	roomRepo data.RepositoryInterface[domain.Room]
}

// NewGameService cria uma nova instância de GameService.
//
// Parâmetros:
//   - gameRepo: repositório das partidas.
//   - userRepo: repositório dos usuários.
//   - roomRepo: repositório das salas.
//
// Retorno:
//   - ponteiro para GameService.
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
//
// Parâmetros:
//   - gameID: identificador da partida.
//
// Retorno:
//   - Game: partida encontrada ou criada.
//   - erro caso não seja possível buscar ou criar a partida.
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

// PlayCard registra a jogada de um jogador em uma partida, validando a carta e o estado do jogo.
//
// Parâmetros:
//   - gameID: identificador da partida.
//   - playerID: identificador do jogador.
//   - card: carta jogada pelo jogador.
//
// Retorno:
//   - erro caso a jogada seja inválida ou não possa ser registrada.
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

// GetOpponentCard retorna a carta jogada pelo oponente, se disponível, e atualiza o estado do jogo.
//
// Parâmetros:
//   - gameID: identificador da partida.
//   - playerID: identificador do jogador.
//
// Retorno:
//   - Card: carta do oponente.
//   - erro caso o oponente não tenha jogado ou haja falha na busca.
func (s *GameService) GetOpponentCard(gameID string, playerID string) (domain.Card, error) {
	game, err := s.gameRepo.Read(gameID)
	if err != nil {
		return domain.Card{}, err
	}

	if _, exists := game.Plays.Get(playerID); !exists {
		return domain.Card{}, errors.New("jogador ainda não jogou")
	}

	var opponentCard domain.Card
	var opponentID string

	game.Plays.ForEach(func(id string, card domain.Card) {
		if id != playerID {
			opponentID = id
			opponentCard = card
		}
	})

	if opponentID == "" {
		// Incrementa tentativas frustradas
		attempts, _ := game.FailedAttempts.Get(playerID)
		attempts++
		game.FailedAttempts.Set(playerID, attempts)
		if attempts >= 3 {
			// Reseta jogada e tentativas para o jogador
			game.Plays.Delete(playerID)
			game.FailedAttempts.Set(playerID, 0)
			s.gameRepo.Update(gameID, game)
			return domain.Card{}, errors.New("oponente ainda não jogou. Seu turno foi resetado após 3 tentativas. Jogue novamente.")
		}
		s.gameRepo.Update(gameID, game)
		return domain.Card{}, errors.New("oponente ainda não jogou")
	}

	game.ResultsSeenBy.Add(playerID)
	// Resetar tentativas frustradas ao sucesso
	game.FailedAttempts.Set(playerID, 0)
	if game.ResultsSeenBy.Size() >= 2 {
		game.Plays.Clear()
		game.ResultsSeenBy.Clear()
		// Limpa tentativas frustradas de todos
		game.FailedAttempts = utils.NewMap[string, int]()
	}

	if err := s.gameRepo.Update(gameID, game); err != nil {
		return domain.Card{}, err
	}

	return opponentCard, nil
}
