package internal

import (
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/handler"
	"github.com/bluebudgetz/common/pkg/config"
	"github.com/bluebudgetz/common/pkg/logging"
	"github.com/bluebudgetz/gate/internal/graphql/impl"
	"github.com/bluebudgetz/gate/internal/graphql/resolver"
	"github.com/bluebudgetz/gate/internal/middleware"
	"github.com/bluebudgetz/gate/internal/migrator"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"net/http"
	"time"
)

type Config struct {
	Log  logging.LogConfig
	Http struct {
		Cors struct {
			Host string
			Port uint16
		}
		Port uint16
	}
	Db struct {
		Username    string
		Password    string
		Host        string
		Port        uint16
		MaxLifetime time.Duration
		MaxIdle     int
		MaxOpen     int
	}
}

type Gate interface {
	Config() Config
	Run() error
}

type gate struct {
	config Config
	db     *sql.DB
	router *chi.Mux
}

func NewGate() Gate {
	conf := Config{}

	// Setup Viper
	v := config.CreateViper("gate")
	v.SetDefault("http.port", 3001)
	v.SetDefault("http.cors.host", "www.bluebudgetz.com")
	v.SetDefault("http.cors.port", 80)
	v.SetDefault("db.username", "")
	v.SetDefault("db.password", "")
	v.SetDefault("db.host", "localhost")
	v.SetDefault("db.port", 5432)
	v.SetDefault("db.maxLifetime", 60 * time.Second)
	v.SetDefault("db.maxIdle", int(5))
	v.SetDefault("db.maxOpen", int(15))
	err := v.Unmarshal(&conf)
	if err != nil {
		panic(errors.Wrap(err, "failed reading configuration"))
	} else if conf.Db.Host == "" {
		panic(errors.New("database host is required"))
	}

	// Setup logging
	logging.ConfigureLogger(&conf.Log)
	logging.Log.Infof("Configuration: %s", spew.Sdump(conf))

	// Setup database connection pool
	dsn := fmt.Sprintf("host=%s port=%d sslmode=disable connect_timeout=30", conf.Db.Host, conf.Db.Port)
	if conf.Db.Username != "" {
		dsn = fmt.Sprintf("%s user=%s", dsn, conf.Db.Username)
	}
	if conf.Db.Password != "" {
		dsn = fmt.Sprintf("%s password=%s", dsn, conf.Db.Password)
	}
	logging.Log.Infof("Connecting to database at: %s", dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(errors.Wrap(err, "failed creating connection pool"))
	}
	db.SetConnMaxLifetime(conf.Db.MaxLifetime)
	db.SetMaxIdleConns(conf.Db.MaxIdle)
	db.SetMaxOpenConns(conf.Db.MaxOpen)

	// Initialize database when in development mode
	if config.GetEnvironment() == config.Dev {
		m, err := migrator.New(db)
		if err != nil {
			panic(err)
		}
		err = m.Migrate()
		if err != nil {
			panic(err)
		}
		err = m.Populate()
		if err != nil {
			panic(err)
		}
	}

	// Setup router
	logging.Log.Info("Setting up router")
	router := chi.NewRouter()
	router.Use(middleware.NewPersistenceMiddleware(db))
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{fmt.Sprintf("http://%s:%d", conf.Http.Cors.Host, conf.Http.Cors.Port)},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)
	router.Handle("/query", createGraphQLHandler())
	router.Get("/playground", handler.Playground("Gate | Bluebudgetz", "/query"))

	return &gate{
		config: conf,
		router: router,
	}
}

func (a *gate) Config() Config {
	return a.config
}

func (a *gate) Run() error {
	port := a.config.Http.Port
	logging.Log.Infof("Starting gate")
	return http.ListenAndServe(fmt.Sprintf(":%d", port), a.router)
}

func createGraphQLHandler() http.HandlerFunc {
	complexity := 0
	if config.GetEnvironment() == config.Prod {
		complexity = 100
	}
	return handler.GraphQL(
		impl.NewExecutableSchema(impl.Config{Resolvers: &resolver.Resolver{}}),
		handler.IntrospectionEnabled(config.GetEnvironment() != config.Prod),
		handler.ComplexityLimit(complexity),
	)
}
