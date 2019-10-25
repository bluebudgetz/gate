package infra

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/bluebudgetz/gate/internal/config"
)

func NewChiRouter(cfg config.Config) chi.Router {
	r := chi.NewRouter()
	//cfg.HTTP.DisableLogRequests
	r.Use(
		middleware.SetHeader("Server", "bluebudgetz/gate"),
		middleware.Heartbeat("/health"),
		middleware.RealIP,
	)
	if !cfg.HTTP.DisableLogRequests {
		r.Use(RequestLogger)
	}
	r.Use(
		chiMetrics,
		chiRequestID,
		middleware.NoCache,
		corsHandler(cfg),
		middleware.GetHead,
		middleware.RedirectSlashes,
		middleware.Compress(cfg.HTTP.GZipLevel),
		middleware.Timeout(30*time.Second),
	)
	return r
}

func corsHandler(cfg config.Config) func(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   cfg.HTTP.CORS.AllowOrigins,
		AllowedMethods:   cfg.HTTP.CORS.AllowMethods,
		AllowedHeaders:   cfg.HTTP.CORS.AllowHeaders,
		ExposedHeaders:   cfg.HTTP.CORS.ExposeHeaders,
		AllowCredentials: cfg.HTTP.CORS.AllowCredentials,
		MaxAge:           cfg.HTTP.CORS.MaxAge, // 300 is the maximum value not ignored by any of major browsers
	}).Handler
}
