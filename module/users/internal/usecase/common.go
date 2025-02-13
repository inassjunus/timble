package usecase

import (
	"context"
	"fmt"
	"time"

	"timble/module/users/entity"
)

var (
	premiumExpCache = 24 * time.Hour
)

type RedisRepository interface {
	ZAdd(ctx context.Context, key string, member interface{}, score float64) error
	ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZCard(ctx context.Context, key string) (int64, error)
	SAdd(ctx context.Context, key string, member interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	Set(ctx context.Context, key string, value interface{}, expire time.Duration) (string, error)
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) (int64, error)
	HSet(ctx context.Context, key string, field string, value interface{}) (int64, error)
	HGet(ctx context.Context, key string, field string) (string, error)
}

type CacheRepository interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte, expired time.Duration) error
}

type PostgresRepository interface {
	GetUserByID(id uint) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	InsertUser(user entity.User) error
	UpdateUser(user entity.User) error
	UpsertUserReaction(reaction entity.ReactionParams) error
}

func buildPremiumCacheKey(userID uint) string {
	return fmt.Sprintf("premium:%d", userID)
}
