package main

import (
	"bytes"
	"context"
	"io/ioutil"
	glog "log"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/golangly/errors"
	"github.com/golangly/log"
	"github.com/jessevdk/go-flags"
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
	Body       string            `long:"body" short:"b" value-name:"FILE" description:"File to read the message from. If omitted, body is read from stdin."`
}

func main() {

	// We must FIRST configure logging properly (pretty/JSON, stdout/stderr, etc)
	log.Printer = log.PrettyPrinter
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

// Runs the program and returns an error, potentially with an exit code attached to it.
// It's critical that the error is a golangly error - due to the way it's formatted.
func mainWithReturnCode() error {

	// Parse environment variables and/or command-line arguments, to form a Config object
	cfg := CLIConfig{}
	parser := flags.NewParser(&cfg, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"
	parser.LongDescription = "Bluebudgetz CLI."
	if _, err := parser.Parse(); err != nil {
		if parseErr, ok := err.(*flags.Error); ok {
			if parseErr.Type == flags.ErrHelp {
				return nil
			}
		}
		return errors.Wrap(err, err.Error())
	}

	// Execute
	switch parser.Active.Name {
	case "pubsub":
		return runPubSub(parser.Active, cfg.PubSub)
	default:
		return errors.Newf("unknown command: %s", parser.Active.Name)
	}
}

func runPubSub(cmd *flags.Command, cfg PubSubConfig) error {
	switch cmd.Active.Name {
	case "publish":
		return runPubSubPublish(cmd, cfg.Publish)
	default:
		return errors.Newf("unknown command: %s", cmd.Active.Name)
	}
}

func runPubSubPublish(_ *flags.Command, cfg PubSubPublishConfig) error {
	ctx := context.Background()

	topicTokens := strings.Split(cfg.Topic, "/")
	if len(topicTokens) != 2 {
		return errors.New("please specify topic in the format of '<project-id>/<topic-name>'")
	}

	projectID := topicTokens[0]
	if projectID == "" {
		return errors.New("please specify topic in the format of '<project-id>/<topic-name>'")
	}

	topicName := topicTokens[1]
	if topicName == "" {
		return errors.New("please specify topic in the format of '<project-id>/<topic-name>'")
	}

	pubsubClient, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		return errors.Wrap(err, "failed creating Pub/Sub client")
	}
	defer pubsubClient.Close()

	topic := pubsubClient.Topic(topicName)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, "failed checking if topic exists")
	} else if !exists {
		topic, err = pubsubClient.CreateTopic(ctx, topicName)
		if err != nil {
			return errors.Wrap(err, "failed creating topic")
		}
	}

	msg := &pubsub.Message{Attributes: cfg.Attributes}
	if cfg.Body == "" {
		body := new(bytes.Buffer)
		if _, err := body.ReadFrom(os.Stdin); err != nil {
			return errors.Wrap(err, "failed reading message from stdin")
		} else {
			msg.Data = body.Bytes()
		}
	} else if body, err := ioutil.ReadFile(cfg.Body); err != nil {
		return errors.Wrapf(err, "failed reading message from file: %s", cfg.Body)
	} else {
		msg.Data = body
	}

	id, err := topic.Publish(ctx, msg).Get(ctx)
	if err != nil {
		return errors.Wrap(err, "failed publishing message")
	}

	log.With("msgID", id).Info("Published message!")
	return nil
}
