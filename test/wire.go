//+build wireinject

package test

import (
	"github.com/google/wire"

	"github.com/bluebudgetz/gate/internal/app"
	"github.com/bluebudgetz/gate/internal/server"
	"github.com/bluebudgetz/gate/internal/server/handlers"
)

func InitializeApp() (*app.App, func(), error) {
	wire.Build(
		app.NewApp,
		NewConfig,
		NewNeo4jDriver,
		handlers.NewRoutes,
		server.NewPrometheusRegistry,
		server.NewChiRouter,
		server.NewHTTPServers,
	)
	return &app.App{}, nil, nil
}
