package internal

import (
	"time"

	"github.com/bluebudgetz/gate/internal/pubsub"
)

type DatabaseConfig struct {
	MongoURI string `long:"mongo-uri" env:"MONGODB_URI" value-name:"URI" default:"mongodb://localhost:27017" description:"Mongo URI"`
	RedisURI string `long:"redis-uri" env:"REDIS_URI" value-name:"URI" default:"localhost:6379" description:"Redis URI"`
}

type MonitoringConfig struct {
	Port              int           `long:"metrics-port" env:"METRICS_PORT" value-name:"PORT" default:"3002" description:"Port to service Prometheus metrics"`
	ReadTimeout       time.Duration `long:"read-timeout" env:"READ_TIMEOUT" default:"3s" description:"Maximum number of seconds to read the entire request, including the body"`
	ReadHeaderTimeout time.Duration `long:"read-header-timeout" env:"READ_HEADER_TIMEOUT" default:"1s" description:"Maximum number of seconds to read the request headers"`
	WriteTimeout      time.Duration `long:"write-timeout" env:"WRITE_TIMEOUT" default:"5s" description:"Maximum number of seconds to write the response"`
	IdleTimeout       time.Duration `long:"idle-timeout" env:"IDLE_TIMEOUT" default:"30s" description:"Maximum number of seconds to let keep-alive connections to live"`
	MaxHeaderBytes    int           `long:"max-header-bytes" env:"MAX_HEADER_BYTES" default:"8192" description:"Maximum number of bytes to read for the request headers"`
}

type HTTPCORSConfig struct {
	AllowOrigins     []string      `long:"allow-origin" env:"ALLOW_ORIGIN" value-name:"ORIGIN" default:"https://www.bluebudgetz.com:443" description:"Origin to allow requests from, e.g. http://my.server.com"`
	AllowMethods     []string      `long:"allow-method" env:"ALLOW_METHOD" value-name:"METHOD" default:"HEAD, GET, POST, PATCH, PUT, DELETE, CONNECT, OPTIONS, TRACE" description:"Methods allowed in CORS requests"`
	AllowHeaders     []string      `long:"allow-header" env:"ALLOW_HEADER" value-name:"HEADER" description:"Headers allowed in CORS requests"`
	AllowCredentials bool          `long:"allow-credentials" env:"ALLOW_CREDENTIALS" description:"Whether to allow client code to access responses when credentials were sent in CORS requests"`
	ExposeHeaders    []string      `long:"expose-header" env:"EXPOSE_HEADER" value-name:"HEADER" description:"Headers exposed to client browser code in CORS requests"`
	MaxAge           time.Duration `long:"max-age" env:"MAX_AGE" value-name:"SECONDS" default:"30s" description:"How long (in seconds) can preflight responses be cached"`
}

type HTTPConfig struct {
	Port              int            `long:"port" env:"PORT" value-name:"PORT" default:"3001" description:"HTTP port to listen on"`
	CORS              HTTPCORSConfig `group:"CORS support" namespace:"cors" env-namespace:"CORS"`
	BodyLimit         string         `long:"body-limit" env:"BODY_LIMIT" default:"2M" description:"Maximum allowed size for a request body, e.g. 500K, 2M, 1G, etc"`
	GZipLevel         int            `long:"gzip-level" env:"GZIP_LEVEL" default:"-1" description:"HTTP GZip compression level"`
	ReadTimeout       time.Duration  `long:"read-timeout" env:"READ_TIMEOUT" default:"5s" description:"Maximum number of seconds to read the entire request, including the body"`
	ReadHeaderTimeout time.Duration  `long:"read-header-timeout" env:"READ_HEADER_TIMEOUT" default:"2s" description:"Maximum number of seconds to read the request headers"`
	WriteTimeout      time.Duration  `long:"write-timeout" env:"WRITE_TIMEOUT" default:"30s" description:"Maximum number of seconds to write the response"`
	IdleTimeout       time.Duration  `long:"idle-timeout" env:"IDLE_TIMEOUT" default:"30s" description:"Maximum number of seconds to let keep-alive connections to live"`
	MaxHeaderBytes    int            `long:"max-header-bytes" env:"MAX_HEADER_BYTES" default:"8192" description:"Maximum number of bytes to read for the request headers"`
}

type Config struct {
	Database   DatabaseConfig   `group:"Databases" namespace:"db" env-namespace:"DB"`
	Monitoring MonitoringConfig `group:"Monitoring" namespace:"monitoring" env-namespace:"MONITORING"`
	HTTP       HTTPConfig       `group:"HTTP server" namespace:"http" env-namespace:"HTTP"`
	PubSub     pubsub.Config    `group:"Pub/Sub" namespace:"pubsub" env-namespace:"PUBSUB"`
}
