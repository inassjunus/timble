package usecase

import (
	"context"

	"github.com/pkg/errors"
	log "go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"timble/internal/utils"
	"timble/module/users/entity"
)

type AuthUsecase interface {
	Login(ctx context.Context, params entity.UserLoginParams) (entity.UserToken, error)
}

type AuthUc struct {
	auth   *utils.AuthConfig
	db     PostgresRepository
	logger *log.Logger
}

func NewAuthUsecase(auth *utils.AuthConfig, db PostgresRepository, logger *log.Logger) *AuthUc {
	return &AuthUc{
		auth:   auth,
		db:     db,
		logger: logger,
	}
}

func (usecase AuthUc) Login(ctx context.Context, params entity.UserLoginParams) (entity.UserToken, error) {
	userToken := entity.UserToken{}
	userData, err := usecase.db.GetUserByUsername(params.Username)
	if err != nil {
		return userToken, errors.WithStack(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userData.HashedPassword), []byte(params.Password))
	if err != nil {
		return userToken, errors.New("Invalid username or password")
	}

	token, err := usecase.auth.GenerateToken(userData.ID)
	if err != nil {
		return userToken, errors.Wrap(err, "Error generating token")
	}

	userToken.Token = token

	return userToken, nil
}
