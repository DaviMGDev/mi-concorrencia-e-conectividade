package model 

import (
	"errors"
)

type Chat struct {
	ChatId   string              `json:"chatId"`
	UsersId  map[string]struct{} `json:"usersId"`
	Messages map[string][]string `json:"messages"`
}

func (chat *Chat) AddUserId(userId string) error {
	_, exists := chat.UsersId[userId]
	if exists {
		return errors.New("User already in chat")
	}
	chat.UsersId[userId] = struct{}
	return nil
}

func (chat *Chat) RemoveUserId(userId string) error {
	_, exists := chat.UsersId[userId]
	if !exists {
		return errors.New("User is not in chat")
	}
	delete(chat.UsersId, userId)
	return nil
}

func (chat *Chat) AddMessage(userId string, message string) error {
	_, exists := chat.UsersId[userId]
	if !exists {
		return errors.New("User is not in chat")
	}
	chat.Messages[userId] = append(chat.Messages[userId], message)
	return nil
}

