package config

import (
	"math/rand"
	"os"
	"server-of-dispair/internal/domain"
	"server-of-dispair/internal/repositories"
	"server-of-dispair/internal/services"

	"github.com/charmbracelet/log"
)

var (
	HOST = "localhost"
	PORT = "8080"

	AuthService  *services.AuthService
	RoomService  *services.RoomService
	StoreService *services.StoreService
	GameService  *services.GameService
	ChatService  *services.ChatService
	Logger       *log.Logger
)

func Initialize() {
	// Logger
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: "[server-of-dispair]",
		Level:  log.DebugLevel,
	})

	// Repositories
	userRepo := repositories.NewInMemoryRepository[domain.User]()
	roomRepo := repositories.NewInMemoryRepository[domain.Room]()

	// domain

	// Services
	AuthService = services.NewAuthService(userRepo)
	RoomService = services.NewRoomService(roomRepo)
	StoreService = services.NewStoreService()
	StoreService.AddPackage(domain.CardPackage{
		Cards: [3]domain.Card{
			{Type: "rock", Stars: rand.Intn(5) + 1},
			{Type: "paper", Stars: rand.Intn(5) + 1},
			{Type: "scissors", Stars: rand.Intn(5) + 1},
		},
	})
	GameService = services.NewGameService(userRepo, roomRepo)
	ChatService = services.NewChatService(userRepo, roomRepo)

}

func Finalize() {
	Logger.Info("Finalizing application.")
	// Add any cleanup logic here
}
