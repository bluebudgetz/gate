package auth

import (
	mongoutil "github.com/bluebudgetz/gate/internal/util/mongo"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func (module *Module) refreshToken(w http.ResponseWriter, r *http.Request) {
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

	// Create JWT token claims
	tokenIssuedAt := time.Now()
	tokenNotBefore := time.Now().Add(-10 * time.Minute)
	tokenExpiresAt := time.Now().Add(24 * time.Hour * 14)
	newClaims := &GateClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    currentClaims.Issuer,
			Subject:   currentClaims.Subject,
			ExpiresAt: tokenExpiresAt.Unix(),
			NotBefore: tokenNotBefore.Unix(),
			IssuedAt:  tokenIssuedAt.Unix(),
		},
	}

	// Save new token in MongoDb
	patchDocument := bson.M{
		"$set": bson.M{
			"expiresAt": tokenExpiresAt,
			"notBefore": tokenNotBefore,
			"createdOn": tokenIssuedAt,
		},
	}
	if _, err := coll.UpdateOne(ctx, bson.M{"_id": ID}, patchDocument); err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).
			Str("id", ID.Hex()).Interface("doc", patchDocument).
			Msg("Failed updating token in MongoDB")
		return
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims).SignedString([]byte(module.jwtKey))
	if err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).
			Str("id", ID.Hex()).Interface("claims", newClaims).
			Msg("Failed signing JWT key.")
		return
	}

	// Write back token to client both in body and in cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session", // TODO: make session cookie name configurable
		Value:   tokenString,
		Expires: time.Now().Add(module.config.HTTP.JWT.ExpDuration),
	})
	w.WriteHeader(http.StatusOK)
}
