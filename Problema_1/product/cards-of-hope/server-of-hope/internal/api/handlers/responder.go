package handlers

import (
	"server-of-hope/internal/api"
	"server-of-hope/internal/api/protocol"
	"server-of-hope/internal/state"
	"server-of-hope/internal/utils"
)

// Responder encapsula a lógica de resposta para um manipulador de requisição.
type Responder struct {
	server   *api.Server
	request  protocol.Request
	response protocol.Response
}

// NewResponder cria uma nova instância de Responder.
func NewResponder(server *api.Server, request protocol.Request) *Responder {
	return &Responder{
		server:  server,
		request: request,
		response: protocol.Response{
			Method: request.Method,
			Status: "error", // O status padrão é 'error'
			Data:   utils.Dict{},
			To:     request.From,
		},
	}
}

// Send envia a resposta para o canal de respostas do servidor.
// Deve ser chamado no final do manipulador, preferencialmente com defer.
func (r *Responder) Send() {
	r.server.Responses <- r.response
}

// SetSuccess atualiza a resposta para um estado de sucesso.
func (r *Responder) SetSuccess(data utils.Dict, logMessage string, logFields ...any) {
	r.response.Status = "ok"
	r.response.Data = data
	state.Logger.Info(logMessage, logFields...)
}

// SetError atualiza a resposta para um estado de erro.
func (r *Responder) SetError(errorMessage string, logMessage string, logFields ...any) {
	r.response.Data["message"] = errorMessage
	state.Logger.Error(logMessage, logFields...)
}
