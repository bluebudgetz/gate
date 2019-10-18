package infra

import (
	"context"
	"github.com/bluebudgetz/gate/internal/config"
	"github.com/go-redis/redis"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(ctx context.Context, cfg config.DatabaseConfig) (*mongo.Client, func(context.Context), error) {
	var mongoClient *mongo.Client
	if client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MongoURI)); err != nil {
		return nil, nil, err
	} else if err := client.Connect(ctx); err != nil {
		// TODO: this means that if MongoDB is not available, we cannot start new Gate instances; a better behavior
		// 	     would be to return errors in runtime, so we can return service-not-available responses instead. This
		//       might affect auto-scaling as well :(
		return nil, nil, err
	} else {
		return client, func(ctx context.Context) {
			if err := mongoClient.Disconnect(ctx); err != nil {
				log.Warn().Err(err).Msg("Failed disconnecting from MongoDB")
			}
		}, nil
	}
}

func NewRedisClient(cfg config.DatabaseConfig) (*redis.Client, func(context.Context), error) {
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.RedisURI})
	return redisClient, func(ctx context.Context) {
		if err := redisClient.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed disconnecting from Redis")
		}
	}, nil
}
