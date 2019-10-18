package infra

import (
	"context"
	"github.com/bluebudgetz/gate/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func NewMonitoringHTTPServer(cfg config.MonitoringConfig) (*http.Server, func(ctx context.Context), error) {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{Addr: ":" + strconv.Itoa(cfg.MetricsPort), Handler: mux}
	return server, func(ctx context.Context) {
		if err := server.Shutdown(ctx); err != nil {
			log.Warn().Err(err).Msg("Failed stopping monitoring HTTP server")
		}
	}, nil
}
