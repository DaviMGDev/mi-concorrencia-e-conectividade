// TO CHECK

package entities

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	HashPassword string `json:"hash_password"`
}

func NewUser(id, username, password string) (*User, error) {
	user := &User{
		ID:       id,
		Username: username,
	}
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	return user, nil
}

type UserInterface interface {
	SetPassword(password string) error
	CheckPassword(password string) bool
	GetID() string
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashPassword = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashPassword), []byte(password))
	return err == nil
}
