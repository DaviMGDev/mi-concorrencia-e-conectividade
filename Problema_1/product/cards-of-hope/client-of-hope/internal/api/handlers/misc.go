package handlers

import (
	"client-of-hope/internal/api"
	"client-of-hope/internal/state"
	"client-of-hope/internal/ui"
	"fmt"
)

func HandleWhoami(client *api.Client, chat *ui.Chat, args []string) {
	if state.Username == "" || state.UserID == "" {
		chat.Outputs <- "You are not logged in."
		return
	}
	chat.Outputs <- fmt.Sprintf("You are logged in as %s (ID: %s)", state.Username, state.UserID)
}

func HandleWhereami(client *api.Client, chat *ui.Chat, args []string) {
	if state.RoomID == "" {
		chat.Outputs <- "You are not currently in a room."
		return
	}
	chat.Outputs <- fmt.Sprintf("You are in room: %s", state.RoomID)
}

func HandleQuit(client *api.Client, chat *ui.Chat, args []string) {
	chat.Outputs <- "Disconnecting..."
	chat.Close()
	client.Close()
}

func HandleHelp(client *api.Client, chat *ui.Chat, args []string) {
	helpMessage := `Note: Any text that does not start with / is treated as a /send command.

  Available commands:

  Auth:
    /register <user> <pass>  - Register a new user.
    /login <user> <pass>     - Log in with an existing user.
    /logout                  - Log out from the current session.

  Chat & Rooms:
    /send <message>          - Send a message to the current room.
    /create <room_name>      - Create a new chat room.
    /join <room_name>        - Join an existing chat room.
    /leave                   - Leave the current room.

  Game:
    /play <card>             - Play a card (rock, paper, scissors).
    /cards                   - Show your current cards.
    /buy                     - Buy a new package of cards.

  Misc:
    /whoami                  - Show your current user information.
    /whereami                - Show the room you are currently in.
    /ping                    - Check the connection with the server.
    /help                    - Show this help message.
    /quit                    - Exit the application.`
	chat.Outputs <- helpMessage
}
