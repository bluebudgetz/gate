package main

import (
	"context"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"

	"github.com/bluebudgetz/gate/internal/infra"
)

const (
	ExitCodeOK            = 0
	ExitCodeBadConfig     = 1
	ExitCodeInternalError = 2
)

type CLIConfig struct {
	PubSub PubSubConfig `command:"pubsub"`
}

type PubSubConfig struct {
	Publish PubSubPublishConfig `command:"publish"`
}

type PubSubPublishConfig struct {
	Topic      string            `long:"topic" env:"TOPIC" value-name:"PATH" required:"true" description:"Topic to publish the message to. Must be in the form of '<project-id>/<topic-id>'"`
	Attributes map[string]string `long:"attribute" short:"a" value-name:"VALUE" description:"Message attributes (can be specified multiple times.)"`
}

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	if err := os.Setenv("LOG_PRETTY", "1"); err != nil {
		log.Error().Err(err).Msg("Failed setting 'LOG_PRETTY' environment variable")
		return ExitCodeInternalError
	}

	// We must FIRST configure logging properly (pretty/JSON, stdout/stderr, etc)
	infra.SetupLogging()

	// Parse environment variables and/or command-line arguments, to form a Config object
	cfg := CLIConfig{}
	parser := flags.NewParser(&cfg, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"
	parser.LongDescription = "Bluebudgetz CLI."
	if _, err := parser.Parse(); err != nil {
		if parseErr, ok := err.(*flags.Error); ok {
			log.Error().Msg(parseErr.Error())
		} else {
			log.Error().Err(err).Msg("Failed loading configuration")
		}
		return ExitCodeBadConfig
	}

	// Defer a panic handler
	defer func() {
		if r := recover(); r != nil {
			log.Error().Str("stack", string(debug.Stack())).Interface("recovered", r).Msg("SYSTEM PANIC!")
			exitCode = ExitCodeInternalError
		}
	}()

	// Execute
	cmd := parser.Active
	if cmd == nil {
		log.Error().Msg("Please specify a command, or pass '-h' for help.")
		return ExitCodeBadConfig
	}
	switch cmd.Name {
	case "pubsub":
		return runPubSub(cmd, cfg.PubSub)
	default:
		log.Error().Msg("Unknown command!")
		return ExitCodeBadConfig
	}
}

func runPubSub(cmd *flags.Command, cfg PubSubConfig) int {
	if cmd.Active == nil {
		log.Error().Msg("Please specify a command, or pass '-h' for help.")
		return ExitCodeBadConfig
	}
	switch cmd.Active.Name {
	case "publish":
		return runPubSubPublish(cmd, cfg.Publish)
	default:
		log.Error().Msg("Unknown command!")
		return ExitCodeBadConfig
	}
}

func runPubSubPublish(_ *flags.Command, cfg PubSubPublishConfig) int {
	ctx := context.Background()

	topicTokens := strings.Split(cfg.Topic, "/")
	if len(topicTokens) != 2 {
		log.Error().Msg("Please specify topic in the format of '<project-id>/<topic-name>'")
		return ExitCodeBadConfig
	}

	projectID := topicTokens[0]
	if projectID == "" {
		log.Error().Msg("Please specify topic in the format of '<project-id>/<topic-name>'")
		return ExitCodeBadConfig
	}

	topicName := topicTokens[1]
	if topicName == "" {
		log.Error().Msg("Please specify topic in the format of '<project-id>/<topic-name>'")
		return ExitCodeBadConfig
	}

	pubsubClient, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		log.Error().Err(err).Msg("Failed creating Pub/Sub client")
		return ExitCodeInternalError
	}
	defer pubsubClient.Close()

	topic := pubsubClient.Topic(topicName)
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed checking if topic exists")
		return ExitCodeInternalError
	} else if !exists {
		topic, err = pubsubClient.CreateTopic(ctx, topicName)
		if err != nil {
			log.Error().Err(err).Msg("Failed creating topic")
			return ExitCodeInternalError
		}
	}

	msgText, err := ioutil.ReadAll(os.Stdin)
	result := topic.Publish(ctx, &pubsub.Message{Data: msgText, Attributes: cfg.Attributes})
	id, err := result.Get(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed publishing message")
		return ExitCodeInternalError
	}

	log.Info().Str("msgID", id).Msgf("Published message!")
	return ExitCodeOK
}
