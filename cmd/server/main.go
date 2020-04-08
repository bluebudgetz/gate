package main

import (
	"context"
	glog "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis"
	"github.com/golangly/errors"
	"github.com/golangly/log"
	"github.com/golangly/webutil"
	"github.com/jessevdk/go-flags"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/go-playground/validator.v9"

	"github.com/bluebudgetz/gate/internal"
	"github.com/bluebudgetz/gate/internal/pubsub"
	"github.com/bluebudgetz/gate/internal/rest"
)

func main() {

	// We must FIRST configure logging properly (pretty/JSON, stdout/stderr, etc)
	log.Root = log.Root.With("svc", "gate")
	glog.SetFlags(0)
	glog.SetOutput(log.Root.Writer())

	// Defer a panic handler
	defer func() {
		if r := recover(); r != nil {
			log.WithPanic(r).Fatal("System panic!")
		}
	}()

	if err := mainWithReturnCode(); err != nil {
		log.WithErr(err).Error(err.Error())
		switch exitCode := errors.LookupTag(err, "exitCode").(type) {
		case int:
			os.Exit(exitCode)
		default:
			os.Exit(1)
		}
	} else {
		os.Exit(0)
	}
}

func mainWithReturnCode() (err error) {

	// Parse environment variables and/or command-line arguments, to form a Config object
	cfg := internal.Config{}
	parser := flags.NewParser(&cfg, flags.HelpFlag|flags.PrintErrors|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"
	parser.LongDescription = "Bluebudgetz API gateway. This is the API micro-service centralizing Bluebudgetz APIs."
	if _, err := parser.Parse(); err != nil {
		if parseErr, ok := err.(*flags.Error); ok {
			if parseErr.Type == flags.ErrHelp {
				return nil
			}
		}
		return errors.Wrap(err, err.Error())
	}

	// Print configuration for reference & debugging
	if skipValue, ok := os.LookupEnv("SKIP_PRINT_CONFIG"); !ok || (skipValue != "1" && skipValue != "y" && skipValue != "yes" && skipValue != "true") {
		log.With("config", cfg).
			With("env", os.Environ()).
			With("args", os.Args[1:]).
			Info("Configuration loaded")
	}

	// Connect to Neo4j
	neo4jDriver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.NoAuth())
	if err != nil {
		return errors.Wrap(err, "failed creating Neo4j client")
	}
	defer neo4jDriver.Close()

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Database.RedisURI})

	// Create the transactions Pub/Sub consumer
	txPubSubConsumer, err := pubsub.NewTransactionProcessor(redisClient, cfg.PubSub)
	if err != nil {
		return errors.Wrap(err, "failed creating transactions pub/sub consumer")
	}

	// Create & start the monitoring HTTP server
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

	// Create the HTTP service router
	router := chi.NewRouter()
	router.Use(
		middleware.SetHeader("Server", "bluebudgetz/gate"),
		middleware.Heartbeat("/health"),
		middleware.RealIP,
		webutil.RequestLogger,
		webutil.Metrics,
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
	router.Route("/", rest.NewRoutes(neo4jDriver))

	// Create the service HTTP server
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

	// Define a channel signaling the termination of the server. The value is an error, but can be nil, and represents
	// the reason the server terminates. In essence, an error or nil is sent on:
	//   - An OS signal like SIGTERM is sent to the process
	//   - The monitoring or service HTTP servers fail
	//   - One of the Pub/Sub consumers fail
	// In any of these cases, an error is sent to this channel, and this "main" will quit.
	quitChan := make(chan error, 1)

	// This bool is turned on when a shutdown initiates. Allows ignoring subsequent errors after shutdown has initiated
	quitting := false

	// Start the monitoring HTTP server goroutine
	go func() {
		if err := monitoringServer.ListenAndServe(); err != nil && err != http.ErrServerClosed && !quitting {
			quitChan <- errors.Wrap(err, "Monitoring HTTP server failed")
		} else {
			quitChan <- nil
		}
	}()

	// Start the service HTTP server goroutine
	go func() {
		if err := serviceServer.ListenAndServe(); err != nil && err != http.ErrServerClosed && !quitting {
			quitChan <- errors.Wrap(err, "Service HTTP server failed")
		} else {
			quitChan <- nil
		}
	}()

	// Start the transactions pub/sub consumer
	go func() {
		if err := txPubSubConsumer.Run(); err != nil && !quitting {
			quitChan <- errors.Wrap(err, "Transactions Pub/Sub consumer failed")
		} else {
			quitChan <- nil
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
		log.With("signal", sig.String()).Info("Received OS signal")
		quitting = true
		quitChan <- nil
	}()

	// Now let's wait until we are told to quit. This can happen through one of the following:
	// - An OS signalling us to stop (e.g. CTRL+C)
	// - One of the HTTP servers failing
	// - One of the Pub/Sub consumers failing
	return <-quitChan
}
