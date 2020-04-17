package services

import (
	"context"
	"net/http"

	"github.com/golangly/errors"
	"github.com/golangly/log"
	"github.com/golangly/webutil"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func NewNeo4jDriver() (neo4j.Driver, func(), error) {
	neo4jDriver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.NoAuth())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed creating Neo4j client")
	}
	return neo4jDriver, func() {
		if err := neo4jDriver.Close(); err != nil {
			log.WithErr(err).Error("Failed disposing Neo4j driver")
		}
	}, nil
}

// Key to use when setting the Neo4j session.
type ctxKeyNeo4jSession int

// The key that holds the Neo4j session in a request context.
const neo4jSessionKey ctxKeyNeo4jSession = 0

func GetNeo4jSession(ctx context.Context) neo4j.Session {
	if ctx == nil {
		return nil
	}
	if session, ok := ctx.Value(neo4jSessionKey).(neo4j.Session); ok {
		return session
	}
	return nil
}

func WithNeo4jSession(ctx context.Context, session neo4j.Session) context.Context {
	return context.WithValue(ctx, neo4jSessionKey, session)
}

func CreateNeo4jSession(neo4jDriver neo4j.Driver, mode neo4j.AccessMode) (neo4j.Session, error) {
	session, err := neo4jDriver.Session(mode)
	if err != nil {
		return nil, errors.Wrap(err, "create session error").AddTag("mode", mode)
	}
	return session, nil
}

func CreateNeo4jSessionMiddleware(neo4jDriver neo4j.Driver) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mode := neo4j.AccessModeRead
			switch r.Method {
			case http.MethodGet, http.MethodHead, http.MethodConnect, http.MethodOptions, http.MethodTrace:
				mode = neo4j.AccessModeRead
			case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
				mode = neo4j.AccessModeWrite
			}

			session, err := CreateNeo4jSession(neo4jDriver, mode)
			if err != nil {
				webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, errors.Wrap(err, "create session error").AddTag("mode", mode))
				return
			}
			defer session.Close()

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), neo4jSessionKey, session)))
		})
	}
}
