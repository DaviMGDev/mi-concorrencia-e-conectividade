package application

import (
	"client-of-hope/internal/api"
	"client-of-hope/internal/api/protocol"
	"client-of-hope/internal/ui"
	"sync"
)

// ServerHandler define o tipo de função responsável por processar mensagens push do servidor.
type ServerHandler func(client *api.Client, chat *ui.Chat, response protocol.Response)

// ServerRouter gerencia o mapeamento de métodos de servidor para handlers.
type ServerRouter struct {
	routes map[string]ServerHandler
	mutex  sync.Mutex
	client *api.Client
	chat   *ui.Chat
}

// NewServerRouter cria uma nova instância de ServerRouter.
func NewServerRouter(client *api.Client, chat *ui.Chat) *ServerRouter {
	return &ServerRouter{
		routes: make(map[string]ServerHandler),
		client: client,
		chat:   chat,
	}
}

// AddRoute adiciona uma nova rota de servidor ao roteador.
func (router *ServerRouter) AddRoute(method string, handler ServerHandler) {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	router.routes[method] = handler
}

// HandleResponse processa uma resposta do servidor, executando o handler correspondente.
func (router *ServerRouter) HandleResponse(response protocol.Response) {
	router.mutex.Lock()
	defer router.mutex.Unlock()

	handler, exists := router.routes[response.Method]
	if !exists {
		// Ignora métodos desconhecidos
		return
	}

	go handler(router.client, router.chat, response)
}

// Start inicia o roteador do servidor, ouvindo continuamente por mensagens push.
func (router *ServerRouter) Start() {
	go func() {
		for response := range router.client.PushedMessages {
			router.HandleResponse(response)
		}
	}()
}
