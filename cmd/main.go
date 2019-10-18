package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bluebudgetz/gate/internal/config"
	"github.com/bluebudgetz/gate/internal/infra"
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

	// Create HTTP router
	var router *gin.Engine
	if r, err := infra.NewRouter(*cfg); err != nil {
		log.Fatal().Err(err).Msg("Failed creating Gin router")
		return
	} else {
		router = r
	}

	// Add handlers
	// TODO: cacheStore := persistence.NewInMemoryStore(time.Second)
	router.GET("/accounts", accounts.GET(mongoClient))

	// Start the HTTP server using our router
	serviceServer := &http.Server{
		Addr:           ":" + strconv.Itoa(cfg.HTTP.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second, // TODO: externalize HTTP read timeout
		WriteTimeout:   10 * time.Second, // TODO: externalize HTTP write timeout
		MaxHeaderBytes: 1 << 20,          // TODO: externalize HTTP max header size
	}
	if err := serviceServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Warn().Err(err).Msg("Service HTTP server failed")
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
