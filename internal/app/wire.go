//+build wireinject

package app

import (
	"github.com/google/wire"

	"github.com/bluebudgetz/gate/internal/config"
	"github.com/bluebudgetz/gate/internal/server"
	"github.com/bluebudgetz/gate/internal/server/handlers"
	"github.com/bluebudgetz/gate/internal/services"
)

func InitializeApp() (*App, func(), error) {
	wire.Build(
		NewApp,
		config.NewConfig,
		services.NewNeo4jDriver,
		handlers.NewRoutes,
		server.NewPrometheusRegistry,
		server.NewChiRouter,
		server.NewHTTPServers,
	)
	return &App{}, nil, nil
}

func InitializeTestApp() (*App, func(), error) {
	wire.Build(
		NewApp,
		config.NewTestConfig,
		services.NewTestNeo4jDriver,
		handlers.NewRoutes,
		server.NewPrometheusRegistry,
		server.NewChiRouter,
		server.NewHTTPServers,
	)
	return &App{}, nil, nil
}
