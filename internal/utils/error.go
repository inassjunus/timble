package utils

import (
	"bytes"
	"fmt"
	"strings"
)

type ErrorDetails struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Field   string `json:"field,omitempty"`
}

type StandardError struct {
	details []*ErrorDetails
}

// Error implement error interface
func (e *ErrorDetails) Error() string {
	return e.Message
}

// Error implement error interface
func (s *StandardError) Error() string {
	buff := bytes.NewBufferString("")

	buff.WriteString("Error on\n")
	for _, err := range s.details {
		buff.WriteString("code: ")
		buff.WriteString(err.Code)
		buff.WriteString("; error: ")
		buff.WriteString(err.Error())
		buff.WriteString("; field: ")
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

var (
	ErrorInternalServerResponse = &StandardError{
		details: []*ErrorDetails{
			&ErrorDetails{
				Message: "internal server error, please check the server logs",
				Code:    "INTERNAL SERVER ERROR",
			},
		}}
)

func BadRequestParamError(message, field string) *StandardError {
	return &StandardError{
		details: []*ErrorDetails{
			&ErrorDetails{
				Message: message,
				Code:    "PARAMETER_PARSING_FAILS",
				Field:   field,
			},
		}}
}

func UserNotFoundError(userID uint) *StandardError {
	return &StandardError{
		details: []*ErrorDetails{
			&ErrorDetails{
				Message: fmt.Sprintf("User not found:%d", userID),
				Code:    "NOT FOUND",
			},
		}}
}
