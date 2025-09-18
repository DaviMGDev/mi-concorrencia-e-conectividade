package entities

import (
	"server-of-dispair/internal/utils"
	"sync"
)

type Room struct {
	ID      string                  `json:"id"`
	Members *utils.Map[string, User] `json:"members"`

	mutex sync.RWMutex `json:"-"`
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Members: utils.NewMap[string, User](),
	}
}

type RoomInterface interface {
	AddMember(user User) error
	RemoveMember(userID string) error
	ListMembers() []User
}

func (r *Room) AddMember(user User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.Members.Set(user.ID, user)
	return nil
}

func (r *Room) RemoveMember(userID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.Members.Delete(userID)
	return nil
}

func (r *Room) ListMembers() []User {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.Members.Values()
}