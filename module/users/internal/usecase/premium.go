package usecase

import (
	"context"

	"github.com/pkg/errors"
	log "go.uber.org/zap"

	"timble/module/users/entity"
)

type PremiumUsecase interface {
	Grant(ctx context.Context, userID uint) error
	Unsubscribe(ctx context.Context, userID uint) error
}

type PremiumUc struct {
	db     PostgresRepository
	cache  CacheRepository
	logger *log.Logger
}

func NewPremiumUsecase(db PostgresRepository, cache CacheRepository, logger *log.Logger) *PremiumUc {
	return &PremiumUc{
		db:     db,
		cache:  cache,
		logger: logger,
	}
}

func (usecase PremiumUc) Grant(ctx context.Context, userID uint) error {
	userData := entity.User{
		ID:      userID,
		Premium: true,
	}
	err := usecase.db.UpdateUserPremium(userData, interface{}(userData.Premium))
	if err != nil {
		return errors.WithStack(err)
	}
	usecase.cache.Set(ctx, BuildPremiumCacheKey(userID), []byte(PREMIUM_TRUE_STRING), premiumExpCache)
	return nil
}

func (usecase PremiumUc) Unsubscribe(ctx context.Context, userID uint) error {
	userData := entity.User{
		ID:      userID,
		Premium: false,
	}
	err := usecase.db.UpdateUserPremium(userData, interface{}(userData.Premium))
	if err != nil {
		return errors.WithStack(err)
	}
	usecase.cache.Set(ctx, BuildPremiumCacheKey(userID), []byte(PREMIUM_FALSE_STRING), premiumExpCache)
	return nil
}
