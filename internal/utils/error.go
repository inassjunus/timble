package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type StandardError struct {
	Message    string `json:"message"`
	Code       string `json:"code,omitempty"`
	Field      string `json:"field,omitempty"`
	HttpStatus int    `json:"-"`
}

var (
	ErrorInternalServerResponse = &StandardError{
		Message: "internal server error, please check the server logs",
		Code:    "INTERNAL SERVER ERROR",
	}

	ErrUnauthenticated = &StandardError{
		Message:    "Invalid or missing required authentication",
		Code:       "Unauthorized",
		HttpStatus: 401,
	}

	ErrInvalidLogin = &StandardError{
		Message:    "Invalid username or password",
		Code:       "Unauthorized",
		HttpStatus: http.StatusUnauthorized,
	}

	ErrDuplicateUser = &StandardError{
		Message:    "Username or email already exists",
		Code:       "DUPLICATE_USER",
		HttpStatus: http.StatusBadRequest,
	}
)

func NewStandardError(message, code, field string) *StandardError {
	return &StandardError{
		Message: message,
		Code:    code,
		Field:   field,
	}
}

// Error implement error interface
func (s *StandardError) Error() string {
	buff := bytes.NewBufferString("")

	buff.WriteString("Error on\n")
	buff.WriteString("code: ")
	buff.WriteString(s.Code)
	buff.WriteString("; error: ")
	buff.WriteString(s.Message)
	buff.WriteString("; field: ")
	buff.WriteString(s.Field)
	buff.WriteString("\n")

	return strings.TrimSpace(buff.String())
}

func BadRequestParamError(message, field string) *StandardError {
	return &StandardError{
		Message:    message,
		Code:       "PARAMETER_PARSING_FAILS",
		Field:      field,
		HttpStatus: http.StatusBadRequest,
	}
}

func UserNotFoundError(userID uint) *StandardError {
	return &StandardError{
		Message: fmt.Sprintf("User not found:%d", userID),
		Code:    "NOT FOUND",
	}
}
