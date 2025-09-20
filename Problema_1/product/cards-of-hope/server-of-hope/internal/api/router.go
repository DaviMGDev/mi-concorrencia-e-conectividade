// Pacote api implementa o roteador de comandos do servidor Cards of Hope, responsável por mapear métodos para handlers.
package api

import (
	"server-of-hope/internal/api/protocol"
	"server-of-hope/internal/state"
	"sync"
)

// HandlerFunc define o tipo de função responsável por processar requisições recebidas pelo servidor.
// Parâmetros:
//   - server: ponteiro para o servidor.
//   - request: requisição recebida.
type HandlerFunc func(server *Server, request protocol.Request)

// RouterInterface define a interface para roteadores de comandos.
type RouterInterface interface {
	AddRoute(method string, handler HandlerFunc)
	HandleRequest(server *Server, request protocol.Request)
	Start()
}

// Router implementa o roteador de comandos do servidor.
//
// Campos:
//   - routes: mapeamento de métodos para handlers.
//   - server: ponteiro para o servidor associado.
//   - mutex: garante acesso concorrente seguro ao mapa de rotas.
type Router struct {
	routes map[string]HandlerFunc
	server *Server
	mutex  sync.Mutex
}

// NewRouter cria e retorna uma nova instância de Router associada ao servidor fornecido.
//
// Parâmetros:
//   - server: ponteiro para o servidor.
//
// Retorno:
//   - *Router: ponteiro para a nova instância de Router.
func NewRouter(server *Server) *Router {
	return &Router{
		routes: make(map[string]HandlerFunc),
		server: server,
	}
}

// AddRoute adiciona um novo método e seu handler ao roteador.
//
// Parâmetros:
//   - method: nome do método/comando.
//   - handler: função handler a ser chamada para o método.
func (router *Router) AddRoute(method string, handler HandlerFunc) {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	router.routes[method] = handler
}

// HandleRequest processa uma requisição recebida, executando o handler correspondente ou retornando erro se desconhecido.
//
// Parâmetros:
//   - server: ponteiro para o servidor.
//   - request: requisição recebida.
func (router *Router) HandleRequest(server *Server, request protocol.Request) {
	router.mutex.Lock()
	defer router.mutex.Unlock()

	state.Logger.Info("Processando requisição", "método", request.Method, "de", request.From)

	if handler, exists := router.routes[request.Method]; exists {
		go handler(server, request)
	} else {
		state.Logger.Warn("Método desconhecido recebido", "método", request.Method, "de", request.From)
		response := protocol.Response{
			Method: request.Method,
			Status: "error",
			Data: map[string]any{
				"message": "Método desconhecido",
			},
			To: request.From,
		}
		server.Responses <- response
	}
}

// Start inicia o roteador, processando continuamente as requisições recebidas do servidor.
func (router *Router) Start() {
	for request := range router.server.Requests {
		router.HandleRequest(router.server, request)
	}
}
