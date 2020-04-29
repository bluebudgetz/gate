package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golangly/log"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/bluebudgetz/gate/internal/config"
)

// Signals causing the app to exit
var signals = []os.Signal{
	syscall.SIGINT,  // Triggered when user wants to interrupt the process (usually via "Ctrl+C"); termination is usually expected here.
	syscall.SIGTERM, // Triggered when user wants to terminate the process (not sure how; probably programmatic)
	syscall.SIGQUIT, // Triggered when user wants to terminate the process w/ a core dump (usually via "Ctrl+\" or programmatic)
	syscall.SIGKILL, // Will never be caught be the program - this is triggered by "kill -9" but we won't receive it
}

type App struct {
	cfg         config.Config
	neo4jDriver neo4j.Driver
	servers     []*http.Server
}

func (app *App) Neo4jDriver() neo4j.Driver {
	return app.neo4jDriver
}

// Run the application, until an error (or nil) is pushed to the given quit channel.
func (app *App) Run(quitChan chan error) error {

	// Start the HTTP servers
	for _, server := range app.servers {
		go func(svr *http.Server) {
			if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				quitChan <- err
			} else {
				quitChan <- nil
			}
		}(server)
	}

	// Install a handler for critical OS signals, which, if called, will send a value to the "quitChan" channel.
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, signals...)
	go func() {
		sig := <-shutdownChan
		log.With("signal", sig.String()).Info("Received OS signal")
		quitChan <- nil
	}()
	defer signal.Stop(shutdownChan)

	// Now let's wait until we are told to quit by the 'quitChan'. This can happen through one of the following:
	// - An OS signalling us to stop, e.g. CTRL+C signal to "shutdownChan" causing nil to be sent to "quitChan"
	// - One of the HTTP servers stopping, sending err or nil to "quitChan"
	err := <-quitChan
	for _, server := range app.servers {
		if err := server.Close(); err != nil {
			log.WithErr(err).Error("Failed closing HTTP server")
		}
	}

	// Return the error from the quitChan
	return err
}

// Returns the application configuration.
func (app *App) GetConfig() config.Config {
	return app.cfg
}

func (app *App) BuildURL(host string, path string, pathArgs ...interface{}) string {
	return fmt.Sprintf("http://%s:%d%s", host, app.GetConfig().HTTP.Port, fmt.Sprintf(path, pathArgs...))
}

func NewApp(cfg config.Config, neo4jDriver neo4j.Driver, servers []*http.Server) (*App, error) {
	return &App{cfg, neo4jDriver, servers}, nil
}
