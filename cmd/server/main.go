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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/go-playground/validator.v9"

	"github.com/bluebudgetz/gate/internal/pubsub"
	"github.com/bluebudgetz/gate/internal/rest"
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

type HTTPCORSConfig struct {
	AllowOrigins     []string      `long:"allow-origin" env:"ALLOW_ORIGIN" value-name:"ORIGIN" default:"https://www.bluebudgetz.com:443" description:"Origin to allow requests from, e.g. http://my.server.com"`
	AllowMethods     []string      `long:"allow-method" env:"ALLOW_METHOD" value-name:"METHOD" default:"HEAD, GET, POST, PATCH, PUT, DELETE, CONNECT, OPTIONS, TRACE" description:"Methods allowed in CORS requests"`
	AllowHeaders     []string      `long:"allow-header" env:"ALLOW_HEADER" value-name:"HEADER" description:"Headers allowed in CORS requests"`
	AllowCredentials bool          `long:"allow-credentials" env:"ALLOW_CREDENTIALS" description:"Whether to allow client code to access responses when credentials were sent in CORS requests"`
	ExposeHeaders    []string      `long:"expose-header" env:"EXPOSE_HEADER" value-name:"HEADER" description:"Headers exposed to client browser code in CORS requests"`
	MaxAge           time.Duration `long:"max-age" env:"MAX_AGE" value-name:"SECONDS" default:"30s" description:"How long (in seconds) can preflight responses be cached"`
}

type HTTPConfig struct {
	Port              int            `long:"port" env:"PORT" value-name:"PORT" default:"3001" description:"HTTP port to listen on"`
	CORS              HTTPCORSConfig `group:"CORS support" namespace:"cors" env-namespace:"CORS"`
	BodyLimit         string         `long:"body-limit" env:"BODY_LIMIT" default:"2M" description:"Maximum allowed size for a request body, e.g. 500K, 2M, 1G, etc"`
	GZipLevel         int            `long:"gzip-level" env:"GZIP_LEVEL" default:"-1" description:"HTTP GZip compression level"`
	ReadTimeout       time.Duration  `long:"read-timeout" env:"READ_TIMEOUT" default:"5s" description:"Maximum number of seconds to read the entire request, including the body"`
	ReadHeaderTimeout time.Duration  `long:"read-header-timeout" env:"READ_HEADER_TIMEOUT" default:"2s" description:"Maximum number of seconds to read the request headers"`
	WriteTimeout      time.Duration  `long:"write-timeout" env:"WRITE_TIMEOUT" default:"30s" description:"Maximum number of seconds to write the response"`
	IdleTimeout       time.Duration  `long:"idle-timeout" env:"IDLE_TIMEOUT" default:"30s" description:"Maximum number of seconds to let keep-alive connections to live"`
	MaxHeaderBytes    int            `long:"max-header-bytes" env:"MAX_HEADER_BYTES" default:"8192" description:"Maximum number of bytes to read for the request headers"`
}

type Config struct {
	Database   DatabaseConfig   `group:"Databases" namespace:"db" env-namespace:"DB"`
	Monitoring MonitoringConfig `group:"Monitoring" namespace:"monitoring" env-namespace:"MONITORING"`
	HTTP       HTTPConfig       `group:"HTTP server" namespace:"http" env-namespace:"HTTP"`
	PubSub     pubsub.Config    `group:"Pub/Sub" namespace:"pubsub" env-namespace:"PUBSUB"`
}

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
	cfg := Config{}
	parser := flags.NewParser(&cfg, flags.HelpFlag|flags.PrintErrors|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"
	parser.LongDescription = "Bluebudgetz API gateway. This is the API micro-service centralizing Bluebudgetz APIs."
	if _, err := parser.Parse(); err != nil {
		if parseErr, ok := err.(*flags.Error); ok {
			if parseErr.Type == flags.ErrHelp {
				return nil
			}
		}
		return err
	}

	// Print configuration for reference & debugging
	if skipValue, ok := os.LookupEnv("SKIP_PRINT_CONFIG"); !ok || (skipValue != "1" && skipValue != "y" && skipValue != "yes" && skipValue != "true") {
		log.With("config", cfg).
			With("env", os.Environ()).
			With("args", os.Args[1:]).
			Info("Configuration loaded")
	}

	// Context for bootstrapping
	startupCtx, cancelStartupCtx := context.WithTimeout(context.Background(), 30*time.Second)

	// Connect to MongoDB
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(cfg.Database.MongoURI))
	if err != nil {
		return errors.Wrap(err, "failed creating MongoDB client")
	} else if err := mongoClient.Connect(startupCtx); err != nil {
		return errors.Wrap(err, "failed connecting MongoDB client")
	}

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
	router.Route("/", rest.NewRoutes(mongoClient))

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

	// If we've reached this far, we can cancel the startup context
	cancelStartupCtx()
	cancelStartupCtx = nil

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
