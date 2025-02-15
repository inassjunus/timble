package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	redis "timble/internal/connection/redis"
	mocksredis "timble/mocks/internal_/connection/redis"
	"timble/module/users/internal/repository"
)

var (
	testKey    = "testkey"
	testMember = "testmember"
	testExpire = 5 * time.Millisecond
)

func TestNewRedisRepository(t *testing.T) {
	t.Run("new redis repository", func(t *testing.T) {
		repo := repository.NewRedisRepository(&redis.RedisClient{})

		assert.IsType(t, &repository.RedisRepository{}, repo)
	})
}

func TestRedisRepository_Incr(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name           string
		key            string
		expire         time.Duration
		expectedResult int64
		expectedError  error
		mockRedisCall  func(redisClient *mocksredis.RedisInterface)
	}{
		{
			name:   "normal case - successfully add new key",
			key:    testKey,
			expire: testExpire,
			mockRedisCall: func(redisClient *mocksredis.RedisInterface) {
				redisClient.On("Incr", ctx, testKey).Return(int64(1), nil)
				redisClient.On("Expire", ctx, testKey, testExpire).Return(true, nil)
			},
			expectedResult: int64(1),
		},
		{
			name:   "error case - error when incrementing",
			key:    testKey,
			expire: testExpire,
			mockRedisCall: func(redisClient *mocksredis.RedisInterface) {
				redisClient.On("Incr", ctx, testKey).Return(int64(0), errors.New("timeout"))
			},
			expectedError:  errors.New("redis client error when incr: timeout"),
			expectedResult: int64(0),
		},
	}

	for _, tc := range tests {
		redisClient := mocksredis.NewRedisInterface(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.mockRedisCall(redisClient)
			repo := repository.NewRedisRepository(redisClient)
			res, err := repo.Incr(ctx, tc.key, tc.expire)

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestRedisRepository_Set(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name           string
		key            string
		member         interface{}
		expire         time.Duration
		expectedError  error
		expectedResult string
		mockRedisCall  func(redisClient *mocksredis.RedisInterface)
	}{
		{
			name:           "normal case - successfully add new key",
			key:            testKey,
			member:         testMember,
			expire:         testExpire,
			expectedResult: "OK",
			mockRedisCall: func(redisClient *mocksredis.RedisInterface) {
				redisClient.On("Set", ctx, testKey, testMember, testExpire).Return("OK", nil)
			},
		},
		{
			name:   "error case - error when adding new value",
			key:    testKey,
			member: testMember,
			expire: testExpire,
			mockRedisCall: func(redisClient *mocksredis.RedisInterface) {
				redisClient.On("Set", ctx, testKey, testMember, testExpire).Return("", errors.New("timeout"))
			},
			expectedError: errors.New("redis client error when set: timeout"),
		},
	}

	for _, tc := range tests {
		redisClient := mocksredis.NewRedisInterface(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.mockRedisCall(redisClient)
			repo := repository.NewRedisRepository(redisClient)
			result, err := repo.Set(ctx, tc.key, tc.member, tc.expire)

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestRedisRepository_Get(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name           string
		key            string
		expectedError  error
		expectedResult string
		mockRedisCall  func(redisClient *mocksredis.RedisInterface)
	}{
		{
			name:           "normal case - successfully get key",
			key:            testKey,
			expectedResult: testMember,
			mockRedisCall: func(redisClient *mocksredis.RedisInterface) {
				redisClient.On("Get", ctx, testKey).Return(testMember, nil)
			},
		},
		{
			name: "error case - error when adding new value",
			key:  testKey,
			mockRedisCall: func(redisClient *mocksredis.RedisInterface) {
				redisClient.On("Get", ctx, testKey).Return("", errors.New("timeout"))
			},
			expectedError: errors.New("redis client error when get: timeout"),
		},
	}

	for _, tc := range tests {
		redisClient := mocksredis.NewRedisInterface(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.mockRedisCall(redisClient)
			repo := repository.NewRedisRepository(redisClient)
			result, err := repo.Get(ctx, tc.key)

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
