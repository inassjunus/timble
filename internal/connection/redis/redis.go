package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"timble/internal/utils"
)

// Interface represents redis connection interface
type RedisInterface interface {
	ZAdd(ctx context.Context, key string, member interface{}, score float64) (int64, error)
	ZCard(ctx context.Context, key string) (int64, error)
	ZRemRangeByRank(ctx context.Context, key string, start int64, stop int64) (int64, error)
	ZRevRange(ctx context.Context, key string, start int64, stop int64) ([]string, error)
	ZRange(ctx context.Context, key string, start int64, stop int64) ([]string, error)
	Del(ctx context.Context, key string) (int64, error)
	SAdd(ctx context.Context, key string, member interface{}) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	Set(ctx context.Context, key string, value interface{}, expire time.Duration) (string, error)
	Get(ctx context.Context, key string) (string, error)
	HSet(ctx context.Context, key string, field string, value interface{}) (int64, error)
	HGet(ctx context.Context, key string, field string) (string, error)
}

var (
	ignoredErrors = map[string]bool{
		"redis: nil": true, // this error is expected for some users, so we can ignore it
	}
)

// Redis represents redis connection object
type RedisClient struct {
	Name   string
	Client *redis.Client
}

// NewClient creates new redis connection
func NewClient(host, port, timeoutString, dbString string) (*RedisClient, error) {
	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		timeout = 100 * time.Millisecond
	}

	db, err := strconv.Atoi(dbString)
	if err != nil {
		db = 0
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
		DB:           db,
	})

	redisClient := &RedisClient{Name: "redis", Client: client}

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return redisClient, err
	}

	return redisClient, nil
}

// ZRevRange returns the specified range of elements in the sorted set
func (r *RedisClient) ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	metricInfo := utils.NewClientMetric(r.Name, "zrevrange")
	result, err := r.Client.ZRevRange(ctx, key, start, stop).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return result, err
}

// ZRange returns the specified range of elements in the sorted set in ascending order
func (r *RedisClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	metricInfo := utils.NewClientMetric(r.Name, "zrange")
	result, err := r.Client.ZRange(ctx, key, start, stop).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return result, err
}

// ZAdd add new values to sorted set
func (r *RedisClient) ZAdd(ctx context.Context, key string, member interface{}, score float64) (int64, error) {
	metricInfo := utils.NewClientMetric(r.Name, "zadd")
	value := redis.Z{
		Score:  score,
		Member: member,
	}
	// res == 1 means a new value is successfully added
	// res == 0 means the value is already exists, so no element was modified
	res, err := r.Client.ZAdd(ctx, key, value).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return res, err
}

// ZCard returns the number of elements in the sorted set
func (r *RedisClient) ZCard(ctx context.Context, key string) (int64, error) {
	metricInfo := utils.NewClientMetric(r.Name, "zcard")
	result, err := r.Client.ZCard(ctx, key).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return result, err
}

// ZRemRangeByRank remove element on sorted set on the given rank
func (r *RedisClient) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error) {
	metricInfo := utils.NewClientMetric(r.Name, "zremrangebyrank")
	// result is the number of deleted values
	result, err := r.Client.ZRemRangeByRank(ctx, key, start, stop).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return result, err
}

// Del remove the given key from redis
func (r *RedisClient) Del(ctx context.Context, key string) (int64, error) {
	metricInfo := utils.NewClientMetric(r.Name, "del")
	// result is the number of deleted keys
	result, err := r.Client.Del(ctx, key).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return result, err
}

// SAdd add new values to a set
func (r *RedisClient) SAdd(ctx context.Context, key string, member interface{}) (int64, error) {
	metricInfo := utils.NewClientMetric(r.Name, "sadd")
	// res == 1 means a new value is successfully added
	// res == 0 means the value is already exists, so no element was modified
	res, err := r.Client.SAdd(ctx, key, member).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return res, err
}

// SMembers returns all elements in the set
func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	metricInfo := utils.NewClientMetric(r.Name, "smembers")
	result, err := r.Client.SMembers(ctx, key).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return result, err
}

// Set add new key-value pair to redis
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expire time.Duration) (string, error) {
	metricInfo := utils.NewClientMetric(r.Name, "set")
	// res == 1 means a new value is successfully added
	// res == 0 means the value is already exists, so no element was modified
	res, err := r.Client.Set(ctx, key, value, expire).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return res, err
}

// Get retrieve value from redis
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	metricInfo := utils.NewClientMetric(r.Name, "get")
	res, err := r.Client.Get(ctx, key).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return res, err
}

// HSet add new key-value pair to redis hash
func (r *RedisClient) HSet(ctx context.Context, key string, field string, value interface{}) (int64, error) {
	metricInfo := utils.NewClientMetric(r.Name, "hset")
	// res == 1 means a new value is successfully added
	// res == 0 means the value is already exists, so no element was modified
	res, err := r.Client.HSet(ctx, key, field, value).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return res, err
}

// HGet retrieve value from redis hash map
func (r *RedisClient) HGet(ctx context.Context, key string, field string) (string, error) {
	metricInfo := utils.NewClientMetric(r.Name, "hget")
	res, err := r.Client.HGet(ctx, key, field).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return res, err
}

func (r *RedisClient) wrapError(err error) error {
	if err != nil && !ignoredErrors[err.Error()] {
		return err
	}

	return nil
}
