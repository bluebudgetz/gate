package internal

import (
	"database/sql"
	"github.com/bluebudgetz/gate/internal/model"
	"github.com/pkg/errors"
	"net/http"
)

func NewPersistenceMiddleware(db *sql.DB) func(http.Handler) http.Handler {
	accountsDataManager, err := model.NewAccountsDataManager(db)
	if err != nil {
		panic(errors.Wrap(err, "could not create accounts data manager"))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(accountsDataManager.PutInContext(r.Context())))
		})
	}
}
