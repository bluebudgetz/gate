package rest

import (
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bluebudgetz/gate/internal/rest/accounts"
	"github.com/bluebudgetz/gate/internal/rest/transactions"
)

func NewRoutes(mongoClient *mongo.Client) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/routes", routesDoc(r))
		r.Route("/accounts", accounts.NewRoutes(mongoClient))
		r.Route("/transactions", transactions.NewRoutes(mongoClient))
	}
}
