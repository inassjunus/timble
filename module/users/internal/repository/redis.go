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

func (repo *RedisRepository) ZAddWithLimit(ctx context.Context, key string, member interface{}, score float64, limit int64) error {
	_, err := repo.redisClient.ZAdd(ctx, key, member, score)
	if err != nil {
		return errors.Wrap(err, "redis client error when zadd")
	}

	count, err := repo.redisClient.ZCard(ctx, key)
	if err != nil {
		return errors.Wrap(err, "redis client error when zcard")
	}

	// if number of elements is more than limit, remove the element with the lowest score
	if count > limit {
		lowestIndex := -1 * (limit + 1)
		_, err := repo.redisClient.ZRemRangeByRank(ctx, key, 0, lowestIndex)
		if err != nil {
			return errors.Wrap(err, "redis client error when zremrangebyrank")
		}
	}

	return nil
}

func (repo *RedisRepository) ZRevRange(ctx context.Context, key string, start int64, stop int64) ([]string, error) {
	res, err := repo.redisClient.ZRevRange(ctx, key, start, stop)
	if err != nil {
		return []string{}, errors.Wrap(err, "redis client error when zrevrange")
	}

	return res, nil
}

func (repo *RedisRepository) Del(ctx context.Context, key string) (int64, error) {
	res, err := repo.redisClient.Del(ctx, key)
	if err != nil {
		return res, errors.Wrap(err, "redis client error when del")
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

func (repo *RedisRepository) ZAdd(ctx context.Context, key string, member interface{}, score float64) error {
	_, err := repo.redisClient.ZAdd(ctx, key, member, score)
	if err != nil {
		return errors.Wrap(err, "redis client error when zadd")
	}

	return nil
}

func (repo *RedisRepository) ZCard(ctx context.Context, key string) (int64, error) {
	count, err := repo.redisClient.ZCard(ctx, key)
	if err != nil {
		return int64(0), errors.Wrap(err, "redis client error when zcard")
	}

	return count, nil
}

func (repo *RedisRepository) ZRange(ctx context.Context, key string, start int64, stop int64) ([]string, error) {
	res, err := repo.redisClient.ZRange(ctx, key, start, stop)
	if err != nil {
		return []string{}, errors.Wrap(err, "redis client error when zrange")
	}

	return res, nil
}

func (repo *RedisRepository) SAdd(ctx context.Context, key string, member interface{}) error {
	_, err := repo.redisClient.SAdd(ctx, key, member)
	if err != nil {
		return errors.Wrap(err, "redis client error when sadd")
	}

	return nil
}

func (repo *RedisRepository) SMembers(ctx context.Context, key string) ([]string, error) {
	res, err := repo.redisClient.SMembers(ctx, key)
	if err != nil {
		return []string{}, errors.Wrap(err, "redis client error when smembers")
	}

	return res, nil
}

func (repo *RedisRepository) HSet(ctx context.Context, key string, field string, value interface{}) (int64, error) {
	res, err := repo.redisClient.HSet(ctx, key, field, value)
	if err != nil {
		return res, errors.Wrap(err, "redis client error when hset")
	}

	return res, nil
}

func (repo *RedisRepository) HGet(ctx context.Context, key string, field string) (string, error) {
	res, err := repo.redisClient.HGet(ctx, key, field)
	if err != nil {
		return "", errors.Wrap(err, "redis client error when hget")
	}

	return res, nil
}
