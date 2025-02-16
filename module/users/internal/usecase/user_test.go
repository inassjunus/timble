package usecase_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	log "go.uber.org/zap"

	"timble/internal/utils"
	mocksrepo "timble/mocks/module/users/internal_/usecase"
	"timble/module/users/entity"
	"timble/module/users/internal/repository"
	uc "timble/module/users/internal/usecase"
)

var (
	testUser = &entity.User{
		ID:             uint(1),
		Email:          "test@email.com",
		Username:       "testuser",
		Premium:        false,
		HashedPassword: "$2a$14$HLqpimP54B8ujZBmWfpawuphlT3PJs1KebGV.ArukdOp9hHAcOfs2",
	}

	defaultAuthConfig = &utils.AuthConfig{
		SecretKey: []byte("secretz"),
		TokenExp:  time.Hour,
	}
)

func TestNewUserUsecase(t *testing.T) {
	t.Run("new user usecase", func(t *testing.T) {

		usecase := uc.NewUserUsecase(
			&utils.AuthConfig{},
			&repository.RedisRepository{},
			&repository.PostgresRepository{},
			&repository.CacheRepository{},
			&log.Logger{},
		)

		assert.IsType(t, &uc.UserUc{}, usecase)
	})
}

func TestUserUc_Create(t *testing.T) {
	type args struct {
		params   entity.UserRegistrationParams
		dbParams entity.User
	}

	type mocked struct {
		dbInsertError error
		dbGetResult   *entity.User
		dbGetError    error
	}
	tests := []struct {
		name           string
		args           args
		mocked         mocked
		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case - successfully create user",
			args: args{
				params: entity.UserRegistrationParams{
					Username: testUser.Username,
					Email:    testUser.Email,
					Password: "testpassword",
				},
				dbParams: entity.User{
					Username:       testUser.Username,
					Email:          testUser.Email,
					HashedPassword: testUser.HashedPassword,
				},
			},
			mocked: mocked{
				dbGetResult: testUser,
			},
			expectedResult: `[a-zA-Z0-9]+\.[a-zA-Z0-9]+\.[a-zA-Z0-9\-\_]+`,
		},
		{
			name: "error case - error during insert",
			args: args{
				params: entity.UserRegistrationParams{
					Username: testUser.Username,
					Email:    testUser.Email,
					Password: "testpassword",
				},
				dbParams: entity.User{
					Username:       testUser.Username,
					Email:          testUser.Email,
					HashedPassword: testUser.HashedPassword,
				},
			},
			mocked: mocked{
				dbInsertError: errors.New("Error from db insert"),
			},
			expectedErr: errors.New("Error from db insert"),
		},
		{
			name: "error case - error during get",
			args: args{
				params: entity.UserRegistrationParams{
					Username: testUser.Username,
					Email:    testUser.Email,
					Password: "testpassword",
				},
				dbParams: entity.User{
					Username:       testUser.Username,
					Email:          testUser.Email,
					HashedPassword: testUser.HashedPassword,
				},
			},
			mocked: mocked{
				dbGetError: errors.New("Error from db get"),
			},
			expectedErr: errors.New("Error from db get"),
		},
	}
	for _, tc := range tests {
		db := mocksrepo.NewPostgresRepository(t)
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db.On("InsertUser", mock.Anything).Return(tc.mocked.dbInsertError)
			if tc.mocked.dbInsertError == nil {
				db.On("GetUserByUsername", tc.args.params.Username).Return(tc.mocked.dbGetResult, tc.mocked.dbGetError)
			}

			usecase := uc.NewUserUsecase(defaultAuthConfig, &repository.RedisRepository{}, db, &repository.CacheRepository{}, &log.Logger{})

			result, err := usecase.Create(ctx, tc.args.params)
			if tc.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				match, _ := regexp.MatchString(tc.expectedResult, result.Token)
				assert.Equal(t, true, match)
			}
		})
	}
}

func TestUserUc_Show(t *testing.T) {
	timestamp, _ := time.Parse("1/2/2006", "2/2/2025")
	userPublic := &entity.UserPublic{
		ID:        testUser.ID,
		Email:     testUser.Email,
		Username:  testUser.Username,
		Premium:   testUser.Premium,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}
	userData := &entity.User{
		ID:             testUser.ID,
		Email:          testUser.Email,
		Username:       testUser.Username,
		Premium:        testUser.Premium,
		HashedPassword: testUser.HashedPassword,
		CreatedAt:      timestamp,
		UpdatedAt:      timestamp,
	}
	type args struct {
		params uint
	}

	type mocked struct {
		dbGetResult *entity.User
		dbGetError  error
	}
	tests := []struct {
		name           string
		args           args
		mocked         mocked
		expectedResult *entity.UserPublic
		expectedErr    error
	}{
		{
			name: "normal case - successfully show user",
			args: args{
				params: 1,
			},
			mocked: mocked{
				dbGetResult: userData,
			},
			expectedResult: userPublic,
		},
		{
			name: "error case - error during get",
			args: args{
				params: 1,
			},
			mocked: mocked{
				dbGetError: errors.New("Error from db get"),
			},
			expectedErr: errors.New("Error from db get"),
		},
	}
	for _, tc := range tests {
		db := mocksrepo.NewPostgresRepository(t)
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db.On("GetUserByID", tc.args.params).Return(tc.mocked.dbGetResult, tc.mocked.dbGetError)

			usecase := uc.NewUserUsecase(defaultAuthConfig, &repository.RedisRepository{}, db, &repository.CacheRepository{}, &log.Logger{})

			result, err := usecase.Show(ctx, tc.args.params)
			if tc.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestUserUc_React(t *testing.T) {
	reactionParams := entity.ReactionParams{
		UserID:   1,
		TargetID: 2,
		Type:     1,
	}

	testUserPremium := &entity.User{
		ID:             uint(2),
		Email:          "test@email.com",
		Username:       "testuser",
		Premium:        true,
		HashedPassword: "$2a$14$HLqpimP54B8ujZBmWfpawuphlT3PJs1KebGV.ArukdOp9hHAcOfs2",
	}

	type args struct {
		params entity.ReactionParams
	}

	type shouldMock struct {
		dbGetUserByID        bool
		redisGetLimit        bool
		dbGetUserByIDTarget  bool
		dbUpsertUserReaction bool
		redisIncr            bool
	}

	type mocked struct {
		cacheGetPremiumResult     []byte
		cacheGetPremiumError      error
		dbGetUserByIDResult       *entity.User
		dbGetUserByIDError        error
		redisGetLimitResult       string
		dbGetUserByIDTargetResult *entity.User
		dbGetUserByIDTargetError  error
		dbUpsertUserReactionError error
	}
	tests := []struct {
		name           string
		args           args
		shouldMock     shouldMock
		mocked         mocked
		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case - successfully add reaction for non premium user, with premium data in cache",
			args: args{
				params: reactionParams,
			},
			shouldMock: shouldMock{
				redisGetLimit:        true,
				dbGetUserByIDTarget:  true,
				dbUpsertUserReaction: true,
				redisIncr:            true,
			},
			mocked: mocked{
				cacheGetPremiumResult:     []byte("false"),
				redisGetLimitResult:       "1",
				dbGetUserByIDTargetResult: testUserPremium,
			},
		},
		{
			name: "normal case - successfully add reaction for non premium user, with premium data NOT in cache",
			args: args{
				params: reactionParams,
			},
			shouldMock: shouldMock{
				dbGetUserByID:        true,
				redisGetLimit:        true,
				dbGetUserByIDTarget:  true,
				dbUpsertUserReaction: true,
				redisIncr:            true,
			},
			mocked: mocked{
				dbGetUserByIDResult:       testUser,
				redisGetLimitResult:       "1",
				dbGetUserByIDTargetResult: testUserPremium,
			},
		},
		{
			name: "normal case - successfully add reaction for premium user, with premium data in cache",
			args: args{
				params: reactionParams,
			},
			shouldMock: shouldMock{
				dbGetUserByIDTarget:  true,
				dbUpsertUserReaction: true,
			},
			mocked: mocked{
				cacheGetPremiumResult:     []byte("true"),
				dbGetUserByIDTargetResult: testUserPremium,
			},
		},
		{
			name: "normal case - successfully add reaction for premium user, with premium data NOT in cache",
			args: args{
				params: reactionParams,
			},
			shouldMock: shouldMock{
				dbGetUserByID:        true,
				dbGetUserByIDTarget:  true,
				dbUpsertUserReaction: true,
			},
			mocked: mocked{
				dbGetUserByIDResult:       testUserPremium,
				dbGetUserByIDTargetResult: testUser,
			},
		},
		{
			name: "error case - failed to get user data when premium info not in cache",
			args: args{
				params: reactionParams,
			},
			shouldMock: shouldMock{
				dbGetUserByID: true,
			},
			mocked: mocked{
				dbGetUserByIDError: errors.New("Error GetUserByID"),
			},
			expectedErr: errors.New("Error GetUserByID"),
		},
		{
			name: "error case case - non premium user exceeds limit",
			args: args{
				params: reactionParams,
			},
			shouldMock: shouldMock{
				redisGetLimit: true,
			},
			mocked: mocked{
				cacheGetPremiumResult: []byte("false"),
				redisGetLimitResult:   "10",
			},
			expectedErr: errors.New("Error on\ncode: LIMIT_EXCEEDED; error: Reaction limit exceeded, try again tommorow; field:"),
		},
		{
			name: "error case - error when retrieving target user",
			args: args{
				params: reactionParams,
			},
			shouldMock: shouldMock{
				dbGetUserByIDTarget: true,
			},
			mocked: mocked{
				cacheGetPremiumResult:    []byte("true"),
				redisGetLimitResult:      "1",
				dbGetUserByIDTargetError: errors.New("Error GetUserByID for target"),
			},
			expectedErr: errors.New("Error GetUserByID for target"),
		},
		{
			name: "error case - empty target user",
			args: args{
				params: reactionParams,
			},
			shouldMock: shouldMock{
				dbGetUserByIDTarget: true,
			},
			mocked: mocked{
				cacheGetPremiumResult: []byte("true"),
				redisGetLimitResult:   "1",
			},
			expectedErr: errors.New("Error on\ncode: NOT FOUND; error: User not found:2; field:"),
		},
		{
			name: "error case - failed to save reaction data",
			args: args{
				params: reactionParams,
			},
			shouldMock: shouldMock{
				dbGetUserByIDTarget:  true,
				dbUpsertUserReaction: true,
			},
			mocked: mocked{
				cacheGetPremiumResult:     []byte("true"),
				dbGetUserByIDTargetResult: testUser,
				dbUpsertUserReactionError: errors.New("Error UpsertUserReaction"),
			},
			expectedErr: errors.New("Error UpsertUserReaction"),
		},
	}
	for _, tc := range tests {
		db := mocksrepo.NewPostgresRepository(t)
		redis := mocksrepo.NewRedisRepository(t)
		cache := mocksrepo.NewCacheRepository(t)
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			cache.On("Get", ctx, fmt.Sprintf("premium:%d", tc.args.params.UserID)).Return(tc.mocked.cacheGetPremiumResult, tc.mocked.cacheGetPremiumError)

			if tc.shouldMock.dbGetUserByID {
				db.On("GetUserByID", tc.args.params.UserID).Return(tc.mocked.dbGetUserByIDResult, tc.mocked.dbGetUserByIDError)
			}

			if tc.shouldMock.redisGetLimit {
				redis.On("Get", ctx, uc.BuildReactionLimitCacheKey(tc.args.params.UserID)).Return(tc.mocked.redisGetLimitResult, nil)
			}

			if tc.shouldMock.dbGetUserByIDTarget {
				db.On("GetUserByID", tc.args.params.TargetID).Return(tc.mocked.dbGetUserByIDTargetResult, tc.mocked.dbGetUserByIDTargetError)
			}

			if tc.shouldMock.dbUpsertUserReaction {
				db.On("UpsertUserReaction", tc.args.params).Return(tc.mocked.dbUpsertUserReactionError)
			}

			if tc.shouldMock.redisIncr {
				redis.On("Incr", ctx, uc.BuildReactionLimitCacheKey(tc.args.params.UserID), time.Hour*24).Return(int64(1), nil)
			}

			usecase := uc.NewUserUsecase(defaultAuthConfig, redis, db, cache, &log.Logger{})

			err := usecase.React(ctx, tc.args.params)
			if tc.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
