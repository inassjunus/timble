package redis_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"

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

func TestRedisClient_NewClient(t *testing.T) {
	s := miniredis.RunT(t)

	tests := []struct {
		name            string
		redisHost       string
		redisPort       string
		redisTimeout    string
		redisDB         string
		expectedAddr    string
		expectedDB      int
		expectedTimeout time.Duration
		expectedErr     error
	}{
		{
			name:            "normal case",
			redisHost:       s.Host(),
			redisPort:       s.Port(),
			redisDB:         "0",
			expectedAddr:    fmt.Sprintf("%s:%s", s.Host(), s.Port()),
			expectedDB:      0,
			expectedTimeout: 200 * time.Millisecond,
			redisTimeout:    testRedisTimeout,
		},
		{
			name:         "error host case",
			redisHost:    ":/:",
			redisPort:    "123",
			redisDB:      "0",
			redisTimeout: testRedisTimeout,
			expectedErr:  errors.New("dial tcp: address :/::123: too many colons in address"),
		},
		{
			name:            "error timeout format case",
			redisHost:       s.Host(),
			redisPort:       s.Port(),
			redisDB:         "0",
			expectedAddr:    fmt.Sprintf("%s:%s", s.Host(), s.Port()),
			expectedDB:      0,
			expectedTimeout: 100 * time.Millisecond,
			redisTimeout:    "abc",
		},
		{
			name:            "error db format case",
			redisHost:       s.Host(),
			redisPort:       s.Port(),
			redisDB:         "abc",
			expectedAddr:    fmt.Sprintf("%s:%s", s.Host(), s.Port()),
			expectedDB:      0,
			expectedTimeout: 200 * time.Millisecond,
			redisTimeout:    testRedisTimeout,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := redis.NewClient(tc.redisHost, tc.redisPort, tc.redisTimeout, tc.redisDB)
			if tc.expectedErr != nil {
				assert.NotEqual(t, err, nil)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NotEqual(t, nil, result)
				assert.Equal(t, nil, err)
				assert.Equal(t, tc.expectedAddr, result.Client.Options().Addr)
				assert.Equal(t, tc.expectedDB, result.Client.Options().DB)
				assert.Equal(t, tc.expectedTimeout, result.Client.Options().ReadTimeout)
				assert.Equal(t, tc.expectedTimeout, result.Client.Options().WriteTimeout)
			}
		})
	}
	defer s.Close()
}

func TestRedisClient_Incr(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		mockErr        string
		expectedValue  string
		expectedResult int64
		expectedErr    error
	}{
		{
			name:           "normal case with new key",
			key:            testKey2,
			expectedValue:  "1",
			expectedResult: 1,
		},
		{
			name:           "normal case with existing key",
			key:            testKey1,
			expectedValue:  "3",
			expectedResult: 3,
		},
		{
			name:           "error case",
			key:            testKey1,
			expectedResult: 0,
			expectedErr:    errors.New("timeout"),
			mockErr:        "timeout",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := miniredis.RunT(t)
			s.Set(testKey1, "2")

			client, _ := redis.NewClient(s.Host(), s.Port(), testRedisTimeout, "0")
			s.SetError(tc.mockErr)

			result, err := client.Incr(context.Background(), tc.key)
			if tc.expectedErr != nil {
				assert.NotEqual(t, err, nil)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, nil, err)
				assert.Equal(t, tc.expectedResult, result)
				val, _ := s.Get(tc.key)
				assert.Equal(t, tc.expectedValue, val)
			}

			defer s.Close()
		})
	}
}

func TestRedisClient_ExpireAt(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		exp            time.Time
		mockErr        string
		expectedValue  string
		expectedResult bool
		expectedErr    error
	}{
		{
			name:           "normal case with new key",
			key:            testKey2,
			exp:            time.Now().Add(time.Minute),
			expectedValue:  "",
			expectedResult: false,
		},
		{
			name:           "normal case with existing key",
			key:            testKey1,
			exp:            time.Now().Add(time.Minute),
			expectedValue:  testMember1,
			expectedResult: true,
		},
		{
			name:           "error case",
			key:            testKey1,
			exp:            time.Now(),
			expectedResult: false,
			expectedErr:    errors.New("timeout"),
			mockErr:        "timeout",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := miniredis.RunT(t)
			s.Set(testKey1, testMember1)

			client, _ := redis.NewClient(s.Host(), s.Port(), testRedisTimeout, "0")
			s.SetError(tc.mockErr)

			result, err := client.ExpireAt(context.Background(), tc.key, tc.exp)
			if tc.expectedErr != nil {
				assert.NotEqual(t, err, nil)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, nil, err)
				assert.Equal(t, tc.expectedResult, result)
				val, _ := s.Get(tc.key)
				assert.Equal(t, tc.expectedValue, val)
			}

			defer s.Close()
		})
	}
}

func TestRedisClient_Set(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		value          string
		mockErr        string
		expectedResult string
		expectedErr    error
	}{
		{
			name:           "normal case with new key",
			key:            testKey2,
			value:          testMember2,
			expectedResult: "OK",
		},
		{
			name:           "normal case with existing key",
			key:            testKey1,
			value:          testMember2,
			expectedResult: "OK",
		},
		{
			name:           "error case",
			key:            testKey1,
			value:          testMember2,
			expectedResult: "",
			expectedErr:    errors.New("timeout"),
			mockErr:        "timeout",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := miniredis.RunT(t)
			s.Set(testKey1, testMember1)

			client, _ := redis.NewClient(s.Host(), s.Port(), testRedisTimeout, "0")
			s.SetError(tc.mockErr)

			result, err := client.Set(context.Background(), tc.key, tc.value, 0)
			if tc.expectedErr != nil {
				assert.NotEqual(t, err, nil)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, nil, err)
				assert.Equal(t, tc.expectedResult, result)
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
		expectedErr    error
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
			expectedErr:    errors.New("timeout"),
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
			if tc.expectedErr != nil {
				assert.NotEqual(t, err, nil)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, nil, err)
				assert.Equal(t, tc.expectedResult, result)
			}

			defer s.Close()
		})
	}
}
