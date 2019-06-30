package auth

import (
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func (module *Module) createToken(w http.ResponseWriter, r *http.Request) {

	// Credentials are sent as a form; Parse the request form payload
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid payload.", http.StatusBadRequest)
		return
	}

	// Extract username
	username := r.PostFormValue("username")
	if username == "" {
		http.Error(w, "Username is required.", http.StatusBadRequest)
		return
	}

	// Extract password
	password := r.PostFormValue("password")
	if password == "" {
		http.Error(w, "Password is required.", http.StatusBadRequest)
		return
	} else {
		password = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	}

	// Fetch user
	var userDocument bson.M
	result := module.mongo.Database("bluebudgetz").Collection("users").FindOne(r.Context(), bson.M{"_id": username})
	if result.Err() != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(result.Err()).Str("username", username).Msg("Failed fetching user from MongoDB")
		return
	} else if err := result.Decode(&userDocument); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Bad credentials.", http.StatusUnauthorized)
			return
		} else {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(result.Err()).Str("username", username).Msg("Failed fetching user from MongoDB")
			return
		}
	}

	// Compare password hashes
	storedPasswordHash := userDocument["password"]
	if storedPasswordHash != password {
		http.Error(w, "Bad credentials.", http.StatusUnauthorized)
		return
	}

	// Create JWT token claims
	tokenIssuedAt := time.Now()
	tokenNotBefore := time.Now().Add(-10 * time.Minute)
	tokenExpiresAt := time.Now().Add(24 * time.Hour * 14)
	claims := &GateClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "gate",
			Subject:   username,
			ExpiresAt: tokenExpiresAt.Unix(),
			NotBefore: tokenNotBefore.Unix(),
			IssuedAt:  tokenIssuedAt.Unix(),
		},
	}

	// Save new token in MongoDb
	tokenDocument := bson.M{
		"issuer":    claims.Issuer,
		"subject":   claims.Subject,
		"expiresAt": tokenExpiresAt,
		"notBefore": tokenNotBefore,
		"createdOn": tokenIssuedAt,
		"revokedOn": nil,
	}
	if insertResult, err := module.mongo.Database("bluebudgetz").Collection("tokens").InsertOne(r.Context(), tokenDocument); err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).
			Str("username", username).Interface("doc", tokenDocument).
			Msg("Failed saving token to MongoDB")
		return
	} else {
		claims.Id = insertResult.InsertedID.(primitive.ObjectID).Hex()
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(module.jwtKey))
	if err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).
			Str("username", username).Interface("claims", claims).
			Msg("Failed signing JWT key.")
		return
	}

	// Write back token to client both in body and in cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session", // TODO: make session cookie name configurable
		Value:   tokenString,
		Expires: time.Now().Add(module.config.HTTP.JWT.ExpDuration),
	})
	w.Header().Set("Location", r.RequestURI+"/"+claims.Id)
	w.WriteHeader(http.StatusOK)
}
