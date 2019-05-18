package internal

import (
	"database/sql"
	"fmt"
	"github.com/bluebudgetz/common/pkg/config"
	"github.com/bluebudgetz/common/pkg/logging"
	"github.com/bluebudgetz/gate/internal/api/accounts"
	"github.com/bluebudgetz/gate/internal/api/transactions"
	"github.com/bluebudgetz/gate/internal/migrator"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"net/http"
	"time"
)

type Gate struct {
	Config struct {
		Log  logging.LogConfig
		Http struct {
			Cors struct {
				Host string
				Port uint16
			}
			Port uint16
		}
		DB struct {
			Postgres struct {
				DSN string
			}
			MaxLifetime time.Duration
			MaxIdle     int
			MaxOpen     int
		}
	}
	DB     *sql.DB
	Router *chi.Mux
}

type logrusLogEntry struct {
	Request *http.Request
}

func (e *logrusLogEntry) Write(status, bytes int, elapsed time.Duration) {
	logging.Log.WithFields(map[string]interface{}{
		"remoteAddr":    e.Request.RemoteAddr,
		"proto":         e.Request.Proto,
		"method":        e.Request.Method,
		"requestURI":    e.Request.RequestURI,
		"headers":       e.Request.Header,
		"url":           e.Request.URL,
		"host":          e.Request.Host,
		"contentLength": e.Request.ContentLength,
		"elapsed":       elapsed,
		"bytesWritten":  bytes,
		"status":        status,
	}).Trace("HTTP request executed")
}

func (e *logrusLogEntry) Panic(v interface{}, stack []byte) {
	if err, ok := v.(error); ok {
		logging.Log.
			WithError(err).
			WithFields(map[string]interface{}{
				"remoteAddr":    e.Request.RemoteAddr,
				"proto":         e.Request.Proto,
				"method":        e.Request.Method,
				"requestURI":    e.Request.RequestURI,
				"headers":       e.Request.Header,
				"url":           e.Request.URL,
				"host":          e.Request.Host,
				"contentLength": e.Request.ContentLength,
				// "stack":         string(stack),
			}).
			Error("HTTP request failed with a panic!")
	} else {
		logging.Log.
			WithFields(map[string]interface{}{
				"remoteAddr":    e.Request.RemoteAddr,
				"proto":         e.Request.Proto,
				"method":        e.Request.Method,
				"requestURI":    e.Request.RequestURI,
				"headers":       e.Request.Header,
				"url":           e.Request.URL,
				"host":          e.Request.Host,
				"contentLength": e.Request.ContentLength,
				"stack":         string(stack),
				"panic":         v,
			}).
			Error("HTTP request failed with a panic!")
	}
}

type logrusLogFormatter struct{}

func (*logrusLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &logrusLogEntry{Request: r}
}

func New() (*Gate, error) {
	_app := Gate{}

	// Setup Viper & read configuration
	v := config.CreateViper("gate")
	v.SetDefault("http.cors.host", "www.bluebudgetz.com")
	v.SetDefault("http.cors.port", 80)
	v.SetDefault("http.port", 3001)
	v.SetDefault("db.postgres.dsn", "") // example: host=localhost port=5432 sslmode=disable connect_timeout=30
	v.SetDefault("db.maxLifetime", 60*time.Second)
	v.SetDefault("db.maxIdle", int(5))
	v.SetDefault("db.maxOpen", int(15))
	err := v.Unmarshal(&_app.Config)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading configuration")
	}

	// Setup logging, and immediately print configuration
	logging.ConfigureLogger(&_app.Config.Log)
	logging.Log.Infof("Configuration: %s", spew.Sdump(_app.Config))

	// Validate configuration
	if _app.Config.DB.Postgres.DSN == "" {
		return nil, errors.New("database DSN is required")
	}

	// Setup database connection pool
	db, err := sql.Open("postgres", _app.Config.DB.Postgres.DSN)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating connection pool")
	}
	db.SetConnMaxLifetime(_app.Config.DB.MaxLifetime)
	db.SetMaxIdleConns(_app.Config.DB.MaxIdle)
	db.SetMaxOpenConns(_app.Config.DB.MaxOpen)
	_app.DB = db

	// Setup router
	logging.Log.Info("Setting up router")
	router := chi.NewRouter()
	router.Use(
		// First provide a "/health" endpoint
		middleware.Heartbeat("/health"),

		// Ensure request is uniquely identified & logged (with the real user IP)
		middleware.RequestID,
		middleware.RealIP,
		middleware.RequestLogger(&logrusLogFormatter{}),
		middleware.Recoverer,

		// Use "GET" handlers if "HEAD"-specific handlers are not found
		middleware.GetHead,

		// Apply common headers
		middleware.SetHeader("server", "bluebudgetz/gate"),
		middleware.NoCache,
		middleware.AllowContentType("application/json", ""),
		middleware.ContentCharset("", "UTF-8"),

		// Add CORS headers
		cors.New(cors.Options{
			AllowedOrigins:   []string{fmt.Sprintf("http://%s:%d", _app.Config.Http.Cors.Host, _app.Config.Http.Cors.Port)},
			AllowedMethods:   []string{"OPTIONS", "HEAD", "GET", "POST", "PATCH", "PUT", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}).Handler,

		// Set content-type
		render.SetContentType(render.ContentTypeJSON),
	)
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/accounts", accounts.New(_app.DB, int(_app.Config.Http.Port)).RoutesV1())
		r.Mount("/transactions", transactions.New(_app.DB).RoutesV1())
	})
	_app.Router = router

	// App is ready to start!
	return &_app, nil
}

func (_app *Gate) Run() error {
	httpErrChan := make(chan error)

	go func() {
		logging.Log.Infof("Starting gate")
		httpErrChan <- http.ListenAndServe(fmt.Sprintf(":%d", _app.Config.Http.Port), _app.Router)
	}()

	// Migrate the database
	if config.GetEnvironment() == config.Dev {

		// Give the HTTP server a chance to fully start
		time.Sleep(500 * time.Millisecond)

		m, err := migrator.New(_app.DB, int(_app.Config.Http.Port))
		if err != nil {
			return errors.Wrap(err, "failed creating database migrator")
		}
		err = m.Migrate()
		if err != nil {
			return errors.Wrap(err, "failed migrating database")
		}
		err = m.Populate()
		if err != nil {
			return errors.Wrap(err, "failed populating database")
		}
	}

	return <-httpErrChan
}
