// Pacote state gerencia o estado global da aplicação, incluindo inicialização de logs e cartas do jogo.
package state

import (
	"log"
	"os"
)

// Initialize inicializa o estado global da aplicação.
//
// Efeitos colaterais:
//   - Abre o arquivo de log e redireciona a saída de log para ele.
//   - Inicializa as cartas do jogo com valores padrão.
func Initialize() {
	LogFile, err := os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Falha ao abrir o arquivo de log: %v", err)
	}
	log.SetOutput(LogFile)
	Cards.Set("rock", 1)
	Cards.Set("paper", 1)
	Cards.Set("scissors", 1)
}

// Finalize encerra o estado global da aplicação (placeholder para futuras finalizações).
func Finalize() {}
