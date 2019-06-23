package transactions

import (
	mongoutil "github.com/bluebudgetz/gate/internal/util/mongo"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func (module *Module) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get ID from URL path
	ID := mongoutil.ObjectIdFromNillableString(chi.URLParam(r, "id"))
	if ID == nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	result, err := module.mongo.Database("bluebudgetz").Collection("transactions").DeleteOne(ctx, &bson.M{"_id": ID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Transaction could not be found.", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed deleting transaction from MongoDB")
			return
		}
	} else if result.DeletedCount < 1 {
		http.Error(w, "Transaction could not be found.", http.StatusNotFound)
		return
	} else if result.DeletedCount > 1 {
		log.Warn().
			Str("id", ID.Hex()).
			Int64("count", result.DeletedCount).
			Msgf("%d MongoDB transactions deleted for a single ID")
	}

	w.WriteHeader(http.StatusNoContent)
}
