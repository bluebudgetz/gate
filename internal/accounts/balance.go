package accounts

import (
	"encoding/json"
	"github.com/bluebudgetz/gate/internal/util"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func (module *Module) getAccountsBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var aggDoc bson.A
	if err := json.Unmarshal(MustAsset("mongodb-balance-query.json"), &aggDoc); err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed loading MongoDB query")
		return
	}

	cur, err := module.mongo.Database("bluebudgetz").Collection("accounts").Aggregate(ctx, aggDoc)
	if err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed fetching accounts from MongoDB")
		return
	}
	defer cur.Close(ctx)

	list := make([]AccountDocument, 0)
	for cur.Next(ctx) {
		var doc AccountDocument
		if err := cur.Decode(&doc); err != nil {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(err).Msg("Failed decoding accounts from MongoDB")
			return
		}
		list = append(list, doc)
	}
	if err := cur.Err(); err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed decoding accounts from MongoDB")
		return
	}

	// Write response back to client
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if module.config.Environment != util.EnvProduction {
		encoder.SetIndent("", "  ")
	}
	w.WriteHeader(http.StatusOK)
	if err := encoder.Encode(list); err != nil {
		log.Error().Err(err).Msg("Failed encoding accounts to client")
		return
	}
}
