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

type Module struct {
	config             util.Config
	jsonSchemaRegistry *schema.Registry
	mongo              *mongo.Client
}

func New(config util.Config, jsonSchemaRegistry *schema.Registry, mongo *mongo.Client) *Module {
	return &Module{config, jsonSchemaRegistry, mongo}
}

func (module *Module) RoutesV1(router chi.Router) {
	// Root
	router.Get("/", module.getAccountsList)
	router.Post("/", module.addAccount)
	router.Get("/tree", module.getAccountsTree)

	// Single account
	router.Delete("/{id}", module.deleteAccount)
	router.Get("/{id}", module.getAccount)
	router.Patch("/{id}", module.patchAccount)
	router.Put("/{id}", module.putAccount)
}
