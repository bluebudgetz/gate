package main

import (
	glog "log"
	"os"

	"github.com/go-chi/chi"
	"github.com/golangly/errors"
	"github.com/golangly/log"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/bluebudgetz/gate/internal/boot"
	"github.com/bluebudgetz/gate/internal/config"
	"github.com/bluebudgetz/gate/internal/server/handlers"
	"github.com/bluebudgetz/gate/internal/services"
)

func main() {

	// We must FIRST configure logging properly (pretty/JSON, stdout/stderr, etc)
	log.Root = log.Root.With("svc", "gate")
	glog.SetFlags(0)
	glog.SetOutput(log.Root.Writer())

	// Defer a panic handler
	defer func() {
		if r := recover(); r != nil {
			log.WithPanic(r).Fatal("System panic!")
		}
	}()

	// Create application
	application, cleanup, err := boot.InitializeApp(
		config.NewConfig,
		func(neo4jDriver neo4j.Driver) func(chi.Router) { return handlers.NewRoutes(neo4jDriver) },
		services.NewNeo4jDriver,
	)
	if err != nil {
		handleError(err)
	}
	defer cleanup()

	// Print configuration for reference & debugging
	if skipValue, ok := os.LookupEnv("SKIP_PRINT_CONFIG"); !ok || (skipValue != "1" && skipValue != "y" && skipValue != "yes" && skipValue != "true") {
		log.With("config", application.GetConfig()).
			With("env", os.Environ()).
			With("args", os.Args[1:]).
			Info("Configuration loaded")
	}

	// Run application
	if err := application.Run(make(chan error, 10)); err != nil {
		handleError(err)
	}

	// Successful run
	os.Exit(0)
}

func handleError(err error) {
	log.WithErr(err).Error(err.Error())
	switch exitCode := errors.LookupTag(err, "exitCode").(type) {
	case int:
		os.Exit(exitCode)
	default:
		os.Exit(1)
	}
}
