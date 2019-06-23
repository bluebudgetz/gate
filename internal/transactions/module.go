package transactions

import (
	"github.com/bluebudgetz/gate/internal/schema"
	"github.com/bluebudgetz/gate/internal/util"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TransactionDocument struct {
	ID        *primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CreatedOn *time.Time          `bson:"createdOn" json:"createdOn,omitempty"`
	UpdatedOn *time.Time          `bson:"updatedOn" json:"updatedOn,omitempty"`
	IssuedOn  *time.Time          `bson:"issuedOn" json:"issuedOn,omitempty"`
	Origin    *string             `bson:"origin" json:"origin,omitempty"`
	SourceID  *primitive.ObjectID `bson:"source" json:"sourceAccountId,omitempty"`
	TargetID  *primitive.ObjectID `bson:"target" json:"targetAccountId,omitempty"`
	Amount    *float64            `bson:"amount" json:"amount,omitempty"`
	Comments  *string             `bson:"comments" json:"comments,omitempty"`
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
	router.Get("/", module.getTransactions)
	router.Post("/", module.addTransaction)
	router.Delete("/{id}", module.deleteTransaction)
	router.Get("/{id}", module.getTransaction)
	router.Patch("/{id}", module.patchTransaction)
	router.Put("/{id}", module.putTransaction)
}
