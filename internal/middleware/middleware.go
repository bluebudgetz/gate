package middleware

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"net/http"
)

const (
	dbKey = "db"
)

func GetDB(ctx context.Context) *sql.DB {
	raw, ok := ctx.Value(dbKey).(*sql.DB)
	if ok {
		return raw
	} else {
		panic(errors.New("could not find sql.DB in context"))
	}
}

func EnrichContext(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, dbKey, db)
}

func NewPersistenceMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = EnrichContext(r.Context(), db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
