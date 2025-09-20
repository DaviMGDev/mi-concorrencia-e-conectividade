package handlers

import (
	"server-of-hope/internal/api"
	"server-of-hope/internal/api/protocol"
	"server-of-hope/internal/state"
	"server-of-hope/internal/utils"
)

func HandleCreateRoom(server *api.Server, request protocol.Request) {
	responder := NewResponder(server, request)
	defer responder.Send()

	roomID, err := state.RoomService.CreateRoom()
	if err != nil {
		responder.SetError("Could not create room", "Failed to create room", "from", request.From, "error", err)
		return
	}

	data := utils.Dict{
		"message": "Room created successfully",
		"room_id": roomID,
	}
	responder.SetSuccess(data, "Room created successfully", "from", request.From, "room_id", roomID)
}

func HandleJoinRoom(server *api.Server, request protocol.Request) {
	responder := NewResponder(server, request)
	defer responder.Send()

	roomID, roomIDOk := request.Data["room_id"].(string)
	userID, userIDOk := request.Data["user_id"].(string)

	if !roomIDOk || !userIDOk {
		responder.SetError("Invalid parameters", "Failed to join room", "from", request.From)
		return
	}

	err := state.RoomService.JoinRoom(roomID, userID)
	if err != nil {
		if err.Error() == "A sala está cheia" {
			responder.SetError("A sala está cheia", "Failed to join room", "from", request.From, "room_id", roomID, "error", err)
		} else {
			responder.SetError("Room does not exist", "Failed to join room", "from", request.From, "room_id", roomID, "error", err)
		}
		return
	}

	data := utils.Dict{"message": "Joined room successfully"}
	responder.SetSuccess(data, "Joined room successfully", "from", request.From, "room_id", roomID)
}

func HandleLeaveRoom(server *api.Server, request protocol.Request) {
	responder := NewResponder(server, request)
	defer responder.Send()

	roomID, roomIDOk := request.Data["room_id"].(string)
	userID, userIDOk := request.Data["user_id"].(string)

	if !roomIDOk || !userIDOk {
		responder.SetError("Invalid parameters", "Failed to leave room", "from", request.From)
		return
	}

	err := state.RoomService.LeaveRoom(roomID, userID)
	if err != nil {
		responder.SetError("Room does not exist or user not in room", "Failed to leave room", "from", request.From, "room_id", roomID, "error", err)
		return
	}

	data := utils.Dict{
		"message": "Left room successfully",
		"room_id": roomID,
	}
	responder.SetSuccess(data, "Left room successfully", "from", request.From, "room_id", roomID)
}
