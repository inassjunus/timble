package usecase

import (
	"context"

	"github.com/pkg/errors"
	log "go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"timble/internal/utils"
	"timble/module/users/entity"
)

type UserUsecase interface {
	Create(ctx context.Context, params entity.UserParams) (entity.UserToken, error)
	Show(ctx context.Context, userID uint) (*entity.UserPublic, error)
	React(ctx context.Context, params entity.ReactionParams) error
}

type UserUc struct {
	auth   *utils.AuthConfig
	cache  CacheRepository
	redis  RedisRepository
	db     PostgresRepository
	logger *log.Logger
}

func NewUserUsecase(auth *utils.AuthConfig, redis RedisRepository, db PostgresRepository, cache CacheRepository, logger *log.Logger) *UserUc {
	return &UserUc{
		auth:   auth,
		redis:  redis,
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

func (usecase UserUc) Create(ctx context.Context, params entity.UserParams) (entity.UserToken, error) {
	userToken := entity.UserToken{}
	bytes, err := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	userData := entity.User{
		Username:       params.Username,
		HashedPassword: string(bytes),
	}

	err = usecase.db.InsertUser(userData)
	if err != nil {
		return userToken, errors.WithStack(err)
	}

	savedData, err := usecase.db.GetUserByUsername(params.Username)
	if err != nil {
		return userToken, errors.WithStack(err)
	}

	token, err := usecase.auth.GenerateToken(savedData.ID)
	if err != nil {
		return userToken, errors.Wrap(err, "Error generating token")
	}

	userToken.Token = token

	return userToken, nil

}

func (usecase UserUc) Show(ctx context.Context, userID uint) (*entity.UserPublic, error) {
	userData, err := usecase.db.GetUserByID(userID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userPublicData := &entity.UserPublic{
		ID:        userData.ID,
		Username:  userData.Username,
		Premium:   userData.Premium,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}

	return userPublicData, nil
}

func (usecase UserUc) React(ctx context.Context, params entity.ReactionParams) error {
	// check premium here
	// if not premium, check for limit in redis

	userData, err := usecase.db.GetUserByID(params.UserID)
	if err != nil {
		return errors.WithStack(err)
	}

	if userData == nil {
		return errors.New("User not found")
	}

	targetUserData, err := usecase.db.GetUserByID(params.TargetID)
	if err != nil {
		return errors.WithStack(err)
	}

	if targetUserData == nil {
		return errors.New("Target user not found")
	}

	// validate reaction type
	if valid := entity.ReactionTypes[params.Type]; !valid {
		return errors.New("Invalid reaction")
	}

	err = usecase.db.UpsertUserReaction(params)
	if err != nil {
		return errors.WithStack(err)
	}

	// update cache here

	return nil
}
