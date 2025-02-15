package utils

import (
	"bytes"
	"fmt"
	"strings"
)

type StandardError struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Field   string `json:"field,omitempty"`
}

var (
	ErrorInternalServerResponse = &StandardError{
		Message: "internal server error, please check the server logs",
		Code:    "INTERNAL SERVER ERROR",
	}

	ErrUnauthenticated = &StandardError{
		Message: "Invalid or missing required authentication",
		Code:    "Unauthorized",
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
		Message: message,
		Code:    "PARAMETER_PARSING_FAILS",
		Field:   field,
	}
}

func UserNotFoundError(userID uint) *StandardError {
	return &StandardError{
		Message: fmt.Sprintf("User not found:%d", userID),
		Code:    "NOT FOUND",
	}
}
