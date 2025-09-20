// Pacote state armazena o estado do jogo, incluindo cartas, jogadas e regras de vitória.
package state

import "client-of-hope/internal/utils"

// Cards armazena o número de cartas disponíveis para o usuário.
// PlayedCard representa a última carta jogada pelo usuário.
// PlayedCardStar representa o valor especial da carta jogada pelo usuário.
// OpponentCard representa a última carta jogada pelo oponente.
// OpponentCardStar representa o valor especial da carta do oponente.
// CardWins define as regras de vitória entre as cartas (pedra, papel, tesoura).
var (
	Cards            *utils.Map[string, int] = utils.NewMap[string, int]()
	PlayedCard       string                  = ""
	PlayedCardStar   int                     = 0
	OpponentCard     string                  = ""
	OpponentCardStar int                     = 0
	CardWins         map[string]string       = map[string]string{
		"rock":     "scissors",
		"paper":    "rock",
		"scissors": "paper",
	}
)
