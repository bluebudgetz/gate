package infra

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func chiMetrics(next http.Handler) http.Handler {
	reqs := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "requests",
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": "gate"},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(reqs)

	latency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "latency",
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": "gate"},
		Buckets:     []float64{300, 1200, 5000},
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(latency)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		defer func() {
			statusCodeStr := strconv.Itoa(ww.Status())
			duration := float64(time.Since(start).Nanoseconds()) / 1000000
			reqs.WithLabelValues(statusCodeStr, r.Method, r.URL.Path).Inc()
			latency.WithLabelValues(statusCodeStr, r.Method, r.URL.Path).Observe(duration)
		}()
		next.ServeHTTP(ww, r)
	})
}
