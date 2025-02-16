package usecase

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	log "go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"timble/internal/utils"
	"timble/module/users/entity"
)

type UserUsecase interface {
	Create(ctx context.Context, params entity.UserRegistrationParams) (entity.UserToken, error)
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

func (usecase UserUc) Create(ctx context.Context, params entity.UserRegistrationParams) (entity.UserToken, error) {
	userToken := entity.UserToken{}
	bytes, _ := bcrypt.GenerateFromPassword([]byte(params.Password), 14)
	userData := entity.User{
		Username:       params.Username,
		Email:          params.Email,
		HashedPassword: string(bytes),
	}

	err := usecase.db.InsertUser(userData)
	if err != nil {
		return userToken, err
	}

	savedData, err := usecase.db.GetUserByUsername(params.Username)
	if err != nil {
		return userToken, errors.WithStack(err)
	}

	token, _ := usecase.auth.GenerateToken(savedData.ID)
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
		Email:     userData.Email,
		Premium:   userData.Premium,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}

	return userPublicData, nil
}

func (usecase UserUc) React(ctx context.Context, params entity.ReactionParams) error {
	// check premium status
	isPremiumBytes, err := usecase.cache.Get(ctx, BuildPremiumCacheKey(params.UserID))
	isPremiumStr := string(isPremiumBytes)
	isPremium := isPremiumStr == PREMIUM_TRUE_STRING
	if err != nil || isPremiumStr == "" {
		// check in db
		userData, err := usecase.db.GetUserByID(params.UserID)
		if err != nil {
			return errors.WithStack(err)
		}
		isPremium = userData.Premium
	}

	if !isPremium {
		limitStr, _ := usecase.redis.Get(ctx, BuildReactionLimitRedisKey(params.UserID))
		limit, _ := strconv.Atoi(limitStr)
		if limit >= REACTION_LIMIT {
			return utils.NewStandardError("Reaction limit exceeded, try again tommorow", "LIMIT_EXCEEDED", "")
		}
	}

	targetUserData, err := usecase.db.GetUserByID(params.TargetID)
	if err != nil {
		return errors.WithStack(err)
	}

	if targetUserData == nil || targetUserData.ID == 0 {
		return utils.UserNotFoundError(params.TargetID)
	}

	err = usecase.db.UpsertUserReaction(params)
	if err != nil {
		return errors.WithStack(err)
	}

	if !isPremium && params.Type != 0 {
		usecase.redis.Incr(ctx, BuildReactionLimitRedisKey(params.UserID), reactionLimitExpCache)
	}

	return nil
}
