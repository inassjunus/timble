package entity

import (
	"encoding/json"
	"errors"
	"io"
)

type User struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	Premium        bool   `json:"premium"`
	HashedPassword string `json:"hashed_password"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type UserPublic struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Premium   bool   `json:"premium"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserToken struct {
	Token string `json:"token"`
}

func NewUserPayload(body io.Reader) (UserParams, error) {
	params := UserParams{}
	err := json.NewDecoder(body).Decode(&params)
	if err != nil {
		return params, err
	}
	if len(params.Username) == 0 {
		return params, errors.New("Username can not be blank")
	}

	if len(params.Password) <= 10 {
		return params, errors.New("Password must be more than 10 characters")
	}
	return params, nil
}
