package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"timble/internal/connection/redis"
)

type RedisRepository struct {
	redisClient redis.RedisInterface
}

func NewRedisRepository(redisClient redis.RedisInterface) *RedisRepository {
	return &RedisRepository{redisClient}
}

func (repo *RedisRepository) Incr(ctx context.Context, key string, expire time.Duration) (int64, error) {
	res, err := repo.redisClient.Incr(ctx, key)
	if err != nil {
		return res, errors.Wrap(err, "redis client error when incr")
	}

	if res == 1 {
		loc, err := time.LoadLocation("Asia/Jakarta")
		expireTime := time.Now()
		if err == nil {
			expireTime = expireTime.In(loc)
		}
		expireTime = expireTime.Add(expire)
		repo.redisClient.ExpireAt(ctx, key, expireTime)
	}

	return res, nil
}

func (repo *RedisRepository) Set(ctx context.Context, key string, value interface{}, expire time.Duration) (string, error) {
	res, err := repo.redisClient.Set(ctx, key, value, expire)
	if err != nil {
		return res, errors.Wrap(err, "redis client error when set")
	}

	return res, nil
}

func (repo *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	res, err := repo.redisClient.Get(ctx, key)
	if err != nil {
		return "", errors.Wrap(err, "redis client error when get")
	}

	return res, nil
}
