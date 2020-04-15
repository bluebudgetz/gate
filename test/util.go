package test

import (
	"testing"
	"time"

	"github.com/bluebudgetz/gate/internal/app"
)

func Run(t *testing.T) (*app.App, func()) {

	// Create application
	application, cleanup, err := InitializeApp()
	if err != nil {
		t.Fatalf("failed creating application: %+v", err)
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
