package main

import (
	"client-of-dispair/internal/api"
	"client-of-dispair/internal/config"
	"client-of-dispair/internal/handlers"
	"client-of-dispair/internal/ui"
)

func main() {
	config.Initialize()
	defer config.Close()

	client := api.NewClient(config.ServerAddress)

	chat := ui.NewChat()

	router := api.NewRouter(client, chat)

	router.AddHandler("/ping", handlers.PingHandler)
	router.AddHandler("/send", handlers.SendMessageHandler)
	router.AddHandler("/get", handlers.GetMessageHandler)
	router.AddHandler("/register", handlers.RegisterHandler)
	router.AddHandler("/login", handlers.LoginHandler)
	router.AddHandler("/logout", handlers.LogoutHandler)
	router.AddHandler("/create", handlers.CreateRoomHandler)
	router.AddHandler("/join", handlers.JoinRoomHandler)
	router.AddHandler("/leave", handlers.LeaveRoomHandler)
	router.AddHandler("/choice", handlers.ChoiceHandler)
	router.AddHandler("/play", handlers.PlayHandler)
	/* 	router.AddHandler("/get_opponent_choice", handlers.GetOpponentChoiceHandler) */
	router.AddHandler("/quit", handlers.QuitHandler)

	err := client.Connect()
	if err != nil {
		chat.Outputs <- "Error connecting to server: " + err.Error()
		return
	}
	defer client.Close()
	chat.Start()
	router.Start()
	select {}
}
