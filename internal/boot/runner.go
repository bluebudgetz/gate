package boot

import (
	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/bluebudgetz/gate/internal/app"
	"github.com/bluebudgetz/gate/internal/config"
	"github.com/bluebudgetz/gate/internal/server"
	"github.com/bluebudgetz/gate/internal/server/metrics"
)

func InitializeApp(
	configFactory func() (config.Config, error),
	routesFactory func(driver neo4j.Driver) func(chi.Router),
	neo4jFactory func() (neo4j.Driver, func(), error)) (*app.App, func(), error) {
	configConfig, err := configFactory()
	if err != nil {
		return nil, nil, err
	}

	prometheusRegistry := metrics.NewPrometheusRegistry()

	neo4jDriver, cleanup, err := neo4jFactory()
	if err != nil {
		return nil, nil, err
	}

	routes := routesFactory(neo4jDriver)

	router := server.NewChiRouter(configConfig, prometheusRegistry, routes)

	servers := server.NewHTTPServers(configConfig, router)

	application, err := app.NewApp(configConfig, neo4jDriver, servers)
	if err != nil {
		cleanup()
		return nil, nil, err
	} else {
		return application, func() { cleanup() }, nil
	}
}
