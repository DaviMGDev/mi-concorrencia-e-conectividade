// Pacote handlers contém funções para lidar com comandos de autenticação do usuário no Cards of Hope.
package handlers

import (
	"client-of-hope/internal/api"
	"client-of-hope/internal/api/protocol"
	"client-of-hope/internal/state"
	"client-of-hope/internal/ui"
	"client-of-hope/internal/utils"
)

// HandleRegister processa o comando de registro de novo usuário.
//
// Parâmetros:
//   - client: ponteiro para o cliente de API conectado ao servidor.
//   - chat: ponteiro para a interface de chat para exibir mensagens ao usuário.
//   - args: slice de strings contendo nome de usuário e senha.
//
// Efeitos colaterais:
//   - Envia requisição de registro ao servidor.
//   - Exibe mensagens de sucesso ou erro no chat.
func HandleRegister(client *api.Client, chat *ui.Chat, args []string) {
	if len(args) < 2 {
		chat.Outputs <- "Uso: /register <usuario> <senha>"
		return
	}
	username := args[0]
	password := args[1]

	request := protocol.Request{
		Method: "register",
		Data: utils.Dict{
			"username": username,
			"password": password,
		},
	}

	response, err := client.DoRequest(request)
	if err != nil {
		state.Log("Falha na requisição de registro: %v", err)
		chat.Outputs <- "Falha na requisição de registro."
		return
	}

	if response.Status != "ok" {
		message, _ := response.Data["message"].(string)
		chat.Outputs <- message
		return
	}

	chat.Outputs <- "Registro realizado com sucesso!"
}

// HandleLogin processa o comando de login do usuário.
//
// Parâmetros:
//   - client: ponteiro para o cliente de API conectado ao servidor.
//   - chat: ponteiro para a interface de chat para exibir mensagens ao usuário.
//   - args: slice de strings contendo nome de usuário e senha.
//
// Efeitos colaterais:
//   - Envia requisição de login ao servidor.
//   - Atualiza o estado global com nome de usuário e ID.
//   - Exibe mensagens de sucesso ou erro no chat.
func HandleLogin(client *api.Client, chat *ui.Chat, args []string) {
	if len(args) < 2 {
		chat.Outputs <- "Uso: /login <usuario> <senha>"
		return
	}
	username := args[0]
	password := args[1]

	request := protocol.Request{
		Method: "login",
		Data: utils.Dict{
			"username": username,
			"password": password,
		},
	}

	response, err := client.DoRequest(request)
	if err != nil {
		state.Log("Falha na requisição de login: %v", err)
		chat.Outputs <- "Falha na requisição de login."
		return
	}

	if response.Status != "ok" {
		message, _ := response.Data["message"].(string)
		chat.Outputs <- message
		return
	}

	userID, ok := response.Data["user_id"].(string)
	if !ok {
		chat.Outputs <- "Falha no login: resposta do servidor não incluiu o ID do usuário."
		return
	}

	state.Username = username
	state.UserID = userID
	chat.Outputs <- "Login realizado com sucesso como " + username
}

// HandleLogout processa o comando de logout do usuário.
//
// Parâmetros:
//   - client: ponteiro para o cliente de API conectado ao servidor.
//   - chat: ponteiro para a interface de chat para exibir mensagens ao usuário.
//   - args: slice de strings (não utilizado neste comando).
//
// Efeitos colaterais:
//   - Limpa o nome de usuário e ID do estado global.
//   - Exibe mensagens de sucesso ou erro no chat.
func HandleLogout(client *api.Client, chat *ui.Chat, args []string) {
	if state.Username == "" {
		chat.Outputs <- "Você não está logado."
		return
	}
	if state.RoomID != "" {
		chat.Outputs <- "Você deve sair da sala antes de fazer logout."
		return
	}
	state.Username = ""
	state.UserID = ""
	chat.Outputs <- "Logout realizado com sucesso."
}
