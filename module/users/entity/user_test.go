package entity_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"timble/module/users/entity"
)

func TestUser_NewUserRegistrationPayload(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedResult entity.UserRegistrationParams
		expectedErr    error
	}{
		{
			name: "normal case",
			body: `
		    {
		      "username":  "testuser",
		      "email": "test@email.com",
		      "password": "testpassword"
		    }
		  `,
			expectedResult: entity.UserRegistrationParams{
				Username: "testuser",
				Email:    "test@email.com",
				Password: "testpassword",
			},
		},
		{
			name: "error case with invalid payload",
			body: `
		    {
		      "email": "test@email.com",
		      "password": "testpassword"`,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: unexpected EOF; field: payload"),
		},
		{
			name: "error case with missing username",
			body: `
		    {
		      "email": "test@email.com",
		      "password": "testpassword"
		    }
		  `,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: Username can not be blank; field: username"),
		},
		{
			name: "error case with missing email",
			body: `
		    {
		      "username":  "testuser",
		      "password": "testpassword"
		    }
		  `,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: Email can not be blank; field: email"),
		},
		{
			name: "error case with invalid email",
			body: `
		    {
		      "username":  "testuser",
		      "email": "testemailcom",
		      "password": "testpassword"
		    }
		  `,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: Invalid email format; field: email"),
		},
		{
			name: "error case with missing password",
			body: `
		    {
		      "username":  "testuser",
		      "email": "test@email.com"
		    }
		  `,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: Password must be more than 10 characters; field: password"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := entity.NewUserRegistrationPayload(strings.NewReader(tc.body))
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}

func TestUser_NewUserLoginPayload(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedResult entity.UserLoginParams
		expectedErr    error
	}{
		{
			name: "normal case",
			body: `
		    {
		      "username":  "testuser",
		      "password": "testpassword"
		    }
		  `,
			expectedResult: entity.UserLoginParams{
				Username: "testuser",
				Password: "testpassword",
			},
		},
		{
			name: "error case with missing username",
			body: `
		    {
		      "password": "testpassword"
		    }
		  `,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: Username can not be blank; field: username"),
		},
		{
			name: "error case with invalid payload",
			body: `
		    {
		      "password": "testpassword" `,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: unexpected EOF; field: payload"),
		},
		{
			name: "error case with missing password",
			body: `
		    {
		      "username":  "testuser"
		    }
		  `,
			expectedErr: errors.New("Error on\ncode: PARAMETER_PARSING_FAILS; error: Password must be more than 10 characters; field: password"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := entity.NewUserLoginPayload(strings.NewReader(tc.body))
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}
