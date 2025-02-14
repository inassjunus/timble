package utils_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"timble/internal/utils"
)

func TestAuth_GenerateToken(t *testing.T) {
	cases := []struct {
		name           string
		cfg            utils.AuthConfig
		expectedResult string
		expectedError  error
	}{
		{
			name: "successfully create token",
			cfg: utils.AuthConfig{
				SecretKey: []byte("secretz"),
				TokenExp:  time.Hour,
			},
			expectedResult: `[a-zA-Z0-9]+\.[a-zA-Z0-9]+\.[a-zA-Z0-9\-\_]+`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.cfg.GenerateToken(uint(1))
			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				match, _ := regexp.MatchString(tc.expectedResult, result)
				assert.Equal(t, true, match)
			}
		})
	}
}

func TestAuth_VerifyToken(t *testing.T) {
	cfg := utils.AuthConfig{
		SecretKey: []byte("secretz"),
		TokenExp:  time.Hour,
	}
	testToken, _ := cfg.GenerateToken(uint(1))
	cases := []struct {
		name           string
		token          string
		cfg            utils.AuthConfig
		expectedResult float64
		expectedError  error
	}{
		{
			name:           "successfully verify token",
			token:          testToken,
			cfg:            cfg,
			expectedResult: 1,
		},
		{
			name:          "error verifying token",
			token:         "thisisinvalidtoken",
			cfg:           cfg,
			expectedError: errors.New("token is malformed: token contains an invalid number of segments"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.cfg.VerifyToken(tc.token)
			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result["user_id"])
			}
		})
	}
}
