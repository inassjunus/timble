package entity_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"timble/module/users/entity"
)

func TestUser_NewUserPayload(t *testing.T) {
	type args struct {
		configName string
		configData string
	}
	tests := []struct {
		name           string
		body           string
		expectedResult entity.UserParams
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
			expectedResult: entity.UserParams{
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
			expectedErr: errors.New("Username can not be blank"),
		},
		{
			name: "error case with missing password",
			body: `
		    {
		      "username":  "testuser"
		    }
		  `,
			expectedErr: errors.New("Password must be more than 10 characters"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := entity.NewUserPayload(strings.NewReader(tc.body))
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, actual)
		})
	}
}
