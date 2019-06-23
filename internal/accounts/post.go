package accounts

import (
	"encoding/json"
	"fmt"
	mongoutil "github.com/bluebudgetz/gate/internal/util/mongo"
	"github.com/qri-io/jsonschema"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func (acc *Accounts) addAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode new account from body
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON payload.", http.StatusBadRequest)
		log.Error().Err(err).Msg("Failed decoding account JSON from request")
		return
	}

	// Validate JSON payload
	var validationErrors []jsonschema.ValError = nil
	acc.jsonSchemaRegistry.V1.Accounts.POST.Validate("/", body, &validationErrors)
	if len(validationErrors) > 0 {
		http.Error(w, fmt.Sprintf("%s (%s)", validationErrors[0].Message, validationErrors[0].PropertyPath), http.StatusBadRequest)
		return
	}

	// Create account
	post := bson.M{
		"createdOn": time.Now(),
		"updatedOn": nil,
		"name":      body["name"],
		"parentId":  mongoutil.ObjectIdFromNillableString(body["parentId"]),
	}
	if log.Debug().Enabled() {
		log.Debug().Interface("body", body).Interface("account", post).Msg("Creating account")
	}
	result, err := acc.mongo.Database("bluebudgetz").Collection("accounts").InsertOne(ctx, &post)
	if err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed creating account")
		return
	}

	w.Header().Set("Location", r.RequestURI+"/"+result.InsertedID.(primitive.ObjectID).Hex())
	w.WriteHeader(http.StatusCreated)
}
