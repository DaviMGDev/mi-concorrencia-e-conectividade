// Pacote application implementa o roteador de comandos do Cards of Hope, responsável por mapear comandos de texto para funções manipuladoras.
package application

import (
	"client-of-hope/internal/api"
	"client-of-hope/internal/ui"
	"strings"
	"sync"
)

// Handler define o tipo de função responsável por processar comandos recebidos do usuário.
// Parâmetros:
//   - client: ponteiro para o cliente de API conectado ao servidor.
//   - chat: ponteiro para a interface de chat.
//   - args: slice de strings com argumentos do comando.
type Handler func(client *api.Client, chat *ui.Chat, args []string)

// Router gerencia o mapeamento de comandos para handlers e executa comandos recebidos do chat.
//
// Campos:
//   - routes: mapa de comandos para funções manipuladoras.
//   - mutex: garante acesso concorrente seguro ao mapa de rotas.
//   - client: cliente de API conectado ao servidor.
//   - chat: interface de chat para entrada e saída de mensagens.
type Router struct {
	routes map[string]Handler
	mutex  sync.Mutex
	client *api.Client
	chat   *ui.Chat
}

// NewRouter cria uma nova instância de Router associada ao cliente e chat fornecidos.
//
// Parâmetros:
//   - client: ponteiro para o cliente de API.
//   - chat: ponteiro para a interface de chat.
//
// Retorno:
//   - *Router: ponteiro para a nova instância de Router.
func NewRouter(client *api.Client, chat *ui.Chat) *Router {
	return &Router{
		routes: make(map[string]Handler),
		client: client,
		chat:   chat,
	}
}

// decodeInput interpreta a entrada do usuário, separando o comando e seus argumentos.
//
// Parâmetros:
//   - input: string de entrada do usuário.
//
// Retorno:
//   - string: comando extraído (sem barra inicial).
//   - []string: argumentos do comando.
func decodeInput(input string) (string, []string) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return "", nil
	}
	if !strings.HasPrefix(args[0], "/") {
		return "send", args
	}
	return args[0][1:], args[1:]
}

// AddRoute adiciona uma nova rota de comando ao roteador.
//
// Parâmetros:
//   - command: nome do comando (sem barra).
//   - handler: função manipuladora a ser chamada para o comando.
func (router *Router) AddRoute(command string, handler Handler) {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	router.routes[command] = handler
}

// HandleCommand processa a entrada do usuário, identifica o comando e executa o handler correspondente.
//
// Parâmetros:
//   - input: string de entrada do usuário.
//
// Efeitos colaterais:
//   - Executa o handler do comando em uma nova goroutine.
//   - Exibe mensagem de erro caso o comando não exista.
func (router *Router) HandleCommand(input string) {
	router.mutex.Lock()
	defer router.mutex.Unlock()
	command, args := decodeInput(strings.TrimSpace(input))
	if command == "" {
		return
	}
	handler, exists := router.routes[command]
	if !exists {
		router.chat.Outputs <- "Comando desconhecido: " + command
		return
	}
	// Executa o handler em uma nova goroutine para não bloquear o chat
	go handler(router.client, router.chat, args)
}

// Start inicia o roteador, ouvindo continuamente por comandos do usuário e processando-os.
//
// Efeitos colaterais:
//   - Executa HandleCommand para cada entrada recebida do chat.
func (router *Router) Start() {
	go func() {
		for input := range router.chat.Inputs {
			router.HandleCommand(input)
		}
	}()
}
