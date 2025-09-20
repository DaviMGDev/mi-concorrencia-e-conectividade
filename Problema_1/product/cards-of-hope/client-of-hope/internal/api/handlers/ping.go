package handlers

import (
	"client-of-hope/internal/api"
	"client-of-hope/internal/api/protocol"
	"client-of-hope/internal/state"
	"client-of-hope/internal/ui"
	"client-of-hope/internal/utils"
	"fmt"
)

func HandlePing(client *api.Client, chat *ui.Chat, args []string) {
       request := protocol.Request{
	       Method: "ping",
	       Data: utils.Dict{
		       "user_id": state.UserID,
	       },
       }

       start := utils.NowMillis()
       response, err := client.DoRequest(request)
       latency := utils.NowMillis() - start
       if err != nil {
	       state.Log("Ping request failed: %v", err)
	       chat.Outputs <- "Ping request failed."
	       return
       }

       chat.Outputs <- fmt.Sprintf("Pong! Server responded: %s | Latência: %d ms", response.Status, latency)
}