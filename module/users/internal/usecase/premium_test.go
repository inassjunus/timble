package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	log "go.uber.org/zap"

	mocksrepo "timble/mocks/module/users/internal_/usecase"
	"timble/module/users/entity"
	"timble/module/users/internal/repository"
	uc "timble/module/users/internal/usecase"
)

func TestNewPremiumUsecase(t *testing.T) {
	t.Run("new premiumh usecase", func(t *testing.T) {

		usecase := uc.NewPremiumUsecase(
			&repository.RedisRepository{},
			&repository.PostgresRepository{},
			&repository.CacheRepository{},
			&log.Logger{},
		)

		assert.IsType(t, &uc.PremiumUc{}, usecase)
	})
}

func TestPremiumUc_Grant(t *testing.T) {
	type args struct {
		params   uint
		dbParams entity.User
	}

	type mocked struct {
		redisResult string
		redisError  error
		dbError     error
	}
	tests := []struct {
		name        string
		args        args
		mocked      mocked
		expectedErr error
	}{
		{
			name: "normal case - successfully grant premium",
			args: args{
				params: 1,
				dbParams: entity.User{
					ID:      1,
					Premium: true,
				},
			},
			mocked: mocked{
				redisResult: uc.PREMIUM_TRUE_STRING,
			},
		},
		{
			name: "error case - user not eligible",
			args: args{
				params: 1,
			},
			expectedErr: errors.New("Error on\ncode: NOT ELIGIBLE FOR PREMIUM; error: You are not eligible for premium for now; field:"),
		},
		{
			name: "error case - error from redis",
			args: args{
				params: 1,
				dbParams: entity.User{
					ID:      1,
					Premium: true,
				},
			},
			mocked: mocked{
				redisError: errors.New("Redis error"),
			},
			expectedErr: errors.New("Redis error"),
		},
		{
			name: "error case - error from db",
			args: args{
				params: 1,
				dbParams: entity.User{
					ID:      1,
					Premium: true,
				},
			},
			mocked: mocked{
				dbError:     errors.New("DB error"),
				redisResult: uc.PREMIUM_TRUE_STRING,
			},
			expectedErr: errors.New("DB error"),
		},
	}
	for _, tc := range tests {
		db := mocksrepo.NewPostgresRepository(t)
		redis := mocksrepo.NewRedisRepository(t)
		cache := mocksrepo.NewCacheRepository(t)
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			redis.On("Get", ctx, "eligible_for_premium:1").Return(tc.mocked.redisResult, tc.mocked.redisError)
			if tc.mocked.redisResult == "true" {
				db.On("UpdateUserPremium", tc.args.dbParams, interface{}(tc.args.dbParams.Premium)).Return(tc.mocked.dbError)
			}

			if tc.expectedErr == nil {
				cache.On("Set", ctx, "premium:1", []byte("true"), 24*time.Hour).Return(nil)
				redis.On("Del", ctx, "eligible_for_premium:1").Return(int64(1), nil)
			}

			usecase := uc.NewPremiumUsecase(redis, db, cache, &log.Logger{})

			err := usecase.Grant(ctx, tc.args.params)
			if tc.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestPremiumUc_Unsubscribe(t *testing.T) {
	type args struct {
		params   uint
		dbParams entity.User
	}

	type mocked struct {
		dbError error
	}
	tests := []struct {
		name        string
		args        args
		mocked      mocked
		expectedErr error
	}{
		{
			name: "normal case - successfully unsubscribe premium",
			args: args{
				params: 1,
				dbParams: entity.User{
					ID:      1,
					Premium: false,
				},
			},
		},
		{
			name: "error case - error from db",
			args: args{
				params: 1,
				dbParams: entity.User{
					ID:      1,
					Premium: false,
				},
			},
			mocked: mocked{
				dbError: errors.New("DB error"),
			},
			expectedErr: errors.New("DB error"),
		},
	}
	for _, tc := range tests {
		db := mocksrepo.NewPostgresRepository(t)
		redis := mocksrepo.NewRedisRepository(t)
		cache := mocksrepo.NewCacheRepository(t)
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db.On("UpdateUserPremium", tc.args.dbParams, interface{}(tc.args.dbParams.Premium)).Return(tc.mocked.dbError)
			if tc.expectedErr == nil {
				cache.On("Set", ctx, "premium:1", []byte("false"), 24*time.Hour).Return(nil)
			}

			usecase := uc.NewPremiumUsecase(redis, db, cache, &log.Logger{})

			err := usecase.Unsubscribe(ctx, tc.args.params)
			if tc.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
