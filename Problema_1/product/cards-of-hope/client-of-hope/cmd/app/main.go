// main é o ponto de entrada do cliente Cards of Hope.
//
// Este programa inicializa o estado global, configura a interface de chat, conecta ao servidor,
// registra rotas de comandos e mantém a aplicação em execução até o encerramento pelo usuário.
//
// Fluxo principal:
//   - Inicializa o logger e o estado global.
//   - Cria e inicia a interface de chat.
//   - Obtém o endereço do servidor a partir da variável de ambiente SERVER_ADDR (ou usa localhost:8080).
//   - Cria o cliente de API e tenta conectar ao servidor.
//   - Registra rotas de comandos para autenticação, chat, sala, jogo e utilitários.
//   - Inicia o roteador e aguarda o encerramento do chat.
//
// Efeitos colaterais:
//   - Pode encerrar o programa caso não consiga conectar ao servidor.
//   - Registra logs de eventos e erros.
//
// Exemplo de uso:
//
//	go run main.go
package main

import (
	"client-of-hope/internal/api"
	"client-of-hope/internal/api/handlers"
	"client-of-hope/internal/application"
	"client-of-hope/internal/state"
	"client-of-hope/internal/ui"
	"os"
)

// getServerAddress retorna o endereço do servidor a partir da variável de ambiente SERVER_ADDR.
// Caso a variável não esteja definida, retorna "localhost:8080" como padrão.
//
// Retorno:
//   - string: endereço do servidor no formato host:porta.
func getServerAddress() string {
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		return "localhost:8080"
	}
	return addr
}

func main() {
	state.Initialize()
	defer state.CloseLogger()

	chat := ui.NewChat()

	serverAddress := getServerAddress()
	client := api.NewClient(serverAddress)
	err := client.Connect()
	if err != nil {
		state.Log("Falha ao conectar ao servidor em %s: %v", serverAddress, err)
		chat.Outputs <- "Falha ao conectar ao servidor. Por favor, certifique-se de que o servidor está em execução."
		os.Exit(1)
	}
	
	chat.Start()
	defer client.Close()

	router := application.NewRouter(client, chat)

	// Autenticação
	router.AddRoute("register", handlers.HandleRegister)
	router.AddRoute("login", handlers.HandleLogin)
	router.AddRoute("logout", handlers.HandleLogout)

	// Chat
	router.AddRoute("send", handlers.HandleSendMessage)
	router.AddRoute("fetch", handlers.HandleFetchMessage)

	// Sala
	router.AddRoute("create", handlers.HandleCreateRoom)
	router.AddRoute("join", handlers.HandleJoinRoom)
	router.AddRoute("leave", handlers.HandleLeaveRoom)

	// Jogo
	router.AddRoute("play", handlers.HandlePlay)
	router.AddRoute("cards", handlers.HandleCards)
	router.AddRoute("buy", handlers.HandleBuy)

	// Diversos
	router.AddRoute("whoami", handlers.HandleWhoami)
	router.AddRoute("whereami", handlers.HandleWhereami)
	router.AddRoute("quit", handlers.HandleQuit)
	router.AddRoute("ping", handlers.HandlePing)
	router.AddRoute("help", handlers.HandleHelp)

	router.Start()

	serverRouter := application.NewServerRouter(client, chat)
	serverRouter.AddRoute("opponent_played", handlers.HandleOpponentPlayed)
	serverRouter.Start()

	// Mantém a goroutine principal viva aguardando o sinal de conclusão do chat.
	<-chat.Done
}