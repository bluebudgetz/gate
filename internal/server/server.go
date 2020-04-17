package server

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/golangly/webutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/go-playground/validator.v9"

	"github.com/bluebudgetz/gate/internal/config"
	"github.com/bluebudgetz/gate/internal/server/handlers"
)

func NewPrometheusRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}

func NewChiRouter(cfg config.Config, metricsRegistry *prometheus.Registry, routes handlers.Routes) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		middleware.SetHeader("Server", "bluebudgetz/gate"),
		middleware.Heartbeat("/health"),
		middleware.RealIP,
	)
	if !cfg.HTTP.DisableAccessLog {
		router.Use(webutil.RequestLogger)
	}
	router.Use(
		webutil.Metrics(metricsRegistry),
		webutil.RequestID,
		middleware.NoCache,
		cors.New(cors.Options{
			AllowedOrigins:   cfg.HTTP.CORS.AllowOrigins,
			AllowedMethods:   cfg.HTTP.CORS.AllowMethods,
			AllowedHeaders:   cfg.HTTP.CORS.AllowHeaders,
			ExposedHeaders:   cfg.HTTP.CORS.ExposeHeaders,
			AllowCredentials: cfg.HTTP.CORS.AllowCredentials,
			MaxAge:           int(cfg.HTTP.CORS.MaxAge.Seconds()), // 300 is the maximum value not ignored by any of major browsers
		}).Handler,
		middleware.GetHead,
		middleware.RedirectSlashes,
		middleware.Compress(cfg.HTTP.GZipLevel),
		middleware.Timeout(30*time.Second),
	)
	router.Route("/", routes)
	return router
}

func NewHTTPServers(cfg config.Config, router *chi.Mux) []*http.Server {
	monitoringMux := http.NewServeMux()
	monitoringMux.Handle("/metrics", promhttp.Handler())
	return []*http.Server{
		{
			Addr:              ":" + strconv.Itoa(cfg.Monitoring.Port),
			Handler:           monitoringMux,
			ReadTimeout:       cfg.Monitoring.ReadTimeout,
			ReadHeaderTimeout: cfg.Monitoring.ReadHeaderTimeout,
			WriteTimeout:      cfg.Monitoring.WriteTimeout,
			IdleTimeout:       cfg.Monitoring.IdleTimeout,
			MaxHeaderBytes:    cfg.Monitoring.MaxHeaderBytes,
		},
		{
			Addr:              ":" + strconv.Itoa(cfg.HTTP.Port),
			Handler:           router,
			ReadTimeout:       cfg.HTTP.ReadTimeout,
			ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
			WriteTimeout:      cfg.HTTP.WriteTimeout,
			IdleTimeout:       cfg.HTTP.IdleTimeout,
			MaxHeaderBytes:    cfg.HTTP.MaxHeaderBytes,
			BaseContext: func(l net.Listener) context.Context {
				var ctx context.Context
				ctx = context.Background()
				ctx = context.WithValue(ctx, "validator", validator.New())
				return ctx
			},
		},
	}
}
