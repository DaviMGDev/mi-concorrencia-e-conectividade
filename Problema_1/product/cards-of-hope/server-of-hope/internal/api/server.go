// Pacote api implementa o servidor TCP do Cards of Hope, gerenciando conexões, clientes e roteamento de comandos.
package api

import (
	"net"
	"server-of-hope/internal/api/protocol"
	"server-of-hope/internal/state"
	"server-of-hope/internal/utils"
)

// Server representa o servidor TCP do Cards of Hope.
//
// Campos:
//   - Address: endereço TCP em que o servidor irá escutar.
//   - Listener: listener TCP ativo do servidor.
//   - Clients: clientes conectados, indexados por ID.
//   - Router: interface responsável pelo roteamento de comandos.
//   - Requests: canal de requisições recebidas.
//   - Responses: canal de respostas a serem enviadas.
type Server struct {
	Address   string
	Listener  net.Listener
	Clients   *utils.Map[string, *Client]
	Router    RouterInterface
	Requests  chan protocol.Request
	Responses chan protocol.Response
}

// NewServer cria e retorna uma nova instância de Server para o endereço fornecido.
//
// Parâmetros:
//   - address: endereço TCP para escutar.
//
// Retorno:
//   - *Server: ponteiro para a nova instância de Server.
func NewServer(address string) *Server {
	return &Server{
		Address:   address,
		Clients:   utils.NewMap[string, *Client](),
		Requests:  make(chan protocol.Request, 100),
		Responses: make(chan protocol.Response, 100),
	}
}

// Start inicia o servidor TCP, configurando o listener, roteador e goroutines de conexão e resposta.
//
// Parâmetros:
//   - router: instância do roteador de comandos.
//
// Retorno:
//   - error: erro ocorrido ao iniciar o servidor, se houver.
func (server *Server) Start(router RouterInterface) error {
	listener, err := net.Listen("tcp", server.Address)
	if err != nil {
		state.Logger.Error("Falha ao iniciar o servidor", "erro", err)
		return err
	}
	server.Listener = listener
	state.Logger.Info("Servidor iniciado", "endereco", server.Address)

	server.Router = router
	go server.Router.Start()
	go server.acceptConnections()
	go server.handleResponses()

	return nil
}

// Stop encerra o listener TCP do servidor, se estiver ativo.
//
// Retorno:
//   - error: erro ocorrido ao encerrar o listener, se houver.
func (server *Server) Stop() error {
	if server.Listener != nil {
		err := server.Listener.Close()
		if err != nil {
			state.Logger.Error("Falha ao encerrar o listener do servidor", "erro", err)
			return err
		}
		state.Logger.Info("Servidor encerrado")
		return nil
	}
	return nil
}

func (server *Server) acceptConnections() {
	for {
		conn, err := server.Listener.Accept()
		if err != nil {
			state.Logger.Error("Failed to accept connection", "error", err)
			continue
		}
		client := NewClient(conn)
		server.Clients.Set(client.Address, client)
		state.Logger.Info("Client connected", "address", client.Address)
		go server.getRequests(client)
	}
}

func (server *Server) getRequests(client *Client) {
	defer func() {
		client.Close()
		server.Clients.Delete(client.Address)
		if client.UserID != "" {
			state.UserConnections.Delete(client.UserID)
		}
		state.Logger.Info("Client disconnected", "address", client.Address)
	}()
	for {
		request, err := client.Receive()
		if err != nil {
			state.Logger.Error("Failed to receive request from client", "address", client.Address, "error", err)
			break
		}
		request.From = client.Address
		server.Requests <- request
	}
}
func (server *Server) handleResponses() {
	for response := range server.Responses {
		client, exists := server.Clients.Get(response.To)
		if exists {
			client.Send(response)
			state.Logger.Info("Response sent", "to", response.To, "method", response.Method, "status", response.Status)
		} else {
			state.Logger.Warn("Client not found for response", "to", response.To, "method", response.Method)
		}
	}
}
