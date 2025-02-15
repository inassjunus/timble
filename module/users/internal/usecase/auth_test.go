package usecase_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	log "go.uber.org/zap"

	"timble/internal/utils"
	mocksrepo "timble/mocks/module/users/internal_/usecase"
	"timble/module/users/entity"
	"timble/module/users/internal/repository"
	uc "timble/module/users/internal/usecase"
)

func TestNewAuthUsecase(t *testing.T) {
	t.Run("new auth usecase", func(t *testing.T) {

		usecase := uc.NewAuthUsecase(
			&utils.AuthConfig{},
			&repository.PostgresRepository{},
			&log.Logger{},
		)

		assert.IsType(t, &uc.AuthUc{}, usecase)
	})
}

func TestAuthUc_Login(t *testing.T) {
	defaultCfg := &utils.AuthConfig{
		SecretKey: []byte("secretz"),
		TokenExp:  time.Hour,
	}

	type args struct {
		params entity.UserLoginParams
		config *utils.AuthConfig
	}

	type mocked struct {
		dbResult *entity.User
		dbError  error
	}
	tests := []struct {
		name           string
		args           args
		mocked         mocked
		expectedResult string
		expectedErr    error
	}{
		{
			name: "normal case - successfully login",
			args: args{
				params: entity.UserLoginParams{
					Username: "testuser",
					Password: "testpassword",
				},
				config: defaultCfg,
			},
			mocked: mocked{
				dbResult: &entity.User{
					ID:             uint(1),
					Email:          "test@email.com",
					Username:       "testuser",
					Premium:        true,
					HashedPassword: "$2a$14$yWjcGVzgVVBZHQV377NA2.R9.Uf7NPoBoHMsBaPboh552vuxhQV06",
				},
			},
			expectedResult: `[a-zA-Z0-9]+\.[a-zA-Z0-9]+\.[a-zA-Z0-9\-\_]+`,
		},
		{
			name: "error case - wrong username",
			args: args{
				params: entity.UserLoginParams{
					Username: "testuser",
					Password: "testpass",
				},
				config: defaultCfg,
			},
			mocked: mocked{
				dbResult: nil,
			},
			expectedErr: errors.New("Invalid username or password"),
		},
		{
			name: "error case - wrong password",
			args: args{
				params: entity.UserLoginParams{
					Username: "testuser",
					Password: "testpass",
				},
				config: defaultCfg,
			},
			mocked: mocked{
				dbResult: &entity.User{
					ID:             uint(1),
					Email:          "test@email.com",
					Username:       "testuser",
					Premium:        true,
					HashedPassword: "$2a$14$yWjcGVzgVVBZHQV377NA2.R9.Uf7NPoBoHMsBaPboh552vuxhQV06",
				},
			},
			expectedErr: errors.New("Invalid username or password"),
		},
		{
			name: "error case - error from db",
			args: args{
				params: entity.UserLoginParams{
					Username: "testuser",
					Password: "testpass",
				},
				config: defaultCfg,
			},
			mocked: mocked{
				dbError: errors.New("DB failed"),
			},
			expectedErr: errors.New("DB failed"),
		},
	}
	for _, tc := range tests {
		db := mocksrepo.NewPostgresRepository(t)
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			db.On("GetUserByUsername", tc.args.params.Username).Return(tc.mocked.dbResult, tc.mocked.dbError)

			usecase := uc.NewAuthUsecase(tc.args.config, db, &log.Logger{})

			result, err := usecase.Login(ctx, tc.args.params)
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
