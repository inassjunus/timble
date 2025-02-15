package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"timble/internal/utils"
)

func TestError_NewStandardError(t *testing.T) {
	cases := []struct {
		name           string
		message        string
		code           string
		field          string
		expectedResult *utils.StandardError
	}{
		{
			name:    "normal case with all the fields",
			message: "test error",
			code:    "test error code",
			field:   "test error field",
			expectedResult: &utils.StandardError{
				Message: "test error",
				Code:    "test error code",
				Field:   "test error field",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := utils.NewStandardError(tc.message, tc.code, tc.field)

			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

func TestError_StandardError_Error(t *testing.T) {
	cases := []struct {
		name           string
		err            *utils.StandardError
		expectedResult string
	}{
		{
			name:           "normal case with all the fields",
			err:            utils.BadRequestParamError("Missing param", "a param"),
			expectedResult: "Error on\ncode: PARAMETER_PARSING_FAILS; error: Missing param; field: a param",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.err.Error()

			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

func TestError_BadRequestParamError(t *testing.T) {
	cases := []struct {
		name           string
		message        string
		field          string
		expectedResult string
	}{
		{
			name:           "normal case",
			message:        "test err message",
			field:          "test field",
			expectedResult: "Error on\ncode: PARAMETER_PARSING_FAILS; error: test err message; field: test field",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := utils.BadRequestParamError(tc.message, tc.field).Error()

			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

func TestError_UserNotFoundError(t *testing.T) {
	cases := []struct {
		name           string
		userID         uint
		expectedResult string
	}{
		{
			name:           "normal case",
			userID:         uint(1),
			expectedResult: "Error on\ncode: NOT FOUND; error: User not found:1; field:",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := utils.UserNotFoundError(tc.userID).Error()

			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}
