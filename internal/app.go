package internal

import (
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/handler"
	"github.com/bluebudgetz/common/pkg/config"
	"github.com/bluebudgetz/common/pkg/logging"
	"github.com/bluebudgetz/gate/internal/assets"
	"github.com/bluebudgetz/gate/internal/graphql"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	Log  logging.LogConfig
	Http struct {
		Port uint16
	}
	Db struct {
		Username string
		Password string
		Host     string
		Port     uint16
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
	v.SetDefault("db.username", "")
	v.SetDefault("db.password", "")
	v.SetDefault("db.host", "localhost")
	v.SetDefault("db.port", 3306)
	v.SetDefault("db.sqlPath", "/db")
	err := v.Unmarshal(&conf)
	if err != nil {
		panic(errors.Wrap(err, "failed reading configuration"))
	} else if conf.Db.Host == "" {
		panic(errors.New("database host is required"))
	}

	// Setup logging
	logging.ConfigureLogger(&conf.Log)

	// Print config
	logging.Log.Infof("Configuration: %s", spew.Sdump(conf))

	// Setup database connection pool
	dbUrl := fmt.Sprintf(
		"%s@tcp(%s:%d)/bb?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		buildDbCredentials(&conf),
		conf.Db.Host,
		conf.Db.Port,
	)
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		panic(errors.Wrap(err, "failed creating connection pool"))
	}

	// Initialize database when in development mode
	if config.GetEnvironment() == config.Dev {
		// TODO: create a "scripts/migrate.go" file for this

		driver, err := mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			panic(errors.Wrap(err, "failed creating database migration driver"))
		}

		// Extract all migrations to a temporary directory
		tempMigrationsPath, err := ioutil.TempDir("", "bluebudgetzMigrations")
		if err != nil {
			panic(errors.Wrap(err, "failed extracting database migration files"))
		}
		defer os.RemoveAll(tempMigrationsPath)
		if err = assets.RestoreAssets(tempMigrationsPath, "deployments/rdbms/migrations"); err != nil {
			panic(errors.Wrap(err, "failed extracting database migration files"))
		}

		// Migrate all the way down, and then all the way up
		migrator, err := migrate.NewWithDatabaseInstance("file://"+tempMigrationsPath+"/deployments/rdbms/migrations", "mysql", driver)
		if err != nil {
			panic(errors.Wrap(err, "failed creating database migrator"))
		}
		if err = migrator.Down(); err != nil && err != migrate.ErrNoChange {
			panic(errors.Wrap(err, "failed to drop database"))
		}
		if err = migrator.Up(); err != nil && err != migrate.ErrNoChange {
			panic(errors.Wrap(err, "failed to migrate database"))
		}

		// Run data initialization
		if _, err = db.Exec(string(assets.MustAsset("assets/rdbms/init.sql"))); err != nil {
			panic(errors.Wrap(err, "failed initializing schema"))
		}
	}

	// Setup router
	router := chi.NewRouter()
	router.Use(NewPersistenceMiddleware(db))
	router.Handle("/query", createGraphQLHandler(db))
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

func buildDbCredentials(conf *Config) string {
	if conf.Db.Username != "" && conf.Db.Password != "" {
		return conf.Db.Username + ":" + conf.Db.Password
	} else if conf.Db.Username != "" {
		return conf.Db.Username
	} else if conf.Db.Password != "" {
		return conf.Db.Password
	} else {
		return ""
	}
}

func createGraphQLHandler(db *sql.DB) http.HandlerFunc {
	resolver, err := NewResolverRoot(db)
	if err != nil {
		panic(err)
	}
	return handler.GraphQL(
		graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}),
		handler.IntrospectionEnabled(config.GetEnvironment() != config.Prod),
		handler.ComplexityLimit(50),
	)
}
