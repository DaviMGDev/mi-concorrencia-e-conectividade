// Pacote protocol define as estruturas de requisição e resposta utilizadas na comunicação entre cliente e servidor.
package protocol

import (
	"client-of-hope/internal/utils"
	"fmt"
)

// Response representa uma resposta enviada do servidor para o cliente.
//
// Campos:
//   - Method: nome do método/comando relacionado à resposta.
//   - Status: status da resposta (ex: "ok", "error").
//   - Data: dicionário de dados adicionais retornados pelo servidor.
type Response struct {
	Method string     `json:"method"`
	Status string     `json:"status"`
	Data   utils.Dict `json:"data,omitempty"`
}

// String retorna uma representação textual da resposta para fins de depuração e logging.
//
// Retorno:
//   - string: representação formatada do status, método e dados da resposta.
func (r Response) String() string {
	return fmt.Sprintf("Status: %s, Method: %s, Data: %v", r.Status, r.Method, r.Data)
}
