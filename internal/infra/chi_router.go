package infra

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

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

func NewChiRouter(cfg HTTPConfig) chi.Router {
	r := chi.NewRouter()
	r.Use(
		middleware.SetHeader("Server", "bluebudgetz/gate"),
		middleware.Heartbeat("/health"),
		middleware.RealIP,
		RequestLogger,
		chiMetrics,
		chiRequestID,
		middleware.NoCache,
		corsHandler(cfg.CORS),
		middleware.GetHead,
		middleware.RedirectSlashes,
		middleware.Compress(cfg.GZipLevel),
		middleware.Timeout(30*time.Second),
	)
	return r
}

func corsHandler(cfg HTTPCORSConfig) func(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowOrigins,
		AllowedMethods:   cfg.AllowMethods,
		AllowedHeaders:   cfg.AllowHeaders,
		ExposedHeaders:   cfg.ExposeHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           int(cfg.MaxAge.Seconds()), // 300 is the maximum value not ignored by any of major browsers
	}).Handler
}
