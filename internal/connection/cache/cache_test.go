package cache_test

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"

	cache "timble/internal/connection/cache"
	redis "timble/internal/connection/redis"
)

var (
	testKey1         = "testkey1"
	testKey2         = "testkey2"
	testMember1      = "testmember1"
	testMember2      = "testmember2"
	testField1       = "testfield1"
	testField2       = "testfield2"
	testScore        = float64(100)
	testRedisTimeout = "200ms"
)

func TestCacheClient_NewClient(t *testing.T) {
	s := miniredis.RunT(t)

	tests := []struct {
		name          string
		redisHost     string
		redisPort     string
		redisTimeout  string
		redisDB       string
		expectedError error
	}{
		{
			name:         "normal case",
			redisHost:    s.Host(),
			redisPort:    s.Port(),
			redisDB:      "0",
			redisTimeout: testRedisTimeout,
		},
		{
			name:          "error host case",
			redisHost:     ":/:",
			redisPort:     "123",
			redisDB:       "0",
			redisTimeout:  testRedisTimeout,
			expectedError: errors.New("dial tcp: address :/::123: too many colons in address"),
		},
		{
			name:         "error timeout format case",
			redisHost:    s.Host(),
			redisPort:    s.Port(),
			redisDB:      "0",
			redisTimeout: "abc",
		},
		{
			name:         "error db format case",
			redisHost:    s.Host(),
			redisPort:    s.Port(),
			redisDB:      "abc",
			redisTimeout: testRedisTimeout,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := cache.NewClient(tc.redisHost, tc.redisPort, tc.redisTimeout, tc.redisDB)
			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.NotNil(t, result)
				assert.Nil(t, err)
			}
		})
	}
	defer s.Close()
}

func TestCacheClient_Set(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		value         string
		mockErr       string
		expectedError error
	}{
		{
			name:  "normal case with new key",
			key:   testKey2,
			value: testMember2,
		},
		{
			name:  "normal case with existing key",
			key:   testKey1,
			value: testMember2,
		},
		{
			name:          "error case",
			key:           testKey1,
			value:         testMember2,
			expectedError: errors.New("timeout"),
			mockErr:       "timeout",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := miniredis.RunT(t)
			s.Set(testKey1, testMember1)

			client, _ := cache.NewClient(s.Host(), s.Port(), testRedisTimeout, "0")
			s.SetError(tc.mockErr)

			err := client.Set(context.Background(), tc.key, tc.value, 0)
			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				val, _ := s.Get(tc.key)
				assert.Equal(t, tc.value, val)
			}

			defer s.Close()
		})
	}
}

func TestRedisClient_Get(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		mockErr        string
		expectedResult string
		expectedError  error
	}{
		{
			name:           "normal case with non existing key",
			key:            testKey2,
			expectedResult: "",
		},
		{
			name:           "normal case with existing key",
			key:            testKey1,
			expectedResult: testMember1,
		},
		{
			name:           "error case",
			key:            testKey1,
			expectedResult: "",
			expectedError:  errors.New("timeout"),
			mockErr:        "timeout",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := miniredis.RunT(t)
			s.Set(testKey1, testMember1)

			client, _ := redis.NewClient(s.Host(), s.Port(), testRedisTimeout, "0")
			s.SetError(tc.mockErr)

			result, err := client.Get(context.Background(), tc.key)
			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}

			defer s.Close()
		})
	}
}
