// Pacote state gerencia o estado global da aplicação, incluindo o sistema de logging.
package state

import (
	"log"
	"os"
)

// LogPath define o caminho do arquivo de log.
// LogFile representa o arquivo de log aberto.
// logger é a instância do logger configurado para o cliente.
var (
	LogPath = "data/log.txt"
	LogFile *os.File
	logger  *log.Logger
)

// InitLogger inicializa o sistema de logging, criando o diretório e arquivo de log se necessário.
//
// Efeitos colaterais:
//   - Cria o diretório de dados se não existir.
//   - Abre ou cria o arquivo de log e inicializa o logger.
func InitLogger() {
	err := os.MkdirAll("data", 0755)
	if err != nil {
		log.Fatalf("Falha ao criar diretório de log: %v", err)
	}

	LogFile, err = os.OpenFile(LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Falha ao abrir arquivo de log: %v", err)
	}
	logger = log.New(LogFile, "CLIENTE: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Log registra uma mensagem formatada no arquivo de log, se o logger estiver inicializado.
//
// Parâmetros:
//   - format: string de formatação.
//   - v: argumentos variádicos para formatação.
func Log(format string, v ...interface{}) {
	if logger != nil {
		logger.Printf(format, v...)
	}
}

// CloseLogger fecha o arquivo de log, se estiver aberto.
func CloseLogger() {
	if LogFile != nil {
		LogFile.Close()
	}
}
