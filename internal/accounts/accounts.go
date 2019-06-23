package accounts

import (
	"github.com/bluebudgetz/gate/internal/schema"
	"github.com/bluebudgetz/gate/internal/util"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type AccountDocument struct {
	ID        *primitive.ObjectID `bson:"_id"`
	CreatedOn *time.Time          `bson:"createdOn"`
	UpdatedOn *time.Time          `bson:"updatedOn"`
	Name      *string             `bson:"name"`
	ParentID  *primitive.ObjectID `bson:"parentId"`
}

type Accounts struct {
	config             util.Config
	jsonSchemaRegistry *schema.Registry
	mongo              *mongo.Client
}

func New(config util.Config, jsonSchemaRegistry *schema.Registry, mongo *mongo.Client) *Accounts {
	return &Accounts{config, jsonSchemaRegistry, mongo}
}

func (acc *Accounts) RoutesV1(router chi.Router) {
	// Root
	router.Get("/", acc.getAccountsList)
	router.Post("/", acc.addAccount)
	router.Get("/tree", acc.getAccountsTree)

	// Single account
	router.Delete("/{id}", acc.deleteAccount)
	router.Get("/{id}", acc.getAccount)
	router.Patch("/{id}", acc.patchAccount)
	router.Put("/{id}", acc.putAccount)
}
