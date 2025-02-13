package config

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/riandyrn/otelchi"
	"moul.io/chizap"

	"timble/internal/utils"
	usersConfig "timble/module/users/config"
)

type RESTServer struct {
	Server                 *http.Server
	ServerConfig           restServerConfig
	PrometheusServer       *http.Server
	PrometheusServerConfig prometheusConfig
}

func NewRestServer() (*RESTServer, error) {
	prometheusConfig := LoadPrometheusConfig()
	restServerConfig := LoadRestServerConfig()

	restRouter := getRESTRoutes()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", restServerConfig.ServerHost, restServerConfig.ServerPort),
		Handler: restRouter,
	}

	reg := prometheus.NewRegistry()

	utils.RegisterCustomMetrics(reg)

	prometheusServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", prometheusConfig.PrometheusPort),
	}

	restServer := &RESTServer{
		Server:                 server,
		ServerConfig:           restServerConfig,
		PrometheusServer:       prometheusServer,
		PrometheusServerConfig: prometheusConfig,
	}
	return restServer, nil
}

func getRESTRoutes() *chi.Mux {
	conns := NewServiceConnections()
	logger := conns.LoggerClient
	cache := conns.CacheClient
	redis := conns.RedisClient
	postgres := conns.PostgresClient
	auth := conns.Auth

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(utils.ReqIDCtx)
	router.Use(utils.ReqBodyCtx)
	router.Use(middleware.Recoverer)
	router.Use(chizap.New(logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))
	router.Use(otelchi.Middleware("timble", otelchi.WithChiRoutes(router)))

	usersHandler := usersConfig.NewUsersHandler(
		auth,
		logger,
		cache,
		redis,
		postgres,
	)

	// Health check function
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		body := utils.NewResponseMessage("ok", utils.Meta{
			HTTPStatus: http.StatusOK,
		})
		body.WriteAPIResponse(w, r, http.StatusOK)
	})

	router.Route("/api/public/users", func(r chi.Router) {
		r.Post("/register", usersHandler.Create)
	})

	router.Route("/api/protected/users", func(r chi.Router) {
		r.Use(utils.Authentication(auth))
		r.Get("/", usersHandler.Show)
		r.Patch("/react", usersHandler.React)
		r.Route("/premium", func(r chi.Router) {
			r.Patch("/grant", usersHandler.GrantPremium)
			r.Patch("/unsubscribe", usersHandler.UnsubscribePremium)
		})
	})

	router.Route("/api/public/auth", func(r chi.Router) {
		r.Post("/login", usersHandler.Login)
	})

	return router
}
