package auth

import (
	"context"
	mongoutil "github.com/bluebudgetz/gate/internal/util/mongo"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

const claimsRequestKey = "___token"

func GetClaims(ctx context.Context) *GateClaims {
	return ctx.Value(claimsRequestKey).(*GateClaims)
}

func (module *Module) Authenticate(ignoredPaths ...string) func(next http.Handler) http.Handler {
	keyFunc := func(token *jwt.Token) (interface{}, error) { return []byte(module.jwtKey), nil }
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Ignore URLs that should not be authenticated (e.g. POST /v1/auth/tokens which is a login...)
			for _, path := range ignoredPaths {
				if path == r.URL.Path {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Get session cookie and read token from it
			cookie, err := r.Cookie("session") // TODO: make session cookie name configurable
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, "Not authenticated.", http.StatusUnauthorized)
				} else {
					http.Error(w, "Could not read authentication.", http.StatusBadRequest)
				}
				return
			}

			// Read and validate that the token is valid
			claims := &GateClaims{}
			if tkn, err := jwt.ParseWithClaims(cookie.Value, claims, keyFunc); err != nil {
				if err == jwt.ErrSignatureInvalid {
					http.Error(w, "Bad credentials.", http.StatusUnauthorized)
				} else {
					http.Error(w, "Bad request.", http.StatusBadRequest)
				}
				return
			} else if !tkn.Valid {
				http.Error(w, "Bad credentials.", http.StatusUnauthorized)
				return
			}

			// Fetch token and ensure it can still be used (exists in our database; not revoked yet; etc)
			var tokenDocument bson.M
			coll := module.mongo.Database("bluebudgetz").Collection("tokens")
			result := coll.FindOne(ctx, bson.M{"_id": mongoutil.ObjectIdFromNillableString(claims.Id)})
			if result.Err() != nil {
				http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
				log.Error().Err(result.Err()).Str("id", claims.Id).Msg("Failed fetching token from MongoDB")
				return
			} else if err := result.Decode(&tokenDocument); err != nil {
				if err == mongo.ErrNoDocuments {
					http.Error(w, "Forbidden.", http.StatusForbidden)
					log.Warn().Err(err).Str("id", claims.Id).Msg("Unknown but valid JWT token encountered!")
					return
				} else {
					http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
					log.Error().Err(err).Str("id", claims.Id).Msg("Failed fetching token from MongoDB")
					return
				}
			} else if _, ok := tokenDocument["revokedOn"].(time.Time); ok {
				// REMEMBER that if we remove this check, we MUST add it to the token refresh handler!
				// Otherwise, users will be able to refresh other users' tokens, and thus impersonate them!
				http.Error(w, "Forbidden.", http.StatusForbidden)
				log.Warn().Err(err).Str("id", claims.Id).Msg("Revoked JWT token used")
				return
			}

			// All is well, continue to next handler
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, claimsRequestKey, claims)))
		})
	}
}
