package state

import (
	"server-of-hope/internal/application"
	"server-of-hope/internal/data"
	"server-of-hope/internal/domain"
)

// AuthService fornece autenticação de usuários.
var AuthService application.AuthServiceInterface

// RoomService gerencia as salas do sistema.
var RoomService application.RoomServiceInterface

// ChatService gerencia o envio e recebimento de mensagens entre usuários.
var ChatService application.ChatServiceInterface

// StoreService gerencia os pacotes de cartas disponíveis na loja.
var StoreService application.StoreServiceInterface

// GameService gerencia a lógica das partidas do jogo.
var GameService application.GameServiceInterface

// UserRepository armazena os dados dos usuários.
var UserRepository data.RepositoryInterface[domain.User]

// RoomRepository armazena os dados das salas.
var RoomRepository data.RepositoryInterface[domain.Room]

// GameRepository armazena os dados das partidas.
var GameRepository data.RepositoryInterface[domain.Game]
