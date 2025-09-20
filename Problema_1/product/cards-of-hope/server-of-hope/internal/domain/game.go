package domain

import "server-of-hope/internal/utils"

// Game representa uma partida do jogo.
//
// Campos:
//   - ID: identificador único da partida.
//   - Plays: jogadas dos jogadores.
//   - ResultsSeenBy: jogadores que visualizaram o resultado.
type Game struct {
	ID             string                   `json:"id"`
	Plays          *utils.Map[string, Card] `json:"plays"`
	ResultsSeenBy  *utils.Set[string]       `json:"results_seen_by"`
	FailedAttempts *utils.Map[string, int]  `json:"failed_attempts"` // tentativas frustradas por jogador
}
