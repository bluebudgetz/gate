package util

import (
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/bluebudgetz/gate/internal/app"
	"github.com/bluebudgetz/gate/internal/boot"
	"github.com/bluebudgetz/gate/internal/config"
	"github.com/bluebudgetz/gate/internal/services"
)

func RunTestApp(t *testing.T, routesFactory func(neo4jDriver neo4j.Driver) func(chi.Router)) (*app.App, func()) {
	// Create application
	application, cleanup, err := boot.InitializeApp(config.NewTestConfig, routesFactory, services.NewTestNeo4jDriver)
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
