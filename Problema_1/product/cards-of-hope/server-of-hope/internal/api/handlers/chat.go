package handlers

import (
	"server-of-hope/internal/api"
	"server-of-hope/internal/api/protocol"
	"server-of-hope/internal/state"
	"server-of-hope/internal/utils"
)

func HandleSendMessage(server *api.Server, request protocol.Request) {
	responder := NewResponder(server, request)
	defer responder.Send()

	roomID, roomIDOk := request.Data["room_id"].(string)
	message, messageOk := request.Data["message"].(string)
	userID, userIDOk := request.Data["user_id"].(string)

	if !roomIDOk || !messageOk || !userIDOk {
		responder.SetError("Invalid parameters", "Failed to send message", "from", request.From)
		return
	}

	err := state.ChatService.SendMessage(roomID, userID, message)
	if err != nil {
		responder.SetError("Room does not exist or user not in room", "Failed to send message", "from", request.From, "room_id", roomID, "error", err)
		return
	}

	data := utils.Dict{
		"message": "Message sent successfully",
		"room_id": roomID,
	}
	responder.SetSuccess(data, "Message sent successfully", "from", request.From, "room_id", roomID)
}

func HandleFetchMessage(server *api.Server, request protocol.Request) {
	responder := NewResponder(server, request)
	defer responder.Send()

	roomID, roomIDOk := request.Data["room_id"].(string)
	userID, userIDOk := request.Data["user_id"].(string)

	if !roomIDOk || !userIDOk {
		responder.SetError("Invalid parameters", "Failed to fetch message", "from", request.From)
		return
	}

	message, err := state.ChatService.ReceiveMessage(roomID, userID)
	if err != nil {
		responder.SetError("Room does not exist or user not in room", "Failed to fetch message", "from", request.From, "room_id", roomID, "user_id", userID, "error", err)
		return
	}

	data := utils.Dict{
		"message": message,
		"room_id": roomID,
	}
	responder.SetSuccess(data, "Message fetched successfully", "from", request.From, "room_id", roomID, "user_id", userID)
}
