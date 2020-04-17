// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package app

import (
	"github.com/bluebudgetz/gate/internal/config"
	"github.com/bluebudgetz/gate/internal/server"
	"github.com/bluebudgetz/gate/internal/server/handlers"
	"github.com/bluebudgetz/gate/internal/services"
)

// Injectors from wire.go:

func InitializeApp() (*App, func(), error) {
	configConfig, err := config.NewConfig()
	if err != nil {
		return nil, nil, err
	}
	registry := server.NewPrometheusRegistry()
	driver, cleanup, err := services.NewNeo4jDriver()
	if err != nil {
		return nil, nil, err
	}
	routes := handlers.NewRoutes(driver)
	mux := server.NewChiRouter(configConfig, registry, routes)
	v := server.NewHTTPServers(configConfig, mux)
	appApp, err := NewApp(configConfig, v)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return appApp, func() {
		cleanup()
	}, nil
}
