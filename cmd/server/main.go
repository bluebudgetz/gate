package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/go-playground/validator.v9"

	"github.com/bluebudgetz/gate/internal/infra"
	"github.com/bluebudgetz/gate/internal/pubsub"
	"github.com/bluebudgetz/gate/internal/rest"
	"github.com/bluebudgetz/gate/internal/rest/accounts"
	"github.com/bluebudgetz/gate/internal/rest/transactions"
)

const (
	ExitCodeOK              = 0
	ExitCodeBadConfig       = 1
	ExitCodeInternalError   = 2
	ExitCodeExternalService = 3
	ExitCodeListen          = 4
)

type DatabaseConfig struct {
	MongoURI string `long:"mongo-uri" env:"MONGODB_URI" value-name:"URI" default:"mongodb://localhost:27017" description:"Mongo URI"`
	RedisURI string `long:"redis-uri" env:"REDIS_URI" value-name:"URI" default:"localhost:6379" description:"Redis URI"`
}

type MonitoringConfig struct {
	Port              int           `long:"metrics-port" env:"METRICS_PORT" value-name:"PORT" default:"3002" description:"Port to service Prometheus metrics"`
	ReadTimeout       time.Duration `long:"read-timeout" env:"READ_TIMEOUT" default:"3s" description:"Maximum number of seconds to read the entire request, including the body"`
	ReadHeaderTimeout time.Duration `long:"read-header-timeout" env:"READ_HEADER_TIMEOUT" default:"1s" description:"Maximum number of seconds to read the request headers"`
	WriteTimeout      time.Duration `long:"write-timeout" env:"WRITE_TIMEOUT" default:"5s" description:"Maximum number of seconds to write the response"`
	IdleTimeout       time.Duration `long:"idle-timeout" env:"IDLE_TIMEOUT" default:"30s" description:"Maximum number of seconds to let keep-alive connections to live"`
	MaxHeaderBytes    int           `long:"max-header-bytes" env:"MAX_HEADER_BYTES" default:"8192" description:"Maximum number of bytes to read for the request headers"`
}

type Config struct {
	Database   DatabaseConfig   `group:"Databases" namespace:"db" env-namespace:"DB"`
	Monitoring MonitoringConfig `group:"Monitoring" namespace:"monitoring" env-namespace:"MONITORING"`
	HTTP       infra.HTTPConfig `group:"HTTP server" namespace:"http" env-namespace:"HTTP"`
	PubSub     pubsub.Config    `group:"Pub/Sub" namespace:"pubsub" env-namespace:"PUBSUB"`
}

func main() {
	os.Exit(mainWithReturnCode())
}

func mainWithReturnCode() (exitCode int) {

	// We must FIRST configure logging properly (pretty/JSON, stdout/stderr, etc)
	infra.SetupLogging()

	// Parse environment variables and/or command-line arguments, to form a Config object
	cfg := Config{}
	parser := flags.NewParser(&cfg, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"
	parser.LongDescription = "Bluebudgetz API gateway. This is the API micro-service centralizing Bluebudgetz APIs."
	if _, err := parser.Parse(); err != nil {
		log.Error().Err(err).Msg("Failed loading configuration")
		return ExitCodeBadConfig
	}
	log.Info().
		Interface("config", cfg).
		Interface("env", os.Environ()).
		Interface("args", os.Args[1:]).
		Msg("Configuration loaded")

	// Defer a panic handler
	defer func() {
		if r := recover(); r != nil {
			log.Error().
				Err(fmt.Errorf("SYSTEM PANIC: %v", r)).
				Str(zerolog.ErrorStackFieldName, string(debug.Stack())).
				Interface("recovered", r).
				Msg("System error")
			exitCode = ExitCodeInternalError
		}
	}()

	// Context for bootstrapping
	startupCtx, cancelStartupCtx := context.WithTimeout(context.Background(), 30*time.Second)

	// Connect to MongoDB
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(cfg.Database.MongoURI))
	if err != nil {
		log.Error().Err(err).Msg("Failed creating MongoDB client")
		return ExitCodeExternalService
	} else if err := mongoClient.Connect(startupCtx); err != nil {
		log.Error().Err(err).Msg("Failed connecting MongoDB client")
		return ExitCodeExternalService
	}

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Database.RedisURI})

	// This main will later block until a value is sent to the "quitChan" channel.
	// Value will be sent when:
	// - An OS signal like SIGTERM is sent to the process
	// - When one of the HTTP servers fails
	quitChan := make(chan int, 1)
	quitting := false

	// Create & start the monitoring HTTP server. When it stops, send an event to the "quitChan" channel to stop the app
	monitoringMux := http.NewServeMux()
	monitoringMux.Handle("/metrics", promhttp.Handler())
	monitoringServer := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Monitoring.Port),
		Handler:           monitoringMux,
		ReadTimeout:       cfg.Monitoring.ReadTimeout,
		ReadHeaderTimeout: cfg.Monitoring.ReadHeaderTimeout,
		WriteTimeout:      cfg.Monitoring.WriteTimeout,
		IdleTimeout:       cfg.Monitoring.IdleTimeout,
		MaxHeaderBytes:    cfg.Monitoring.MaxHeaderBytes,
	}
	go func() {
		if err := monitoringServer.ListenAndServe(); err != nil && err != http.ErrServerClosed && !quitting {
			log.Warn().Err(err).Msg("Monitoring HTTP server failed")
			quitChan <- ExitCodeListen
		} else {
			quitChan <- ExitCodeOK
		}
	}()

	// Create & start the service HTTP server. When it stops, send an event to the "quitChan" channel to stop the app
	router := infra.NewChiRouter(cfg.HTTP)
	router.Route("/", rest.NewRoutes(accounts.NewManager(mongoClient), transactions.NewManager(mongoClient)))
	serviceServer := &http.Server{
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
	}
	go func() {
		if err := serviceServer.ListenAndServe(); err != nil && err != http.ErrServerClosed && !quitting {
			log.Warn().Err(err).Msg("Service HTTP server failed")
			quitChan <- ExitCodeListen
		} else {
			quitChan <- ExitCodeOK
		}
	}()

	// Create & start the Pub/Sub transactions consumer; when it stops, send an event to the "quitChan" channel to stop
	// the entire app
	txPubSubConsumer, shutdownTxPubSubConsumer, err := pubsub.NewTransactionProcessor(redisClient, cfg.PubSub)
	if err != nil {
		log.Error().Err(err).Msg("Failed creating transactions pub/sub consumer")
		return ExitCodeListen
	}
	defer shutdownTxPubSubConsumer(context.Background())
	go func() {
		if err := txPubSubConsumer.Run(); err != nil && !quitting {
			log.Warn().Err(err).Msg("Transactions Pub/Sub consumer failed")
			quitChan <- ExitCodeListen
		} else {
			quitChan <- ExitCodeOK
		}
	}()

	// Install a handler for critical OS signals, which, if called, will send a value to the "quitChan" channel.
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(
		shutdownChan,
		syscall.SIGINT,  // Triggered when user wants to interrupt the process (usually via "Ctrl+C"); termination is usually expected here.
		syscall.SIGTERM, // Triggered when user wants to terminate the process (not sure how; probably programmatic)
		syscall.SIGQUIT, // Triggered when user wants to terminate the process w/ a core dump (usually via "Ctrl+\" or programmatic)
		syscall.SIGKILL, // Will never be caught be the program - this is triggered by "kill -9" but we won't receive it
	)
	go func() {
		sig := <-shutdownChan
		log.Info().Str("signal", sig.String()).Msg("Received OS signal")
		quitting = true
		quitChan <- 0
	}()

	// If we've reached this far, we can cancel the startup context
	cancelStartupCtx()

	// Now let's wait until we are told to quit. This can happen through one of the following:
	// - An OS signalling us to stop (e.g. CTRL+C)
	// - One of the HTTP servers failing
	// - One of the Pub/Sub consumers failing
	exitCode = <-quitChan
	log.Info().Int("exitCode", exitCode).Msg("Done")
	return
}
