// TO FIX
// MAYBE UNUSED
package entities

type Message struct {
	UserID    string `json:"from"`
	RoomID    string `json:"room_id"`
	Content   string `json:"content"`
	TimeStamp string `json:"timestamp"`
}

func NewMessage(userID, roomID, content, timestamp string) *Message {
	return &Message{
		UserID:    userID,
		RoomID:    roomID,
		Content:   content,
		TimeStamp: timestamp,
	}
}

type MessageInterface interface {
}
