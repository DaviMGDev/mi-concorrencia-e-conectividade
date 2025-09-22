// Pacote api fornece a implementação do cliente TCP para comunicação com o servidor do Cards of Hope.
//
// Este pacote define a estrutura Client, responsável por gerenciar a conexão, envio e recebimento de mensagens
// utilizando o protocolo definido em protocol.Request e protocol.Response.
package api

import (
	"client-of-hope/internal/api/protocol"
	"encoding/json"
	"errors"
	"net"
	"sync"
	"time"
)

// Client representa um cliente TCP que se comunica com o servidor do Cards of Hope.
//
// Campos:
//   - Address: endereço do servidor (host:porta).
//   - Connection: conexão TCP ativa com o servidor.
//   - Mutex: garante acesso concorrente seguro à conexão.
//   - Encoder: codificador JSON para envio de mensagens.
//   - Decoder: decodificador JSON para recebimento de mensagens.
type Client struct {
	Address            string
	Connection         net.Conn
	Mutex              sync.Mutex
	Encoder            *json.Encoder
	Decoder            *json.Decoder
	PushedMessages     chan protocol.Response
	requestResponseMap sync.Map // map[string]chan protocol.Response
}

// NewClient cria uma nova instância de Client para o endereço fornecido.
func NewClient(address string) *Client {
	return &Client{
		Address:        address,
		Mutex:          sync.Mutex{},
		PushedMessages: make(chan protocol.Response, 10),
	}
}

// Send envia uma requisição para o servidor de forma thread-safe.
func (client *Client) Send(request protocol.Request) error {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	return client.Encoder.Encode(request)
}

// Receive recebe uma resposta do servidor.
func (client *Client) Receive() (protocol.Response, error) {
	var response protocol.Response
	err := client.Decoder.Decode(&response)
	return response, err
}

// Close encerra a conexão TCP com o servidor.
func (client *Client) Close() error {
	return client.Connection.Close()
}

// Connect estabelece la conexión TCP con el servidor e inicializa los codificadores/decodificadores JSON.
func (client *Client) Connect() error {
	conn, err := net.Dial("tcp", client.Address)
	if err != nil {
		return err
	}
	client.Connection = conn
	client.Encoder = json.NewEncoder(conn)
	client.Decoder = json.NewDecoder(conn)
	go client.Listen()
	return nil
}

// Listen escuta continuamente por mensagens do servidor e as distribui.
func (client *Client) Listen() {
	for {
		var response protocol.Response
		err := client.Decoder.Decode(&response)
		if err != nil {
			// Conexão provavelmente fechada, encerra o listener.
			close(client.PushedMessages)
			return
		}

		// Verifica se é uma resposta para uma requisição DoRequest
		if ch, ok := client.requestResponseMap.Load(response.Method); ok {
			responseChan, _ := ch.(chan protocol.Response)
			responseChan <- response
			client.requestResponseMap.Delete(response.Method)
		} else {
			// Se não, é uma mensagem push do servidor
			client.PushedMessages <- response
		}
	}
}

// DoRequest envia uma requisição ao servidor e aguarda a resposta.
func (client *Client) DoRequest(request protocol.Request) (protocol.Response, error) {
	responseChan := make(chan protocol.Response, 1)
	client.requestResponseMap.Store(request.Method, responseChan)

	if err := client.Send(request); err != nil {
		return protocol.Response{}, err
	}

	// Aguarda a resposta com um timeout
	select {
	case response := <-responseChan:
		return response, nil
	case <-time.After(15 * time.Second):
		client.requestResponseMap.Delete(request.Method)
		return protocol.Response{}, errors.New("request timed out")
	}
}
