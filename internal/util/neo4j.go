package util

import (
	"context"
	"net/http"

	"github.com/golangly/errors"
	"github.com/golangly/webutil"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

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
			session, err := neo4jDriver.Session(mode)
			if err != nil {
				webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, errors.Wrap(err, "create session error").AddTag("mode", mode))
				return
			}
			defer session.Close()

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), neo4jSessionKey, session)))
		})
	}
}

type Neo4jSessionMock struct{}

func (m *Neo4jSessionMock) LastBookmark() string { panic("LastBookmark() not implemented") }

//noinspection GoUnusedParameter
func (m *Neo4jSessionMock) BeginTransaction(configurers ...func(*neo4j.TransactionConfig)) (neo4j.Transaction, error) {
	panic("BeginTransaction(...func(*neo4j.TransactionConfig)) not implemented")
}

//noinspection GoUnusedParameter
func (m *Neo4jSessionMock) ReadTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	panic("ReadTransaction(neo4j.TransactionWork, ...func(*neo4j.TransactionConfig)) not implemented")
}

//noinspection GoUnusedParameter
func (m *Neo4jSessionMock) WriteTransaction(work neo4j.TransactionWork, configurers ...func(*neo4j.TransactionConfig)) (interface{}, error) {
	panic("WriteTransaction(neo4j.TransactionWork, ...func(*neo4j.TransactionConfig)) not implemented")
}

//noinspection GoUnusedParameter
func (m *Neo4jSessionMock) Run(cypher string, params map[string]interface{}, configurers ...func(*neo4j.TransactionConfig)) (neo4j.Result, error) {
	panic("Run(string, map[string]interface{}, ...func(*neo4j.TransactionConfig)) not implemented")
}

//noinspection GoUnusedParameter
func (m *Neo4jSessionMock) Close() error {
	panic("Close() not implement")
}
