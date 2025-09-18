package handlers

import (
	"server-of-dispair/internal/config"
	"server-of-dispair/internal/protocol"
)

func HandleCreateRoom(server *protocol.Server, request *protocol.Request) {
	roomID, err := config.RoomService.CreateRoom()
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	responseData := map[string]any{"room_id": roomID}
	response := protocol.NewResponse(request.From, request.Method, "success", "Room created successfully", responseData)
	server.Responses <- response
}

func HandleJoinRoom(server *protocol.Server, request *protocol.Request) {
	roomID, roomIDOk := request.Data["room_id"].(string)
	playerID, playerIDOk := request.Data["player_id"].(string)

	if !roomIDOk || !playerIDOk {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data: missing room_id or player_id", nil)
		server.Responses <- response
		return
	}

	err := config.RoomService.JoinRoom(roomID, playerID)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	response := protocol.NewResponse(request.From, request.Method, "success", "Joined room successfully", nil)
	server.Responses <- response
}

func HandleLeaveRoom(server *protocol.Server, request *protocol.Request) {
	roomID, roomIDOk := request.Data["room_id"].(string)
	playerID, playerIDOk := request.Data["player_id"].(string)

	if !roomIDOk || !playerIDOk {
		response := protocol.NewResponse(request.From, request.Method, "error", "Invalid request data: missing room_id or player_id", nil)
		server.Responses <- response
		return
	}

	err := config.RoomService.LeaveRoom(roomID, playerID)
	if err != nil {
		response := protocol.NewResponse(request.From, request.Method, "error", err.Error(), nil)
		server.Responses <- response
		return
	}

	response := protocol.NewResponse(request.From, request.Method, "success", "Left room successfully", nil)
	server.Responses <- response
}