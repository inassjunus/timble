package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"timble/internal/connection/cache"
	mockscache "timble/mocks/internal_/connection/cache"
	"timble/module/users/internal/repository"
)

var (
	testCacheKey    = "testkey"
	testCacheValue  = []byte("testmember")
	testCacheExpire = 5 * time.Millisecond
)

func TestNewCacheRepository(t *testing.T) {
	t.Run("new cache repository", func(t *testing.T) {
		repo := repository.NewCacheRepository(&cache.CacheClient{})

		assert.IsType(t, &repository.CacheRepository{}, repo)
	})
}

func TestCacheRepository_Set(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name          string
		key           string
		member        []byte
		expire        time.Duration
		expectedError error
		mockCacheCall func(cacheClient *mockscache.CacheInterface)
	}{
		{
			name:   "normal case - successfully add new key",
			key:    testCacheKey,
			member: testCacheValue,
			expire: testCacheExpire,
			mockCacheCall: func(cacheClient *mockscache.CacheInterface) {
				cacheClient.On("Set", ctx, testCacheKey, testCacheValue, testCacheExpire).Return(nil)
			},
		},
		{
			name:   "error case - error when adding new value",
			key:    testCacheKey,
			member: testCacheValue,
			expire: testCacheExpire,
			mockCacheCall: func(cacheClient *mockscache.CacheInterface) {
				cacheClient.On("Set", ctx, testCacheKey, testCacheValue, testCacheExpire).Return(errors.New("timeout"))
			},
			expectedError: errors.New("cache client error when set: timeout"),
		},
	}

	for _, tc := range tests {
		cacheClient := mockscache.NewCacheInterface(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.mockCacheCall(cacheClient)
			repo := repository.NewCacheRepository(cacheClient)
			err := repo.Set(ctx, tc.key, tc.member, tc.expire)

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestCacheRepository_Get(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name           string
		key            string
		expectedError  error
		expectedResult []byte
		mockCacheCall  func(cacheClient *mockscache.CacheInterface)
	}{
		{
			name:           "normal case - successfully get key",
			key:            testCacheKey,
			expectedResult: testCacheValue,
			mockCacheCall: func(cacheClient *mockscache.CacheInterface) {
				cacheClient.On("Get", ctx, testCacheKey).Return(testCacheValue, nil)
			},
		},
		{
			name: "error case - error when adding new value",
			key:  testCacheKey,
			mockCacheCall: func(cacheClient *mockscache.CacheInterface) {
				cacheClient.On("Get", ctx, testCacheKey).Return([]byte(""), errors.New("timeout"))
			},
			expectedError: errors.New("cache client error when get: timeout"),
		},
	}

	for _, tc := range tests {
		cacheClient := mockscache.NewCacheInterface(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.mockCacheCall(cacheClient)
			repo := repository.NewCacheRepository(cacheClient)
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
