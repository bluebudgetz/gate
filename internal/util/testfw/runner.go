package testfw

import (
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/golangly/errors"
	"github.com/golangly/log"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/bluebudgetz/gate/internal/app"
	"github.com/bluebudgetz/gate/internal/boot"
	"github.com/bluebudgetz/gate/internal/config"
	"github.com/bluebudgetz/gate/internal/services"
)

func NewTestNeo4jDriver() (neo4j.Driver, func(), error) {
	driver, cleanup, err := services.NewNeo4jDriver()
	if err != nil {
		return nil, cleanup, err
	}

	session, err := services.CreateNeo4jSession(driver, neo4j.AccessModeWrite)
	if err != nil {
		return nil, cleanup, err
	}
	defer session.Close()

	// Cleanup database
	log.Info("Cleaning Neo4j database")
	if result, err := session.Run("MATCH (n) DETACH DELETE (n)", nil); err != nil {
		return nil, cleanup, errors.Wrap(err, "failed cleaning database")
	} else if err := result.Err(); err != nil {
		return nil, cleanup, errors.Wrap(err, "result error")
	}
	return driver, cleanup, err
}

func Run(t *testing.T, routesFactory func(neo4jDriver neo4j.Driver) func(chi.Router)) (*app.App, func()) {
	// Create application
	application, cleanup, err := boot.InitializeApp(config.NewTestConfig, routesFactory, NewTestNeo4jDriver)
	if err != nil {
		t.Fatalf("failed starting app: %+v", err)
	}

	// Run application
	quitChan := make(chan error, 10)
	go func() {
		if err := application.Run(quitChan); err != nil {
			t.Fatalf("failed running application: %+v", err)
		}
	}()

	// Sleep to give the app a chance to fully load
	time.Sleep(500 * time.Millisecond)

	// Return application & composite cleanup function
	return application, func() {
		quitChan <- nil
		cleanup()
	}
}
