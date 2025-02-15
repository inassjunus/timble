package entity

import (
	"encoding/json"
	"io"
	"net/mail"
	"timble/internal/utils"
	"time"
)

type User struct {
	ID             uint      `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Premium        bool      `json:"premium"`
	HashedPassword string    `json:"hashed_password"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserPublic struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Premium   bool      `json:"premium"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRegistrationParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserToken struct {
	Token string `json:"token"`
}

func NewUserRegistrationPayload(body io.Reader) (UserRegistrationParams, error) {
	params := UserRegistrationParams{}
	err := json.NewDecoder(body).Decode(&params)
	if err != nil {
		return params, err
	}

	if len(params.Username) == 0 {
		return params, utils.BadRequestParamError("Username can not be blank", "username")
	}

	if len(params.Email) == 0 {
		return params, utils.BadRequestParamError("Email can not be blank", "email")
	}

	_, err = mail.ParseAddress(params.Email)

	if err != nil {
		return params, utils.BadRequestParamError("Invalid email format", "email")
	}

	if len(params.Password) <= 10 {
		return params, utils.BadRequestParamError("Password must be more than 10 characters", "password")
	}
	return params, nil
}

func NewUserLoginPayload(body io.Reader) (UserLoginParams, error) {
	params := UserLoginParams{}
	err := json.NewDecoder(body).Decode(&params)
	if err != nil {
		return params, err
	}
	if len(params.Username) == 0 {
		return params, utils.BadRequestParamError("Username can not be blank", "username")
	}

	if len(params.Password) <= 10 {
		return params, utils.BadRequestParamError("Password must be more than 10 characters", "password")
	}
	return params, nil
}
