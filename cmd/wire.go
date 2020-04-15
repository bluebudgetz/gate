//+build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/bluebudgetz/gate/internal/app"
	"github.com/bluebudgetz/gate/internal/config"
	"github.com/bluebudgetz/gate/internal/server"
	"github.com/bluebudgetz/gate/internal/server/handlers"
	"github.com/bluebudgetz/gate/internal/services"
)

func InitializeApp() (*app.App, func(), error) {
	wire.Build(
		app.NewApp,
		config.NewConfig,
		services.NewNeo4jDriver,
		handlers.NewRoutes,
		server.NewPrometheusRegistry,
		server.NewChiRouter,
		server.NewHTTPServers,
	)
	return &app.App{}, nil, nil
}
