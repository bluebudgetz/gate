package transactions

import (
	"encoding/json"
	"github.com/bluebudgetz/gate/internal/util"
	mongoutil "github.com/bluebudgetz/gate/internal/util/mongo"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func (module *Module) getTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cur, err := module.mongo.Database("bluebudgetz").Collection("transactions").Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed fetching transactions from MongoDB")
		return
	}
	defer cur.Close(ctx)

	list := make([]TransactionDocument, 0)
	for cur.Next(ctx) {
		var doc TransactionDocument
		if err := cur.Decode(&doc); err != nil {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(err).Msg("Failed decoding transactions from MongoDB")
			return
		}
		list = append(list, doc)
	}
	if err := cur.Err(); err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed decoding transactions from MongoDB")
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
		log.Error().Err(err).Msg("Failed encoding transactions to client")
		return
	}
}

func (module *Module) getTransaction(w http.ResponseWriter, r *http.Request) {
	ID := mongoutil.ObjectIdFromNillableString(chi.URLParam(r, "id"))
	if ID == nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	result := module.mongo.Database("bluebudgetz").Collection("transactions").FindOne(r.Context(), bson.M{"_id": ID})
	if result.Err() != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(result.Err()).Str("id", ID.Hex()).Msg("Failed fetching transaction from MongoDB")
		return
	}

	var account TransactionDocument
	if err := result.Decode(&account); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Transaction could not be found.", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed fetching transaction from MongoDB")
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
	if err := encoder.Encode(account); err != nil {
		log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed encoding transaction to client")
		return
	}
}
