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
	Set(ctx context.Context, key string, value interface{}, expire time.Duration) (string, error)
	Get(ctx context.Context, key string) (string, error)
	ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
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

// Increment a key-value pair to redis
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	metricInfo := utils.NewClientMetric(r.Name, "incr")
	res, err := r.Client.Incr(ctx, key).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return res, err
}

// Expire a key-value pair to redis
func (r *RedisClient) ExpireAt(ctx context.Context, key string, tm time.Time) (bool, error) {
	metricInfo := utils.NewClientMetric(r.Name, "expire-at")
	res, err := r.Client.ExpireAt(ctx, key, tm).Result()
	err = r.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return res, err
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

func (r *RedisClient) wrapError(err error) error {
	if err != nil && !ignoredErrors[err.Error()] {
		return err
	}

	return nil
}
