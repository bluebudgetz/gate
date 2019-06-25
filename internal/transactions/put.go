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

func (module *Module) putTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body map[string]interface{}

	// Get ID from URL path
	ID := mongoutil.ObjectIdFromNillableString(chi.URLParam(r, "id"))
	if ID == nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// Decode updated transaction from body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, "Invalid JSON payload.", http.StatusBadRequest)
		log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed decoding transaction JSON from request")
		return
	}

	// Validate JSON payload
	var validationErrors []jsonschema.ValError = nil
	module.jsonSchemaRegistry.V1.Transactions.PUT.Validate("/", body, &validationErrors)
	if len(validationErrors) > 0 {
		http.Error(w, fmt.Sprintf("%s (%s)", validationErrors[0].Message, validationErrors[0].PropertyPath), http.StatusBadRequest)
		return
	}

	// Update
	upsert := false
	returnDocument := options.After
	result := module.mongo.Database("bluebudgetz").Collection("transactions").FindOneAndUpdate(ctx, &bson.M{"_id": ID}, &bson.M{
		"$currentDate": &bson.M{
			"updatedOn": true,
		},
		"$set": &bson.M{
			"issuedOn":        body["issuedOn"],
			"origin":          body["origin"],
			"sourceAccountId": mongoutil.ObjectIdFromNillableString(body["sourceAccountId"]),
			"targetAccountId": mongoutil.ObjectIdFromNillableString(body["targetAccountId"]),
			"amount":          body["amount"],
			"comments":        body["comments"],
		},
	}, &options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &returnDocument})
	var updated TransactionDocument
	if result.Err() != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(result.Err()).Str("id", ID.Hex()).Msg("Failed updating transaction in MongoDB")
		return
	} else if err := result.Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Transaction could not be found.", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed fetching updated transaction from MongoDB")
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
		log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed encoding transaction to client")
		return
	}
}