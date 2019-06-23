package main

import (
	"fmt"
	chiprometheus "github.com/766b/chi-prometheus"
	"github.com/bluebudgetz/gate/internal/accounts"
	"github.com/bluebudgetz/gate/internal/schema"
	"github.com/bluebudgetz/gate/internal/transactions"
	"github.com/bluebudgetz/gate/internal/util"
	mongoutil "github.com/bluebudgetz/gate/internal/util/mongo"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	stdlog "log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	config := util.Config{}

	// Use Viper for configuration
	v := viper.New()
	v.SetDefault("environment", util.EnvDevelopment)
	v.SetDefault("http.port", 3001)
	v.SetDefault("http.cors.host", "www.bluebudgetz.com")
	v.SetDefault("http.cors.port", 80)
	v.SetDefault("loglevel", "info")
	v.SetDefault("metrics.port", 3002)
	v.SetDefault("database.mongodb.uri", "mongodb://localhost:27017")
	v.SetConfigName("gate")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/gate/")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.SetEnvPrefix("GATE")
	v.AutomaticEnv()
	if err := v.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("Failed configuring application.")
	}

	// Configuration validations
	if config.Environment != util.EnvProduction && config.Environment != util.EnvDevelopment {
		log.Fatal().
			Interface("env", os.Environ()).
			Msgf("environment must be either '%s' or '%s'", util.EnvProduction, util.EnvDevelopment)
	}

	// Logging level
	switch config.LogLevel {
	case "disabled":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		log.Fatal().Msg("log level must be one of: disabled, panic, fatal, error, warn, info or debug")
	}

	// Add common metadata
	log.Logger = log.With().
		Str("service", "gate").
		Str("env", config.Environment).
		Caller().
		Stack(). // TODO: replace standard zerolog Stacktrace hook with one that also supports causing errors
		Logger()
	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)

	// Pretty logging if in development mode
	if config.Environment == util.EnvDevelopment {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	// Print configuration
	log.Info().
		Interface("args", os.Args).
		Interface("env", os.Environ()).
		Interface("config", config).
		Interface("keys", v.AllKeys()).
		Msg("Configured")

	// Prepare MongoDB client
	mongoClient, err := mongoutil.CreateMongoClient(config.Database.MongoDB.URI)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating MongoDB client.")
	}

	// Prepare JSON schema registry
	schemaRegistry, err := schema.NewSchemaRegistry()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating JSON schema registry.")
	}

	// Create router
	router := chi.NewRouter()
	router.Use(
		middleware.SetHeader("Server", "bluebudgetz-gate"),
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.Heartbeat("/health"),
		middleware.RealIP,
		middleware.RequestID,
		middleware.RequestLogger(&zeroLogLogFormatter{}),
		middleware.Recoverer,
		chiprometheus.NewMiddleware("serviceName", 50, 200, 500, 1000, 2000, 5000, 10000, 30000),
		middleware.GetHead,
		middleware.NoCache,
		middleware.ContentCharset("", "UTF-8"),
		middleware.AllowContentType("application/json"),
		middleware.SetHeader("Content-Type", "application/json; charset=UTF-8"),
	)
	if config.HTTP.CORS.Host != "" && config.HTTP.CORS.Port != 0 {
		corsHost := config.HTTP.CORS.Host
		corsPort := config.HTTP.CORS.Port
		corsInstance := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://" + corsHost + ":" + strconv.Itoa(corsPort), "https://" + corsHost + ":" + strconv.Itoa(corsPort)},
			AllowedMethods:   []string{"OPTIONS", "HEAD", "GET", "POST", "PATCH", "PUT", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		})
		if log.Debug().Enabled() {
			corsInstance.Log = stdlog.New(log.Logger.With().Str("module", "cors").Logger(), "", 0)
		}
		router.Use(corsInstance.Handler)
	}
	router.Route("/v1/accounts", accounts.New(config, schemaRegistry, mongoClient).RoutesV1)
	router.Route("/v1/transactions", transactions.New(config, schemaRegistry, mongoClient).RoutesV1)

	// Start an HTTP server that only provides "/metrics" on a different port
	// This port SHOULD NOT be exposed externally
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsHttpServer := &http.Server{Addr: fmt.Sprintf(":%d", config.Metrics.Port), Handler: metricsMux}
	log.Info().Msgf("Starting metrics server on port %d", config.Metrics.Port)
	go func() {
		if err := metricsHttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Metrics HTTP server failed")
		}
	}()

	// Start the application HTTP server, and wait until it exits
	server := &http.Server{Addr: ":" + strconv.Itoa(config.HTTP.Port), Handler: router}
	log.Info().Msgf("Starting API server on port %d", config.HTTP.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("HTTP server exited with an error")
	} else {
		log.Info().Msg("Done")
		os.Exit(0)
	}
}

type zeroLogLogFormatter struct{}

type zeroLogLogEntry struct {
	r *http.Request
}

func (f *zeroLogLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &zeroLogLogEntry{r}
}

func (e *zeroLogLogEntry) Write(status, bytes int, elapsed time.Duration) {
	log.Info().
		Str("remoteAddr", e.r.RemoteAddr).
		Str("proto", e.r.Proto).
		Str("method", e.r.Method).
		Str("requestURI", e.r.RequestURI).
		Str("headers", fmt.Sprintf("%v", e.r.Header)).
		Str("url", e.r.URL.String()).
		Str("host", e.r.Host).
		Int("status", status).
		Int("bytesWritten", bytes).
		Dur("elapsed", elapsed).
		Msg("HTTP request completed")
}

func (e *zeroLogLogEntry) Panic(v interface{}, stack []byte) {
	log.Error().
		Str("remoteAddr", e.r.RemoteAddr).
		Str("proto", e.r.Proto).
		Str("method", e.r.Method).
		Str("requestURI", e.r.RequestURI).
		Str("headers", fmt.Sprintf("%v", e.r.Header)).
		Str("url", e.r.URL.String()).
		Str("host", e.r.Host).
		Interface("recovered", v).
		Bytes("providedStack", stack).
		Msg("HTTP request failed")
}
