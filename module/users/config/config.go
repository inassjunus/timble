package config

import (
	"net/http"

	cacheLib "github.com/go-redis/cache/v9"
	"go.uber.org/zap"

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

func NewUsersHandler(auth *utils.AuthConfig, logger *zap.Logger, cache *cacheLib.Cache, redisClient *redis.RedisClient, postgresClient *postgres.PostgresClient) *handler.UsersResource {
	redisRepository := repository.NewRedisRepository(redisClient)
	cacheRepository := repository.NewCacheRepository(cache)
	postgresRepository := repository.NewPostgresRepository(postgresClient)

	authUsecase := usecase.NewAuthUsecase(auth, postgresRepository, logger)
	premiumUsecase := usecase.NewPremiumUsecase(postgresRepository, cacheRepository, logger)
	userUsecase := usecase.NewUserUsecase(auth, redisRepository, postgresRepository, cacheRepository, logger)

	return handler.NewUsersResource(authUsecase, premiumUsecase, userUsecase, logger)
}
