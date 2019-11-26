package transactions

import (
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate go-bindata -o ./assets_gen.go -ignore ".*\\.go" -pkg transactions ./...

type (
	Transaction struct {
		ID              string     `json:"id" yaml:"id"`
		CreatedOn       time.Time  `json:"createdOn" yaml:"createdOn"`
		UpdatedOn       *time.Time `json:"updatedOn" yaml:"updatedOn"`
		IssuedOn        time.Time  `json:"issuedOn" yaml:"issuedOn"`
		Origin          string     `json:"origin" yaml:"origin"`
		SourceAccountID string     `json:"sourceAccountId" yaml:"sourceAccountId"`
		TargetAccountID string     `json:"targetAccountId" yaml:"targetAccountId"`
		Amount          float64    `json:"amount" yaml:"amount"`
		Comment         string     `json:"comment" yaml:"comment"`
	}
)

func coll(mongoClient *mongo.Client) *mongo.Collection {
	return mongoClient.Database("bluebudgetz").Collection("transactions")
}

func NewRoutes(mongoClient *mongo.Client) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", GetTransactionList(mongoClient))
		r.Post("/", CreateTransaction(mongoClient))
		r.Route("/{ID}", func(r chi.Router) {
			r.Get("/", GetTransaction(mongoClient))
			r.Put("/", UpdateTransaction(mongoClient))
			r.Patch("/", PatchTransaction(mongoClient))
			r.Delete("/", DeleteTransaction(mongoClient))
		})
	}
}
