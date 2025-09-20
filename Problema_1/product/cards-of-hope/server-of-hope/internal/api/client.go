// Pacote api implementa o cliente TCP utilizado pelo servidor para comunicação com clientes conectados.
package api

import (
	"encoding/json"
	"net"
	"server-of-hope/internal/api/protocol"
)

// Client representa um cliente TCP conectado ao servidor.
//
// Campos:
//   - Address: armazena o endereço remoto do cliente.
//   - Connection: conexão TCP ativa com o cliente.
//   - Encoder: codifica respostas em JSON para envio ao cliente.
//   - Decoder: decodifica requisições JSON recebidas do cliente.
type Client struct {
	Address    string
	Connection net.Conn
	Encoder    *json.Encoder
	Decoder    *json.Decoder
}

// ClientInterface define a interface para comunicação com clientes TCP.
//
// Métodos:
//   - Send: envia uma resposta para o cliente.
//   - Receive: recebe uma requisição do cliente.
//   - Close: encerra a conexão com o cliente.
type ClientInterface interface {
	Send(response protocol.Response) error
	Receive() (protocol.Request, error)
	Close() error
}

// NewClient cria e retorna uma nova instância de Client para a conexão fornecida.
//
// Parâmetros:
//   - connection: conexão TCP ativa com o cliente.
//
// Retorno:
//   - *Client: ponteiro para a nova instância de Client.
func NewClient(connection net.Conn) *Client {
	return &Client{
		Address:    connection.RemoteAddr().String(),
		Connection: connection,
		Encoder:    json.NewEncoder(connection),
		Decoder:    json.NewDecoder(connection),
	}
}

// Send envia uma resposta para o cliente codificada em JSON.
//
// Parâmetros:
//   - response: resposta a ser enviada.
//
// Retorno:
//   - error: erro ocorrido no envio, se houver.
func (client *Client) Send(response protocol.Response) error {
	return client.Encoder.Encode(response)
}

// Receive recebe uma requisição do cliente decodificada de JSON.
//
// Retorno:
//   - protocol.Request: requisição recebida.
//   - error: erro ocorrido na leitura, se houver.
func (client *Client) Receive() (protocol.Request, error) {
	var request protocol.Request
	err := client.Decoder.Decode(&request)
	return request, err
}

// Close encerra a conexão TCP com o cliente.
//
// Retorno:
//   - error: erro ocorrido ao fechar a conexão, se houver.
func (client *Client) Close() error {
	return client.Connection.Close()
}
