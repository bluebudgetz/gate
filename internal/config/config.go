package config

import (
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
	"os"
)

type DatabaseConfig struct {
	MongoURI string `long:"mongo-uri" env:"MONGODB_URI" value-name:"URI" default:"mongodb://localhost:27017" description:"Mongo URI"`
	RedisURI string `long:"redis-uri" env:"REDIS_URI" value-name:"URI" default:"localhost:6379" description:"Redis URI"`
}

type MonitoringConfig struct {
	MetricsPort int `long:"metrics-port" env:"METRICS_PORT" value-name:"PORT" default:"3002" description:"Port to service Prometheus metrics"`
}

type HTTPCORSConfig struct {
	AllowOrigins     []string `long:"allow-origin" env:"ALLOW_ORIGIN" value-name:"ORIGIN" default:"https://www.bluebudgetz.com:80" description:"Origin to allow requests from, e.g. http://my.server.com"`
	AllowMethods     []string `long:"allow-method" env:"ALLOW_METHOD" value-name:"METHOD" default:"HEAD, GET, POST, PATCH, PUT, DELETE, CONNECT, OPTIONS, TRACE" description:"Methods allowed in CORS requests"`
	AllowHeaders     []string `long:"allow-header" env:"ALLOW_HEADER" value-name:"HEADER" description:"Headers allowed in CORS requests"`
	AllowCredentials bool     `long:"allow-credentials" env:"ALLOW_CREDENTIALS" description:"Whether to allow client code to access responses when credentials were sent in CORS requests"`
	ExposeHeaders    []string `long:"expose-header" env:"EXPOSE_HEADER" value-name:"HEADER" description:"Headers exposed to client browser code in CORS requests"`
	MaxAge           int64    `long:"max-age" env:"MAX_AGE" value-name:"SECONDS" default:"30" description:"How long (in seconds) can preflight responses be cached"`
}

type HTTPConfig struct {
	Port               int            `long:"port" default:"3001" env:"PORT" value-name:"PORT" description:"HTTP port to listen on"`
	CORS               HTTPCORSConfig `group:"CORS support" namespace:"cors" env-namespace:"CORS"`
	RequireHTTPS       bool           `long:"require-https" env:"REQUIRE_HTTPS" description:"Redirect HTTP requests to HTTPS"`
	DisableLogRequests bool           `long:"disable-request-log" env:"DISABLE_REQUEST_LOG" description:"Disable HTTP request log"`
	BodyLimit          string         `long:"body-limit" env:"BODY_LIMIT" default:"2M" description:"Maximum allowed size for a request body, e.g. 500K, 2M, 1G, etc"`
	GZipLevel          int            `long:"gzip-level" env:"GZIP_LEVEL" default:"-1" description:"HTTP GZip compression level"`
}

type PubSubConfig struct {
	GCPProjectID                 string `long:"project-id" env:"PROJECT_ID" value-name:"GCP_PROJECT_ID" description:"The GCP project ID where the Pub/Sub topics & subscriptions reside."`
	MaxMessageProcessingFailures uint8  `long:"max-process-failures" env:"MAX_PROCESS_FAILURES" value-name:"COUNT" default:"3" description:"Maximum number of message processing failures before forwarding a message to the dead letterbox."`
	DeadLetterboxTopic           string `long:"dead-letterbox-topic" env:"DEAD_LETTERBOX_TOPIC" value-name:"topic-name" default:"deadLetterBox" description:"Name of Pub/Sub dead letterbox topic."`
}

type Config struct {
	Database   DatabaseConfig   `group:"Databases" namespace:"db" env-namespace:"DB"`
	Monitoring MonitoringConfig `group:"Monitoring" namespace:"monitoring" env-namespace:"MONITORING"`
	HTTP       HTTPConfig       `group:"HTTP server" namespace:"http" env-namespace:"HTTP"`
	PubSub     PubSubConfig     `group:"Pub/Sub" namespace:"pubsub" env-namespace:"PUBSUB"`
}

func NewConfig() (*Config, error) {
	config := Config{}
	parser := flags.NewParser(&config, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"
	parser.LongDescription = "Bluebudgetz API gateway. This is the API micro-service centralizing Bluebudgetz APIs."
	if _, err := parser.Parse(); err != nil {
		return nil, err
	}

	// Print configuration early, before potential initialization errors will prevent this log (down the line..)
	log.Info().
		Interface("config", config).
		Interface("env", os.Environ()).
		Interface("args", os.Args).
		Msg("Configuration loaded")

	return &config, nil
}
