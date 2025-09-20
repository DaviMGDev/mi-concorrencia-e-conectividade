package state

import (
	"server-of-hope/internal/application"
	"server-of-hope/internal/data"
	"server-of-hope/internal/domain"
)

// Initialize inicializa os repositórios e serviços globais do servidor.
func Initialize() {
	/* 	InitializeLogger() */

	UserRepository = data.NewInMemoryRepository[domain.User]()
	RoomRepository = data.NewInMemoryRepository[domain.Room]()
	StoreService = application.NewStoreService()
	GameRepository = data.NewInMemoryRepository[domain.Game]()

	AuthService = application.NewAuthService(UserRepository)
	RoomService = application.NewRoomService(RoomRepository)
	ChatService = application.NewChatService(RoomRepository, UserRepository)
	GameService = application.NewGameService(GameRepository, UserRepository, RoomRepository)
}

// Finalize libera os recursos e limpa os repositórios e serviços globais.
func Finalize() {
	/* 	FinalizeLogger() */

	UserRepository = nil
	RoomRepository = nil

	AuthService = nil
	RoomService = nil
	ChatService = nil

}
