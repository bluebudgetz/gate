package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"

	"github.com/bluebudgetz/gate/internal/config"
	"github.com/bluebudgetz/gate/internal/infra"
	"github.com/bluebudgetz/gate/internal/rest"
	"github.com/bluebudgetz/gate/internal/rest/accounts"
)

func main() {

	// We must FIRST configure logging properly (pretty/JSON, stdout/stderr, etc)
	infra.SetupLogging()

	// Parse environment variables and/or command-line arguments, to form a Config object
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed loading configuration")
	}

	// Context for bootstrapping
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to MongoDB
	var mongoClient *mongo.Client
	if client, shutdown, err := infra.NewMongoClient(ctx, cfg.Database); err != nil {
		log.Fatal().Err(err).Msg("Failed creating MongoDB client")
	} else {
		mongoClient = client
		defer shutdown(ctx)
		var _ = mongoClient
	}

	// Connect to Redis
	var redisClient *redis.Client
	if client, shutdown, err := infra.NewRedisClient(cfg.Database); err != nil {
		log.Fatal().Err(err).Msg("Failed creating Redis client")
	} else {
		redisClient = client
		defer shutdown(ctx)
		var _ = redisClient
	}

	// Start monitoring HTTP server
	if server, shutdown, err := infra.NewMonitoringHTTPServer(cfg.Monitoring); err != nil {
		log.Fatal().Err(err).Msg("Failed creating monitoring HTTP server")
	} else {
		defer shutdown(ctx)
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Warn().Err(err).Msg("Monitoring HTTP server failed")
			}
		}()
	}

	// Managers
	accountsMgr := accounts.NewManager(mongoClient)

	// Create Chi router
	router := infra.NewChiRouter(*cfg)
	router.Route("/", rest.NewRoutes(accountsMgr))

	// Start the HTTP server using our router
	serviceServer := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler:           router,
		ReadTimeout:       time.Duration(cfg.HTTP.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.HTTP.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(cfg.HTTP.IdleTimeout) * time.Second,
		MaxHeaderBytes:    cfg.HTTP.MaxHeaderBytes,
		BaseContext: func(l net.Listener) context.Context {
			var ctx context.Context
			ctx = context.Background()
			ctx = context.WithValue(ctx, "validator", validator.New())
			return ctx
		},
	}
	if err := serviceServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Warn().Err(err).Msg("Service HTTP server failed")
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
