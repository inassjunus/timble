package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	cacheLib "github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"

	"timble/internal/utils"
)

// Interface represents redis connection interface
type CacheInterface interface {
	Set(ctx context.Context, key string, value interface{}, expire time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
}

var (
	ignoredErrors = map[string]bool{
		"cache miss": true, // this error is ignorable since it is expected that cache is sometimes empty
	}
)

// Cache represents redis connection object
type CacheClient struct {
	Name   string
	Client *cacheLib.Cache
}

// NewClient creates new redis connection
func NewClient(host, port, timeoutString, dbString string) (*CacheClient, error) {
	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		timeout = 100 * time.Millisecond
	}

	db, err := strconv.Atoi(dbString)
	if err != nil {
		db = 0
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
		DB:           db,
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	client := cacheLib.New(&cacheLib.Options{
		Redis:      redisClient,
		LocalCache: cacheLib.NewTinyLFU(1000, time.Minute),
	})

	cacheClient := &CacheClient{Name: "cache", Client: client}

	return cacheClient, nil
}

// Set add new key-value pair to redis
func (c *CacheClient) Set(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	metricInfo := utils.NewClientMetric("cache", "set")
	cacheItem := &cacheLib.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   expire,
	}
	err := c.Client.Set(cacheItem)
	metricInfo.TrackClientWithError(c.wrapError(err))
	return err
}

// Get retrieve value from redis
func (c *CacheClient) Get(ctx context.Context, key string) (res []byte, err error) {
	metricInfo := utils.NewClientMetric("cache", "get-first")
	err = c.Client.Get(context.Background(), key, &res)
	metricInfo.TrackClientWithError(c.wrapError(err))
	return res, err
}

func (c *CacheClient) wrapError(err error) error {
	if err != nil && !ignoredErrors[err.Error()] {
		return err
	}

	return nil
}
