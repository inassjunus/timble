package config

import (
	"net/http"

	"go.uber.org/zap"

	cache "timble/internal/connection/cache"
	postgres "timble/internal/connection/postgres"
	redis "timble/internal/connection/redis"
	"timble/internal/utils"
	"timble/module/users/internal/handler"
	"timble/module/users/internal/repository"
	"timble/module/users/internal/usecase"
)

// rest handler
type UsersRESTInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	React(w http.ResponseWriter, r *http.Request)
	GrantPremium(w http.ResponseWriter, r *http.Request)
	UnsubscribePremium(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

func NewUsersHandler(auth *utils.AuthConfig, logger *zap.Logger, cache cache.CacheInterface, redisClient redis.RedisInterface, postgresClient postgres.PostgresInterface) *handler.UsersResource {
	redisRepository := repository.NewRedisRepository(redisClient)
	cacheRepository := repository.NewCacheRepository(cache)
	postgresRepository := repository.NewPostgresRepository(postgresClient)

	authUsecase := usecase.NewAuthUsecase(auth, postgresRepository, logger)
	premiumUsecase := usecase.NewPremiumUsecase(redisRepository, postgresRepository, cacheRepository, logger)
	userUsecase := usecase.NewUserUsecase(auth, redisRepository, postgresRepository, cacheRepository, logger)

	return handler.NewUsersResource(authUsecase, premiumUsecase, userUsecase, logger)
}
