package auth

import (
	mongoutil "github.com/bluebudgetz/gate/internal/util/mongo"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func (module *Module) revokeToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentClaims := GetClaims(ctx)
	coll := module.mongo.Database("bluebudgetz").Collection("tokens")

	// Get ID from URL path
	ID := mongoutil.ObjectIdFromNillableString(chi.URLParam(r, "id"))
	if ID == nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	// Read token from database
	// NOTE: we DO NOT check if the token has been revoked, since the middleware does this for us
	var tokenDocument bson.M
	result := coll.FindOne(ctx, bson.M{"_id": ID})
	if result.Err() != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(result.Err()).Str("id", ID.Hex()).Msg("Failed fetching token from MongoDB")
		return
	} else if err := result.Decode(&tokenDocument); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed fetching token from MongoDB")
			return
		}
	}

	// Compare currently authenticate subject owns the token that belongs to this ID
	tokenSubject := tokenDocument["subject"]
	if tokenSubject != currentClaims.Subject {
		http.Error(w, "Forbidden.", http.StatusForbidden)
		return
	}

	// Save new token in MongoDb
	patchDocument := bson.M{
		"$set": bson.M{
			"revokedOn": time.Now(),
		},
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": ID}, patchDocument); err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).
			Str("id", ID.Hex()).Interface("doc", patchDocument).
			Msg("Failed updating token in MongoDB")
		return
	}

	// Delete the session cookie by setting its expiration time to now
	http.SetCookie(w, &http.Cookie{
		Name:    "session", // TODO: make session cookie name configurable
		Value:   "",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
}
