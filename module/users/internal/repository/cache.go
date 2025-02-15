package repository

import (
	"context"
	"time"

	"timble/internal/connection/cache"

	"github.com/pkg/errors"
)

type CacheRepository struct {
	cacheClient cache.CacheInterface
}

func NewCacheRepository(cacheClient cache.CacheInterface) *CacheRepository {
	return &CacheRepository{
		cacheClient: cacheClient,
	}
}

// Get data stored from cache
func (repo *CacheRepository) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := repo.cacheClient.Get(ctx, key)
	if err != nil {
		return nil, errors.Wrap(err, "cache client error when get")
	}

	return res, nil
}

// Set data in cache
func (repo *CacheRepository) Set(ctx context.Context, key string, data []byte, expire time.Duration) error {
	err := repo.cacheClient.Set(ctx, key, data, expire)
	if err != nil {
		return errors.Wrap(err, "cache client error when set")
	}

	return nil
}
