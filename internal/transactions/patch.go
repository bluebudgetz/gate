package transactions

import (
	"encoding/json"
	"fmt"
	"github.com/bluebudgetz/gate/internal/util"
	mongoutil "github.com/bluebudgetz/gate/internal/util/mongo"
	"github.com/go-chi/chi"
	"github.com/qri-io/jsonschema"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

func (module *Module) patchTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get ID from URL path
	ID := mongoutil.ObjectIdFromNillableString(chi.URLParam(r, "id"))
	if ID == nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// Decode patch spec from body
	var body map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "Invalid JSON payload.", http.StatusBadRequest)
		log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed decoding transaction JSON from request")
		return
	}

	// Validate JSON payload
	var validationErrors []jsonschema.ValError = nil
	module.jsonSchemaRegistry.V1.Transactions.PATCH.Validate("/", body, &validationErrors)
	if len(validationErrors) > 0 {
		http.Error(w, fmt.Sprintf("%s (%s)", validationErrors[0].Message, validationErrors[0].PropertyPath), http.StatusBadRequest)
		return
	}

	// Build patch query
	patch := bson.M{}
	if issuedOn, ok := body["issuedOn"]; ok {
		patch["issuedOn"] = issuedOn
	}
	if origin, ok := body["origin"]; ok {
		patch["origin"] = origin
	}
	if sourceAccountId := mongoutil.ObjectIdFromNillableString(body["sourceAccountId"]); sourceAccountId != nil {
		patch["sourceAccountId"] = sourceAccountId
	}
	if targetAccountId := mongoutil.ObjectIdFromNillableString(body["targetAccountId"]); targetAccountId != nil {
		patch["targetAccountId"] = targetAccountId
	}
	if amount, ok := body["amount"]; ok {
		patch["amount"] = amount
	}
	if comments, ok := body["comments"]; ok {
		patch["comments"] = comments
	}

	// Update
	upsert := false
	returnDocument := options.After
	result := module.mongo.Database("bluebudgetz").Collection("transactions").FindOneAndUpdate(ctx, &bson.M{"_id": ID}, &bson.M{
		"$currentDate": &bson.M{
			"updatedOn": true,
		},
		"$set": &patch,
	}, &options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &returnDocument})
	var updated TransactionDocument
	if result.Err() != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(result.Err()).Str("id", ID.Hex()).Msg("Failed patching transaction in MongoDB")
		return
	} else if err := result.Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Transaction could not be found.", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed fetching patched transaction from MongoDB")
			return
		}
	}

	// Write response back to client
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if module.config.Environment != util.EnvProduction {
		encoder.SetIndent("", "  ")
	}
	w.WriteHeader(http.StatusOK)
	if err := encoder.Encode(updated); err != nil {
		log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed encoding transactio to client")
		return
	}
}
