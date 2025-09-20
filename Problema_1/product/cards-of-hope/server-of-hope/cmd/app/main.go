// main é o ponto de entrada do servidor Cards of Hope.
//
// Este programa inicializa o estado global, configura o servidor TCP, registra rotas de comandos e mantém o servidor em execução.
//
// Fluxo principal:
//   - Inicializa o estado global e recursos do servidor.
//   - Cria o servidor TCP e o roteador de comandos.
//   - Registra rotas para autenticação, sala, chat, jogo e utilidades.
//   - Inicia o servidor e aguarda indefinidamente.
//
// Efeitos colaterais:
//   - Pode encerrar o programa caso haja falha na inicialização.
//   - Registra logs de eventos e erros.
//
// Exemplo de uso:
//
//	go run main.go
package main

import (
	"server-of-hope/internal/api"
	"server-of-hope/internal/api/handlers"
	"server-of-hope/internal/state"
)

func main() {
	state.Initialize()
	defer state.Finalize()

	server := api.NewServer(state.HOST + ":" + state.PORT)
	router := api.NewRouter(server)

	router.AddRoute("register", handlers.HandleRegisterUser)
	router.AddRoute("login", handlers.HandleLoginUser)

	router.AddRoute("create", handlers.HandleCreateRoom)
	router.AddRoute("join", handlers.HandleJoinRoom)
	router.AddRoute("leave", handlers.HandleLeaveRoom)

	router.AddRoute("send", handlers.HandleSendMessage)
	router.AddRoute("fetch", handlers.HandleFetchMessage)

	router.AddRoute("play", handlers.HandlePlayCard)
	router.AddRoute("get_opponent_card", handlers.HandleGetOpponentCard)

	router.AddRoute("buy", handlers.HandleBuyPackage)

	router.AddRoute("ping", handlers.HandlePing)
	server.Start(router)
	defer server.Stop()

	select {}
}
