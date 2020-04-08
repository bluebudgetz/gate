package rest

//go:generate go-bindata -o ./assets_gen.go -ignore ".*\\.go" -pkg rest ./...

import (
	"time"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/bluebudgetz/gate/internal/util"
)

var (
	deleteAccountQuery   string
	deleteTxQuery        string
	getAccountQuery      string
	getAccountsTreeQuery string
	getTxListQuery       string
	getTxQuery           string
	patchAccountQuery    string
	patchTxQuery         string
	postAccountQuery     string
	postTxQuery          string
	putAccountQuery      string
	putTxQuery           string
)

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

func init() {
	deleteAccountQuery = string(MustAsset("delete_account.cyp"))
	deleteTxQuery = string(MustAsset("delete_tx.cyp"))
	getAccountQuery = string(MustAsset("get_account.cyp"))
	getAccountsTreeQuery = string(MustAsset("get_accounts_tree.cyp"))
	getTxListQuery = string(MustAsset("get_tx_list.cyp"))
	getTxQuery = string(MustAsset("get_tx.cyp"))
	patchAccountQuery = string(MustAsset("patch_account.cyp"))
	patchTxQuery = string(MustAsset("patch_tx.cyp"))
	postAccountQuery = string(MustAsset("post_account.cyp"))
	postTxQuery = string(MustAsset("post_tx.cyp"))
	putAccountQuery = string(MustAsset("put_account.cyp"))
	putTxQuery = string(MustAsset("put_tx.cyp"))
}

func NewRoutes(neo4jDriver neo4j.Driver) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/routes", routesDoc(r))
		r.Route("/accounts", func(r chi.Router) {
			r.Use(util.CreateNeo4jSessionMiddleware(neo4jDriver))

			r.Get("/", GetAccountTree())
			r.Post("/", CreateAccount())
			r.Route("/{ID}", func(r chi.Router) {
				r.Get("/", GetAccount())
				r.Put("/", UpdateAccount())
				r.Patch("/", PatchAccount())
				r.Delete("/", DeleteAccount())
			})
		})
		r.Route("/transactions", func(r chi.Router) {
			r.Use(util.CreateNeo4jSessionMiddleware(neo4jDriver))

			r.Get("/", GetTransactionList())
			r.Post("/", CreateTransaction())
			r.Route("/{ID}", func(r chi.Router) {
				r.Get("/", GetTransaction())
				r.Put("/", UpdateTransaction())
				r.Patch("/", PatchTransaction())
				r.Delete("/", DeleteTransaction())
			})
		})
	}
}
