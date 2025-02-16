package usecase

import (
	"context"
	"fmt"
	"time"

	"timble/module/users/entity"
)

const (
	PREMIUM_TRUE_STRING  = "true"
	PREMIUM_FALSE_STRING = "false"
	REACTION_LIMIT       = 10
)

var (
	premiumExpCache       = 24 * time.Hour
	reactionLimitExpCache = 24 * time.Hour
)

type RedisRepository interface {
	Set(ctx context.Context, key string, value interface{}, expire time.Duration) (string, error)
	Get(ctx context.Context, key string) (string, error)
	Incr(ctx context.Context, key string, expire time.Duration) (int64, error)
}

type CacheRepository interface {
	Get(ctx context.Context, key string) (res []byte, err error)
	Set(ctx context.Context, key string, data []byte, exp time.Duration) error
}

type PostgresRepository interface {
	GetUserByID(id uint) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	InsertUser(user entity.User) error
	UpdateUserPremium(user entity.User, value interface{}) error
	UpsertUserReaction(reaction entity.ReactionParams) error
}

func BuildPremiumCacheKey(userID uint) string {
	return fmt.Sprintf("premium:%d", userID)
}

func BuildReactionLimitCacheKey(userID uint) string {
	currentTime := time.Now()
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err == nil {
		currentTime = currentTime.In(loc)
	}
	return fmt.Sprintf("reaction:%s:%d", currentTime.Format("2006-01-02"), userID)
}
