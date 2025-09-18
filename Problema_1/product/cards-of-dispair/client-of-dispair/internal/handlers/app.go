package handlers

import (
	"client-of-dispair/internal/api"
	"client-of-dispair/internal/ui"
	"os"
)

func QuitHandler(client *api.Client, chat *ui.Chat, arguments []string) {
	chat.Outputs <- "Quitting the application. Goodbye!"
	client.Close()
	os.Exit(0)
}
