package accounts

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

func (acc *Accounts) patchAccount(w http.ResponseWriter, r *http.Request) {
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
		log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed decoding account JSON from request")
		return
	}

	// Validate JSON payload
	var validationErrors []jsonschema.ValError = nil
	acc.jsonSchemaRegistry.V1.Accounts.PATCH.Validate("/", body, &validationErrors)
	if len(validationErrors) > 0 {
		http.Error(w, fmt.Sprintf("%s (%s)", validationErrors[0].Message, validationErrors[0].PropertyPath), http.StatusBadRequest)
		return
	}

	// Build patch query
	patch := bson.M{}
	if name, ok := body["name"]; ok {
		patch["name"] = name
	}
	if parentId := mongoutil.ObjectIdFromNillableString(body["parentId"]); parentId != nil {
		patch["parentId"] = parentId
	}

	// Update
	upsert := false
	returnDocument := options.After
	result := acc.mongo.Database("bluebudgetz").Collection("accounts").FindOneAndUpdate(ctx, &bson.M{"_id": ID}, &bson.M{
		"$currentDate": &bson.M{
			"updatedOn": true,
		},
		"$set": &patch,
	}, &options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &returnDocument})
	var updated AccountDocument
	if result.Err() != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(result.Err()).Str("id", ID.Hex()).Msg("Failed patching account in MongoDB")
		return
	} else if err := result.Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Account could not be found.", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed fetching patched account from MongoDB")
			return
		}
	}

	// Write response back to client
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if acc.config.Environment != util.EnvProduction {
		encoder.SetIndent("", "  ")
	}
	w.WriteHeader(http.StatusOK)
	if err := encoder.Encode(updated); err != nil {
		log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed encoding account to client")
		return
	}
}
