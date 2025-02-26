package config

import (
	"time"

	"github.com/caarlos0/env/v6"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	cache "timble/internal/connection/cache"
	postgres "timble/internal/connection/postgres"
	redis "timble/internal/connection/redis"
	"timble/internal/utils"
)

type ServiceConnections struct {
	LoggerClient   *zap.Logger
	CacheClient    *cache.CacheClient
	RedisClient    *redis.RedisClient
	PostgresClient *postgres.PostgresClient

	Auth *utils.AuthConfig
}

type authConfig struct {
	SecretKey       string `env:"SECRET"`
	TokenExpiration string `env:"TOKEN_EXPIRATION"`
}

type restServerConfig struct {
	ServerHost string `env:"SERVER_HOST"`
	ServerPort int    `env:"SERVER_PORT"`
}

type prometheusConfig struct {
	PrometheusPort int `env:"PROMETHEUS_PORT"`
}

type redisConfig struct {
	Host    string `env:"REDIS_HOST"`
	Port    string `env:"REDIS_PORT"`
	Timeout string `env:"REDIS_TIMEOUT"`
	DB      string `env:"REDIS_DB"`
}

type cacheConfig struct {
	Host    string `env:"CACHE_HOST"`
	Port    string `env:"CACHE_PORT"`
	Timeout string `env:"CACHE_TIMEOUT"`
	DB      string `env:"CACHE_DB"`
}

type databaseConfig struct {
	Host         string `env:"DB_HOST" envDefault:"127.0.0.1"`
	Port         int    `env:"DB_PORT" envDefault:"5432"`
	Username     string `env:"DB_USERNAME" envDefault:"user"`
	Password     string `env:"DB_PASSWORD" envDefault:"pass"`
	Database     string `env:"DB_NAME" envDefault:"database"`
	MaxLifetime  int    `env:"DB_MAX_LIFETIME" envDefault:"5"`
	MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
	MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"100"`
}

func LoadAuthConfig() authConfig {
	authConfig := authConfig{}
	env.Parse(&authConfig)
	return authConfig
}

func LoadRestServerConfig() restServerConfig {
	restServerCfg := restServerConfig{}
	env.Parse(&restServerCfg)
	return restServerCfg
}

func LoadPrometheusConfig() prometheusConfig {
	prometheusCfg := prometheusConfig{}
	env.Parse(&prometheusCfg)
	return prometheusCfg
}

func LoadRedisConfig() redisConfig {
	redisCfg := redisConfig{}
	env.Parse(&redisCfg)
	return redisCfg
}

func LoadCacheConfig() cacheConfig {
	cacheCfg := cacheConfig{}
	env.Parse(&cacheCfg)
	return cacheCfg
}

func LoadDatabaseConfig() databaseConfig {
	dbConfig := databaseConfig{}
	env.Parse(&dbConfig)
	return dbConfig
}

func NewServiceConnections() *ServiceConnections {
	redisConfig := LoadRedisConfig()
	cacheConfig := LoadCacheConfig()
	databaseConfig := LoadDatabaseConfig()
	authConfig := LoadAuthConfig()

	tokenExp := 1 * time.Hour // Token valid for 1 hour by default
	if t, err := time.ParseDuration(authConfig.TokenExpiration); err == nil {
		tokenExp = t
	}

	auth := &utils.AuthConfig{
		SecretKey: []byte(authConfig.SecretKey),
		TokenExp:  tokenExp,
	}

	logger, _ := zap.NewProduction(zap.AddStacktrace(zapcore.FatalLevel + 1))

	// redis for data storage
	redisClient, err := redis.NewClient(redisConfig.Host, redisConfig.Port, redisConfig.Timeout, redisConfig.DB)
	if err != nil {
		panic(err)
	}

	// redis for cache
	cacheClient, err := cache.NewClient(cacheConfig.Host, cacheConfig.Port, cacheConfig.Timeout, cacheConfig.DB)
	if err != nil {
		panic(err)
	}

	postgresClient := &postgres.PostgresClient{
		Name:             "postgres",
		GormOpenFunc:     postgres.OpenGorm,
		PostgresOpenFunc: postgres.OpenPostgres,
		GormGetDBFunc:    postgres.GetSQLDB,
	}

	wrappedPostgresClient, _ := postgres.NewClient(
		postgresClient,
		databaseConfig.Host,
		databaseConfig.Port,
		databaseConfig.Database,
		databaseConfig.Username,
		databaseConfig.Password,
		databaseConfig.MaxIdleConns,
		databaseConfig.MaxOpenConns,
	)

	return &ServiceConnections{
		LoggerClient:   logger,
		CacheClient:    cacheClient,
		RedisClient:    redisClient,
		PostgresClient: wrappedPostgresClient,
		Auth:           auth,
	}
}
