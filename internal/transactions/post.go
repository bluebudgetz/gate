package transactions

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

func (module *Module) addTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode new transaction from body
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON payload.", http.StatusBadRequest)
		log.Error().Err(err).Msg("Failed decoding transaction JSON from request")
		return
	}

	// Validate JSON payload
	var validationErrors []jsonschema.ValError = nil
	module.jsonSchemaRegistry.V1.Transactions.POST.Validate("/", body, &validationErrors)
	if len(validationErrors) > 0 {
		http.Error(w, fmt.Sprintf("%s (%s)", validationErrors[0].Message, validationErrors[0].PropertyPath), http.StatusBadRequest)
		return
	}

	// Create transaction
	post := bson.M{
		"createdOn":       time.Now(),
		"updatedOn":       nil,
		"issuedOn":        body["issuedOn"],
		"origin":          body["origin"],
		"sourceAccountId": mongoutil.ObjectIdFromNillableString(body["sourceAccountId"]),
		"targetAccountId": mongoutil.ObjectIdFromNillableString(body["targetAccountId"]),
		"amount":          body["amount"],
		"comments":        body["comments"],
	}
	if log.Debug().Enabled() {
		log.Debug().Interface("body", body).Interface("transaction", post).Msg("Creating transaction")
	}
	result, err := module.mongo.Database("bluebudgetz").Collection("transactions").InsertOne(ctx, &post)
	if err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed creating transaction")
		return
	}

	w.Header().Set("Location", r.RequestURI+"/"+result.InsertedID.(primitive.ObjectID).Hex())
	w.WriteHeader(http.StatusCreated)
}
