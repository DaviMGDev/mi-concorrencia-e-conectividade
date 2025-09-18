package main

import (
	"os"
	"os/signal"
	"server-of-dispair/internal/config"
	"server-of-dispair/internal/handlers"
	"server-of-dispair/internal/protocol"
	"syscall"
)

func main() {
	config.Initialize()
	defer config.Finalize()

	server := protocol.NewServer(config.HOST, config.PORT)

	/* 	// Register routes
	   	server.Router.AddRoute("ping", handlers.HandlePing)
	   	server.Router.AddRoute("REGISTER", handlers.HandleRegister)
	   	server.Router.AddRoute("LOGIN", handlers.HandleLogin)
	   	server.Router.AddRoute("CREATE_ROOM", handlers.HandleCreateRoom)
	   	server.Router.AddRoute("JOIN_ROOM", handlers.HandleJoinRoom)
	   	server.Router.AddRoute("LEAVE_ROOM", handlers.HandleLeaveRoom)
	   	server.Router.AddRoute("READY", handlers.HandleReady)     // Note: GameService not fully implemented
	   	server.Router.AddRoute("UNREADY", handlers.HandleUnready) // Note: GameService not fully implemented
	   	server.Router.AddRoute("PLAY", handlers.HandlePlay)       // Note: GameService not fully implemented
	   	server.Router.AddRoute("BUY", handlers.HandleBuyPackage)

	   	// Chat routes
	   	server.Router.AddRoute("send_message", handlers.HandleSendMessage)
	   	server.Router.AddRoute("get_messages", handlers.HandleGetMessages) */

	server.Router.AddRoute("ping", handlers.HandlePing)
	server.Router.AddRoute("register", handlers.HandleRegister)
	server.Router.AddRoute("login", handlers.HandleLogin)
	server.Router.AddRoute("create", handlers.HandleCreateRoom)
	server.Router.AddRoute("join", handlers.HandleJoinRoom)
	server.Router.AddRoute("leave", handlers.HandleLeaveRoom)
	server.Router.AddRoute("send", handlers.HandleSendMessage)
	server.Router.AddRoute("get", handlers.HandleGetMessages)
	server.Router.AddRoute("buy", handlers.HandleBuyPackage)
	server.Router.AddRoute("choose", handlers.HandleChooseCard)
	server.Router.AddRoute("get_opponent_choice", handlers.HandleGetOpponentChoice)
	server.Router.AddRoute("quit", handlers.HandleQuitGame)

	// Start server
	if err := server.Start(); err != nil {
		config.Logger.Fatalf("Failed to start server: %v", err)
	}
	config.Logger.Infof("Server started on %s:%s", config.HOST, config.PORT)

	// Wait for a signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	config.Logger.Info("Shutting down server...")
	if err := server.Stop(); err != nil {
		config.Logger.Errorf("Failed to stop server gracefully: %v", err)
	}
}
