package repository_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"timble/internal/connection/postgres"
	mockspostgres "timble/mocks/internal_/connection/postgres"
	"timble/module/users/entity"
	"timble/module/users/internal/repository"
)

var (
	testUser = &entity.User{
		ID:             uint(1),
		Email:          "test@email.com",
		Username:       "testuser",
		Premium:        true,
		HashedPassword: "testhashespassword",
	}
)

func TestNewPostgresRepository(t *testing.T) {
	t.Run("new postgres repository", func(t *testing.T) {
		repo := repository.NewPostgresRepository(&postgres.PostgresClient{})

		assert.IsType(t, &repository.PostgresRepository{}, repo)
	})
}

func TestPostgresRepository_GetUserByID(t *testing.T) {
	blankResult := &entity.User{}
	tests := []struct {
		name             string
		args             uint
		expectedError    error
		expectedResult   *entity.User
		mockPostgresCall func(postgresClient *mockspostgres.PostgresInterface)
	}{
		{
			name:           "normal case - successfully get user",
			args:           testUser.ID,
			expectedResult: testUser,
			mockPostgresCall: func(postgresClient *mockspostgres.PostgresInterface) {
				postgresClient.On("GetFirst", blankResult, fmt.Sprintf("id='%d'", testUser.ID)).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*entity.User)
					arg.ID = testUser.ID
					arg.Email = testUser.Email
					arg.Username = testUser.Username
					arg.Premium = testUser.Premium
					arg.HashedPassword = testUser.HashedPassword
				}).Return(nil)
			},
		},
		{
			name:           "error case - error when querying",
			args:           testUser.ID,
			expectedResult: blankResult,
			mockPostgresCall: func(postgresClient *mockspostgres.PostgresInterface) {
				postgresClient.On("GetFirst", blankResult, fmt.Sprintf("id='%d'", testUser.ID)).Return(errors.New("timeout"))
			},
			expectedError: errors.New("postgres client error when get user by ID: timeout"),
		},
	}

	for _, tc := range tests {
		postgresClient := mockspostgres.NewPostgresInterface(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.mockPostgresCall(postgresClient)
			repo := repository.NewPostgresRepository(postgresClient)
			result, err := repo.GetUserByID(tc.args)

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

func TestPostgresRepository_GetUserByUsername(t *testing.T) {
	blankResult := &entity.User{}
	tests := []struct {
		name             string
		args             string
		expectedError    error
		expectedResult   *entity.User
		mockPostgresCall func(postgresClient *mockspostgres.PostgresInterface)
	}{
		{
			name:           "normal case - successfully get user",
			args:           testUser.Username,
			expectedResult: testUser,
			mockPostgresCall: func(postgresClient *mockspostgres.PostgresInterface) {
				postgresClient.On("GetFirst", blankResult, fmt.Sprintf("username='%s'", testUser.Username)).Run(func(args mock.Arguments) {
					arg := args.Get(0).(*entity.User)
					arg.ID = testUser.ID
					arg.Email = testUser.Email
					arg.Username = testUser.Username
					arg.Premium = testUser.Premium
					arg.HashedPassword = testUser.HashedPassword
				}).Return(nil)
			},
		},
		{
			name:           "error case - error when querying",
			args:           testUser.Username,
			expectedResult: blankResult,
			mockPostgresCall: func(postgresClient *mockspostgres.PostgresInterface) {
				postgresClient.On("GetFirst", blankResult, fmt.Sprintf("username='%s'", testUser.Username)).Return(errors.New("timeout"))
			},
			expectedError: errors.New("postgres client error when get user by username: timeout"),
		},
	}

	for _, tc := range tests {
		postgresClient := mockspostgres.NewPostgresInterface(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.mockPostgresCall(postgresClient)
			repo := repository.NewPostgresRepository(postgresClient)
			result, err := repo.GetUserByUsername(tc.args)

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

func TestPostgresRepository_InsertUser(t *testing.T) {
	postgreParams := []interface{}{
		testUser.Username,
		testUser.Email,
		testUser.Premium,
		testUser.HashedPassword,
	}
	tests := []struct {
		name             string
		args             entity.User
		expectedError    error
		mockPostgresCall func(postgresClient *mockspostgres.PostgresInterface)
	}{
		{
			name: "normal case - successfully insert user",
			args: *testUser,
			mockPostgresCall: func(postgresClient *mockspostgres.PostgresInterface) {
				postgresClient.On("Exec", repository.INSERT_USER_QUERY, postgreParams).Return(nil)
			},
		},
		{
			name: "error case - unexpected error during insert",
			args: *testUser,
			mockPostgresCall: func(postgresClient *mockspostgres.PostgresInterface) {
				postgresClient.On("Exec", repository.INSERT_USER_QUERY, postgreParams).Return(errors.New("timeout"))
			},
			expectedError: errors.New("postgres client error when insert to users: timeout"),
		},
	}

	for _, tc := range tests {
		postgresClient := mockspostgres.NewPostgresInterface(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.mockPostgresCall(postgresClient)
			repo := repository.NewPostgresRepository(postgresClient)
			err := repo.InsertUser(tc.args)

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestPostgresRepository_UpdateUserPremium(t *testing.T) {
	tests := []struct {
		name             string
		args             entity.User
		expectedError    error
		mockPostgresCall func(postgresClient *mockspostgres.PostgresInterface)
	}{
		{
			name: "normal case - successfully update user",
			args: *testUser,
			mockPostgresCall: func(postgresClient *mockspostgres.PostgresInterface) {
				postgresClient.On("Exec", repository.UPDATE_USER_PREMIUM_QUERY, testUser.Premium, testUser.ID).Return(nil)
			},
		},
		{
			name: "error case - unexpected error during update",
			args: *testUser,
			mockPostgresCall: func(postgresClient *mockspostgres.PostgresInterface) {
				postgresClient.On("Exec", repository.UPDATE_USER_PREMIUM_QUERY, testUser.Premium, testUser.ID).Return(errors.New("timeout"))
			},
			expectedError: errors.New("postgres client error when update premium to users: timeout"),
		},
	}

	for _, tc := range tests {
		postgresClient := mockspostgres.NewPostgresInterface(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.mockPostgresCall(postgresClient)
			repo := repository.NewPostgresRepository(postgresClient)
			err := repo.UpdateUserPremium(tc.args, tc.args.Premium)

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestPostgresRepository_UpsertUserReaction(t *testing.T) {
	reaction := entity.ReactionParams{
		UserID:   testUser.ID,
		TargetID: uint(2),
		Type:     1,
	}
	postgreParams := []interface{}{
		reaction.UserID,
		reaction.TargetID,
		reaction.Type,
	}
	tests := []struct {
		name             string
		args             entity.ReactionParams
		expectedError    error
		mockPostgresCall func(postgresClient *mockspostgres.PostgresInterface)
	}{
		{
			name: "normal case - successfully upsert user reaction",
			args: reaction,
			mockPostgresCall: func(postgresClient *mockspostgres.PostgresInterface) {
				postgresClient.On("Exec", repository.UPSERT_USER_REACTION, postgreParams, reaction.Type).Return(nil)
			},
		},
		{
			name: "error case - unexpected error during upsert",
			args: reaction,
			mockPostgresCall: func(postgresClient *mockspostgres.PostgresInterface) {
				postgresClient.On("Exec", repository.UPSERT_USER_REACTION, postgreParams, reaction.Type).Return(errors.New("timeout"))
			},
			expectedError: errors.New("postgres client error when upsert to user_reactions: timeout"),
		},
	}

	for _, tc := range tests {
		postgresClient := mockspostgres.NewPostgresInterface(t)

		t.Run(tc.name, func(t *testing.T) {
			tc.mockPostgresCall(postgresClient)
			repo := repository.NewPostgresRepository(postgresClient)
			err := repo.UpsertUserReaction(tc.args)

			if tc.expectedError != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectedError.Error(), err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
