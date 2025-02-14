package config_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	redis "timble/internal/connection/redis"
	"timble/internal/utils"
	mocksCache "timble/mocks/internal_/connection/cache"
	mocksPostgre "timble/mocks/internal_/connection/postgres"
	"timble/module/users/config"
	"timble/module/users/internal/handler"
)

func TestNewSearchTuningRESTHandler(t *testing.T) {
	s := miniredis.RunT(t)
	tests := []struct {
		name string
	}{
		{
			name: "normal case",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			redisClient, _ := redis.NewClient(s.Host(), s.Port(), "200ms", "0")
			cacheClient := mocksCache.NewCacheInterface(t)
			postgresClient := mocksPostgre.NewPostgresInterface(t)

			result := config.NewUsersHandler(&utils.AuthConfig{}, &zap.Logger{}, cacheClient, redisClient, postgresClient)

			assert.NotNil(t, result)
			assert.IsType(t, &handler.UsersResource{}, result)
		})
	}
	defer s.Close()
}
