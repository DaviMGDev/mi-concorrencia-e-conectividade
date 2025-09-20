// Pacote protocol define as estruturas de requisição e resposta utilizadas na comunicação entre cliente e servidor.
package protocol

import "client-of-hope/internal/utils"

// Request representa uma requisição enviada do cliente para o servidor.
//
// Campos:
//   - Method: nome do método/comando a ser executado no servidor.
//   - Data: dicionário de dados adicionais necessários para o comando.
type Request struct {
	Method string     `json:"method"`
	Data   utils.Dict `json:"data,omitempty"`
}
