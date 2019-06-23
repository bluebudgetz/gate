package mongo

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func CreateMongoClient(uri string) (*mongo.Client, error) {
	if client, err := mongo.NewClient(options.Client().ApplyURI(uri)); err != nil {
		return nil, errors.Wrapf(err, "failed connecting to MongoDB")
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Connect(ctx); err != nil {
			return nil, errors.Wrapf(err, "failed connecting to MongoDB")
		} else {
			return client, nil
		}
	}
}

func ObjectIdFromNillableString(any interface{}) *primitive.ObjectID {
	if any != nil {
		if stringValue, ok := any.(string); ok {
			if objectID, err := primitive.ObjectIDFromHex(stringValue); err != nil {
				log.Debug().Interface("any", any).Str("stringValue", stringValue).Msg("Failed creating ObjectID")
			} else {
				return &objectID
			}
		}
	}
	return nil
}
