package repository

import (
	"context"
	"time"

	"timble/internal/utils"

	cacheLib "github.com/go-redis/cache/v9"
)

var (
	ignoredErrors = map[string]bool{
		"cache miss": true, // this error is ignorable since it is expected that cache is empty
	}
)

type CacheRepository struct {
	cache *cacheLib.Cache
}

func NewCacheRepository(cache *cacheLib.Cache) *CacheRepository {
	return &CacheRepository{
		cache: cache,
	}
}

// Get retrieves purchase history stored in cache by customer ID
func (repo *CacheRepository) Get(key string) (res []byte, err error) {
	metricInfo := utils.NewClientMetric("cache", "get-first")
	err = repo.cache.Get(context.Background(), key, &res)
	err = repo.wrapError(err)
	metricInfo.TrackClientWithError(err)
	return res, err
}

// Set purchase history in cache by customer ID
func (repo *CacheRepository) Set(key string, data []byte, exp time.Duration) (err error) {
	metricInfo := utils.NewClientMetric("cache", "set")

	cacheItem := &cacheLib.Item{
		Ctx:   context.Background(),
		Key:   key,
		Value: data,
		TTL:   exp,
	}
	err = repo.cache.Set(cacheItem)
	metricInfo.TrackClientWithError(err)
	return err
}

func (repo *CacheRepository) wrapError(err error) error {
	if err != nil && !ignoredErrors[err.Error()] {
		return err
	}

	return nil
}
