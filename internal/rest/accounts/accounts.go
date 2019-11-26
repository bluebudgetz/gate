package accounts

//go:generate go-bindata -o ./assets_gen.go -ignore ".*\\.go" -pkg accounts ./...

import (
	"encoding/json"
	"time"

	"github.com/go-chi/chi"
	"github.com/golangly/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Account struct {
		ID        string     `json:"id" yaml:"id"`
		CreatedOn time.Time  `json:"createdOn" yaml:"createdOn"`
		UpdatedOn *time.Time `json:"updatedOn" yaml:"updatedOn"`
		Name      string     `json:"name" yaml:"name"`
		ParentID  *string    `json:"parentId" yaml:"parentId"`
	}

	AccountWithBalance struct {
		Account
		TotalIncomingAmount float64 `json:"totalIncomingAmount" yaml:"totalIncomingAmount"`
		TotalOutgoingAmount float64 `json:"totalOutgoingAmount" yaml:"totalOutgoingAmount"`
		Balance             float64 `json:"balance" yaml:"balance"`
	}
)

var (
	getAccountsListQueryDoc bson.A
)

func init() {
	if err := json.Unmarshal(MustAsset("get_accounts_list_query.json"), &getAccountsListQueryDoc); err != nil {
		log.WithErr(err).Fatal("Failed loading 'get_accounts_query.json'")
	}
}

func coll(mongoClient *mongo.Client) *mongo.Collection {
	return mongoClient.Database("bluebudgetz").Collection("accounts")
}

func NewRoutes(mongoClient *mongo.Client) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", GetAccountList(mongoClient))
		r.Post("/", CreateAccount(mongoClient))
		r.Route("/{ID}", func(r chi.Router) {
			r.Get("/", GetAccount(mongoClient))
			r.Put("/", UpdateAccount(mongoClient))
			r.Patch("/", PatchAccount(mongoClient))
			r.Delete("/", DeleteAccount(mongoClient))
		})
	}
}
