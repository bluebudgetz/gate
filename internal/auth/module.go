package auth

import (
	"github.com/bluebudgetz/gate/internal/schema"
	"github.com/bluebudgetz/gate/internal/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

type GateClaims struct {
	jwt.StandardClaims
}

type Module struct {
	config             util.Config
	jsonSchemaRegistry *schema.Registry
	mongo              *mongo.Client
	jwtKey             string
}

func New(jwtKey string, config util.Config, jsonSchemaRegistry *schema.Registry, mongo *mongo.Client) *Module {
	return &Module{config, jsonSchemaRegistry, mongo, jwtKey}
}

func (module *Module) RoutesV1(router chi.Router) {
	// TODO: use separate APIs for token management and session (cookie) management
	router.Post("/tokens", module.createToken)
	router.Patch("/tokens/{id}", module.refreshToken)
	router.Delete("/tokens/{id}", module.revokeToken)
}
