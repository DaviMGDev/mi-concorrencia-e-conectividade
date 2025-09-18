package handlers

import (
	"server-of-dispair/internal/config"
	"server-of-dispair/internal/entities"
	"server-of-dispair/internal/protocol"
	"time"
)

func HandleSendMessage(server *protocol.Server, request *protocol.Request) {
	userID, userIDOk := request.Data["user_id"].(string)
	roomID, roomIDOk := request.Data["room_id"].(string)
	content, contentOk := request.Data["message"].(string)

	if !userIDOk || !roomIDOk || !contentOk {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data", nil)
		server.Responses <- response
		return
	}

	// Use the server timestamp for consistency
	timestamp := time.Now().Format(time.RFC3339)

	msg := entities.NewMessage(userID, roomID, content, timestamp)

	err := config.ChatService.SendMessage(msg)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	response := protocol.NewResponse(request.From, request.Method, "success", "Message sent successfully", nil)
	// The server should broadcast this message to all members of the room.
	// This logic is not yet implemented.
	server.Responses <- response
}

func HandleGetMessages(server *protocol.Server, request *protocol.Request) {
	roomID, roomIDOk := request.Data["room_id"].(string)
	sinceStr, sinceOk := request.Data["since"].(string)

	if !roomIDOk || !sinceOk {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data", nil)
		server.Responses <- response
		return
	}

	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid 'since' timestamp format. Use RFC3339.", nil)
		server.Responses <- response
		return
	}

	messages, err := config.ChatService.GetMessages(roomID, since)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	responseData := map[string]any{"messages": messages}
	response := protocol.NewResponse(request.From, request.Method, "success", "Messages retrieved successfully", responseData)
	server.Responses <- response
}
