package models

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	GUID     string `json:"guid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser(username string, password string) (User, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("cannot create user: %w", err)
	}

	return User{
		GUID:     uuid.NewString(),
		Username: username,
		Password: string(hashPassword),
	}, nil
}

func (u User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
