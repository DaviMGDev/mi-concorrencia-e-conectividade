// Pacote api fornece a implementação do cliente TCP para comunicação com o servidor do Cards of Hope.
//
// Este pacote define a estrutura Client, responsável por gerenciar a conexão, envio e recebimento de mensagens
// utilizando o protocolo definido em protocol.Request e protocol.Response.
package api

import (
	"client-of-hope/internal/api/protocol"
	"encoding/json"
	"net"
	"sync"
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
	Address    string
	Connection net.Conn
	Mutex      sync.Mutex
	Encoder    *json.Encoder
	Decoder    *json.Decoder
}

// NewClient cria uma nova instância de Client para o endereço fornecido.
//
// Parâmetros:
//   - address: endereço do servidor (host:porta).
//
// Retorno:
//   - *Client: ponteiro para a nova instância de Client.
func NewClient(address string) *Client {
	return &Client{
		Address: address,
		Mutex:   sync.Mutex{},
	}
}

// Send envia uma requisição para o servidor de forma thread-safe.
//
// Parâmetros:
//   - request: estrutura protocol.Request a ser enviada.
//
// Retorno:
//   - error: erro ocorrido no envio, se houver.
func (client *Client) Send(request protocol.Request) error {
	client.Mutex.Lock()
	defer client.Mutex.Unlock()
	return client.Encoder.Encode(request)
}

// Receive recebe uma resposta do servidor.
//
// Retorno:
//   - protocol.Response: resposta recebida do servidor.
//   - error: erro ocorrido na leitura, se houver.
func (client *Client) Receive() (protocol.Response, error) {
	var response protocol.Response
	err := client.Decoder.Decode(&response)
	return response, err
}

// Close encerra a conexão TCP com o servidor.
//
// Retorno:
//   - error: erro ocorrido ao fechar a conexão, se houver.
func (client *Client) Close() error {
	return client.Connection.Close()
}

// Connect estabelece a conexão TCP com o servidor e inicializa os encoders/decoders JSON.
//
// Retorno:
//   - error: erro ocorrido ao tentar conectar, se houver.
func (client *Client) Connect() error {
	conn, err := net.Dial("tcp", client.Address)
	if err != nil {
		return err
	}
	client.Connection = conn
	client.Encoder = json.NewEncoder(conn)
	client.Decoder = json.NewDecoder(conn)
	return nil
}

// DoRequest envia uma requisição ao servidor e aguarda a resposta.
//
// Parâmetros:
//   - request: estrutura protocol.Request a ser enviada.
//
// Retorno:
//   - protocol.Response: resposta recebida do servidor.
//   - error: erro ocorrido durante o envio ou recebimento.
func (client *Client) DoRequest(request protocol.Request) (protocol.Response, error) {
	if err := client.Send(request); err != nil {
		return protocol.Response{}, err
	}
	return client.Receive()
}
