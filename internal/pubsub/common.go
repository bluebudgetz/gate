package pubsub

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type (
	MessageProcessor func(m *pubsub.Message) error

	Consumer interface {
		Run() error
	}

	consumer struct {
		processor              MessageProcessor
		pubsubClient           *pubsub.Client
		redisClient            *redis.Client
		maxMessageFailures     uint8
		deadLetterboxTopicName string
		projectID              string
		topicName              string
		subscriptionName       string
	}
)

func newConsumer(redisClient *redis.Client, processor MessageProcessor, maxMessageFailures uint8, deadLetterboxTopicName, subscriptionPath string) (Consumer, func(context.Context), error) {
	if processor == nil {
		return nil, nil, errors.New("message processor is required")
	} else if subscriptionPath == "" {
		return nil, nil, errors.New("subscription path is required")
	}

	subTokens := strings.Split(subscriptionPath, "/")
	if len(subTokens) != 3 {
		return nil, nil, errors.New("invalid subscription path: " + subscriptionPath)
	}

	projectID := subTokens[0]
	topicName := subTokens[1]
	subscriptionName := subTokens[2]
	if pubsubClient, err := pubsub.NewClient(context.Background(), projectID); err != nil {
		return nil, nil, errors.New("failed creating Pub/Sub client")
	} else {
		return &consumer{
				processor,
				pubsubClient,
				redisClient,
				maxMessageFailures,
				deadLetterboxTopicName,
				projectID,
				topicName,
				subscriptionName,
			}, func(ctx context.Context) {
				if err := pubsubClient.Close(); err != nil {
					log.Warn().Err(err).Msg("Failed closing Pub/Sub client")
				}
			}, nil
	}
}

func (p *consumer) getOrCreateTopic(ctx context.Context, name string) (*pubsub.Topic, error) {
	topic := p.pubsubClient.Topic(name)
	if exists, err := topic.Exists(ctx); err != nil {
		return nil, errors.Wrapf(err, "failed checking if topic '%s' exists", name)

	} else if exists {
		return topic, nil

	} else if topic, err := p.pubsubClient.CreateTopic(ctx, name); err != nil {
		return nil, errors.Wrapf(err, "failed creating topic '%s'", name)

	} else {
		return topic, nil
	}
}

func (p *consumer) getOrCreateSubscription(ctx context.Context, topicName string, subscriptionName string) (*pubsub.Topic, *pubsub.Subscription, error) {
	subscription := p.pubsubClient.Subscription(subscriptionName)
	if exists, err := subscription.Exists(ctx); err != nil {
		return nil, nil, errors.Wrapf(err, "failed checking if subscription '%s' exists", subscriptionName)

	} else if exists {
		if subscriptionConfig, err := subscription.Config(ctx); err != nil {
			return nil, nil, errors.Wrapf(err, "failed fetching configuration of subscription '%s'", subscriptionName)
		} else {
			return subscriptionConfig.Topic, subscription, nil
		}

	} else if topic, err := p.getOrCreateTopic(ctx, topicName); err != nil {
		return nil, nil, errors.Wrapf(err, "failed getting topic '%s' for subscription '%s'", topicName, subscriptionName)

	} else if subscription, err = p.pubsubClient.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{Topic: topic}); err != nil {
		return nil, nil, errors.Wrapf(err, "failed creating subscription '%s' in topic '%s'", subscriptionName, topicName)

	} else {
		return topic, subscription, nil
	}
}

func (p *consumer) increaseMessageFailureCounter(messageID string) (uint8, error) {
	key := fmt.Sprintf("%s:%s:%s:%s", p.projectID, p.topicName, p.subscriptionName, messageID)
	count, err := p.redisClient.Incr(key).Result()
	if err != nil {
		return math.MaxUint8, err
	}
	p.redisClient.Expire(key, time.Minute*1) // expire failures counter one minute from last failure
	return uint8(count), nil
}

func (p *consumer) Run() error {
	ctx := context.Background()

	// Lookup a reference to the dead-letter topic
	deadLetterTopic, err := p.getOrCreateTopic(ctx, p.deadLetterboxTopicName)
	if err != nil {
		return err
	} else {
		defer deadLetterTopic.Stop()
	}

	// Lookup or create the subscription
	topic, subscription, err := p.getOrCreateSubscription(ctx, p.topicName, p.subscriptionName)
	if err != nil {
		return err
	} else {
		defer topic.Stop()
	}

	// Subscribe and start receive messages
	log.Info().
		Str("projectID", p.projectID).
		Str("topic", p.topicName).
		Str("subscription", p.subscriptionName).
		Msg("Subscribing to Pub/Sub topic")
	return subscription.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		info := map[string]interface{}{
			"msgID":        m.ID,
			"topic":        p.topicName,
			"subscription": p.subscriptionName,
			"attributes":   m.Attributes,
		}
		logger := log.With().Fields(info).Logger()
		if err := p.processor(m); err != nil {
			logger.Error().Err(err).Msg("Failed processing message")

			// Message processing failed:
			//  - increase failure count for this message ID
			//  - if number of failures reached the maximum allowed number of failures, stop & send to dead-letter box
			//  - otherwise, NAck it to try again
			if failures, err := p.increaseMessageFailureCounter(m.ID); err != nil {
				logger.Warn().Err(err).Msg("Failed increasing failure count for message")
				m.Nack()
			} else if failures < p.maxMessageFailures {
				m.Nack()
			} else {
				logger.Warn().Uint8("failures", failures).Err(err).Msg("Sending message to dead-letter topic")
				attributes := make(map[string]string, 0)
				attributes["originalTopic"] = p.topicName
				attributes["originalSubscription"] = p.subscriptionName
				if m.Attributes != nil {
					for k, v := range m.Attributes {
						attributes[k] = v
					}
				}
				deadLetterTopic.Publish(ctx, &pubsub.Message{Data: m.Data, Attributes: attributes})
				m.Ack()
			}
		} else {
			logger.Info().Str("msgID", m.ID).Msg("Processed message")
			m.Ack()
		}
	})
}
