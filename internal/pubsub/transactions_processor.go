package pubsub

import (
	"context"
	"mime"

	"github.com/pkg/errors"

	"cloud.google.com/go/pubsub"
	"github.com/go-redis/redis"
)

type Config struct {
	TxSubscription               string `long:"tx-subscription" env:"TX_SUBSCRIPTION" value-name:"PATH" description:"Subscription to fetch transaction processing requests from. Must be in the form of '<project-id>/<topic-id>/<subscription-name>'"`
	MaxMessageProcessingFailures uint8  `long:"max-process-failures" env:"MAX_PROCESS_FAILURES" value-name:"COUNT" default:"3" description:"Maximum number of message processing failures before forwarding a message to the dead letterbox."`
	DeadLetterboxTopic           string `long:"dead-letterbox-topic" env:"DEAD_LETTERBOX_TOPIC" value-name:"topic-name" default:"" description:"Name of Pub/Sub dead letterbox topic."`
}

func NewTransactionProcessor(redisClient *redis.Client, cfg Config) (Consumer, func(context.Context), error) {
	return newConsumer(
		redisClient,
		processTransactions,
		cfg.MaxMessageProcessingFailures,
		cfg.DeadLetterboxTopic,
		cfg.TxSubscription,
	)
}

func processTransactions(m *pubsub.Message) error {
	if contentTypeValue, ok := m.Attributes["content-type"]; !ok {
		return errors.New("message is missing 'content-type' header")
	} else if contentType, _, err := mime.ParseMediaType(contentTypeValue); err != nil {
		return errors.New("bad 'content-type' header encountered")
	} else if contentType == "application/vnd.ms-excel" {
		// TODO: implement support for *.xls transactions files
		return errors.New("opening *.xls is not yet implemented")
	} else if contentType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		// TODO: implement support for *.xlsx transactions files
		return errors.New("opening *.xlsx is not yet implemented")
	} else {
		return errors.New("unsupported content-type")
	}
}
