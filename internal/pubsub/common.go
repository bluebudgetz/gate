package pubsub

import (
	"context"
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/go-redis/redis"
	"github.com/golangly/errors"

	"github.com/golangly/log"
)

type (
	MessageProcessor func(m *pubsub.Message) error

	Consumer interface {
		io.Closer
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

func newConsumer(redisClient *redis.Client, processor MessageProcessor, maxMessageFailures uint8, deadLetterboxTopicName, subscriptionPath string) (Consumer, error) {
	if processor == nil {
		return nil, errors.New("message processor is required")
	} else if subscriptionPath == "" {
		return nil, errors.New("subscription path is required")
	}

	subTokens := strings.Split(subscriptionPath, "/")
	if len(subTokens) != 3 {
		return nil, errors.Newf("invalid subscription path").AddTag("subscriptionPath", subscriptionPath)
	}

	projectID := subTokens[0]
	topicName := subTokens[1]
	subscriptionName := subTokens[2]
	if pubsubClient, err := pubsub.NewClient(context.Background(), projectID); err != nil {
		return nil, errors.Wrapf(err, "failed creating Pub/Sub client").AddTag("projectID", projectID)
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
		}, nil
	}
}

func (p *consumer) Close() error { return p.pubsubClient.Close() }

func (p *consumer) getOrCreateTopic(ctx context.Context, name string) (*pubsub.Topic, error) {
	topic := p.pubsubClient.Topic(name)
	if exists, err := topic.Exists(ctx); err != nil {
		return nil, errors.Wrap(err, "failed checking if topic exists").
			AddTag("projectID", p.projectID).
			AddTag("topic", name)

	} else if exists {
		return topic, nil

	} else if topic, err := p.pubsubClient.CreateTopic(ctx, name); err != nil {
		return nil, errors.Wrap(err, "failed creating topic").
			AddTag("projectID", p.projectID).
			AddTag("topic", name)

	} else {
		return topic, nil
	}
}

func (p *consumer) getOrCreateSubscription(ctx context.Context, topicName string, subscriptionName string) (*pubsub.Topic, *pubsub.Subscription, error) {
	topic, err := p.getOrCreateTopic(ctx, topicName)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed getting subscription").
			AddTag("projectID", p.projectID).
			AddTag("topic", topicName).
			AddTag("subscription", subscriptionName)
	}

	subscription := p.pubsubClient.Subscription(subscriptionName)
	exists, err := subscription.Exists(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed checking if subscription exists").
			AddTag("projectID", p.projectID).
			AddTag("topic", topicName).
			AddTag("subscription", subscriptionName)
	}

	if !exists {
		subscription, err = p.pubsubClient.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{Topic: topic})
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed creating subscription").
				AddTag("projectID", p.projectID).
				AddTag("topic", topicName).
				AddTag("subscription", subscriptionName)
		}
	}

	subscriptionConfig, err := subscription.Config(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed getting subscription configuration").
			AddTag("projectID", p.projectID).
			AddTag("topic", topicName).
			AddTag("subscription", subscriptionName)
	}

	if subscriptionConfig.Topic.ID() != topicName {
		return nil, nil, errors.Wrap(err, "subscription belongs to wrong topic").
			AddTag("projectID", p.projectID).
			AddTag("topicFound", subscriptionConfig.Topic.ID()).
			AddTag("topicExpected", topicName).
			AddTag("subscription", subscriptionName)
	} else {
		return subscriptionConfig.Topic, subscription, nil
	}
}

func (p *consumer) increaseMessageFailureCounter(messageID string) (uint8, error) {
	key := fmt.Sprintf("%s:%s:%s:%s", p.projectID, p.topicName, p.subscriptionName, messageID)
	count, err := p.redisClient.Incr(key).Result()
	if err != nil {
		return math.MaxUint8, errors.Wrap(err, "failed increasing message failure counter").
			AddTag("projectID", p.projectID).
			AddTag("topic", p.topicName).
			AddTag("subscription", p.subscriptionName).
			AddTag("msgID", messageID)
	}

	// expire failures counter one minute from last failure
	if result := p.redisClient.Expire(key, time.Minute*1); result.Err() != nil {
		log.WithErr(err).
			With("projectID", p.projectID).
			With("topic", p.projectID).
			With("subscription", p.projectID).
			Warn("Failed expiring failure counter for a Pub/Sub message - this might result in more messages sent to the dead letterbox.")
	}
	return uint8(count), nil
}

func (p *consumer) Run() error {
	ctx := context.Background()

	// Lookup a reference to the dead-letter topic
	deadLetterTopic, err := p.getOrCreateTopic(ctx, p.deadLetterboxTopicName)
	if err != nil {
		return errors.Wrap(err, "failed getting dead letterbox topic").
			AddTag("projectID", p.projectID).
			AddTag("topic", p.deadLetterboxTopicName)
	}
	defer deadLetterTopic.Stop()

	// Lookup or create the subscription
	topic, subscription, err := p.getOrCreateSubscription(ctx, p.topicName, p.subscriptionName)
	if err != nil {
		return errors.Wrap(err, "failed getting consumer's subscription").
			AddTag("projectID", p.projectID).
			AddTag("topic", p.topicName).
			AddTag("subscription", p.subscriptionName)
	}
	defer topic.Stop()

	// Subscribe and start receive messages
	log.With("projectID", p.projectID).
		With("topic", p.topicName).
		With("subscription", p.subscriptionName).
		Info("Subscribing to Pub/Sub topic")
	return subscription.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		logger := log.
			With("topic", p.topicName).
			With("subscription", p.subscriptionName).
			With("msgID", m.ID).
			With("msgAttributes", m.Attributes)
		if err := p.processor(m); err != nil {
			logger.WithErr(err).Error("Failed processing message")

			// Message processing failed:
			//  - increase failure count for this message ID
			//  - if number of failures reached the maximum allowed number of failures, stop & send to dead-letter box
			//  - otherwise, NAck it to try again
			if failures, err := p.increaseMessageFailureCounter(m.ID); err != nil {
				logger.WithErr(err).Warn("Failed increasing failure count for message")
				m.Nack()
			} else if failures < p.maxMessageFailures {
				m.Nack()
			} else {
				logger.With("failures", failures).WithErr(err).Warn("Sending message to dead-letter topic")
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
			logger.Info("Processed message")
			m.Ack()
		}
	})
}
