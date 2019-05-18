package transactions

import (
	"database/sql"
	"encoding/json"
	"github.com/bluebudgetz/gate/internal/api/util"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type Transactions struct {
	db *sql.DB
}

func New(db *sql.DB) *Transactions {
	return &Transactions{db}
}

func (t *Transactions) getTransactions(w http.ResponseWriter, r *http.Request) {
	_sql := `
		SELECT
	        tx.id,
		    tx.created_on, 
		    tx.updated_on, 
		    tx.origin, 
		    tx.source_account_id,
		    tx.target_account_id,
		    tx.amount,
		    tx.comments
		FROM
	        bb.transactions AS tx 
        ORDER BY tx.created_on DESC, tx.id `
	rows, err := t.db.QueryContext(r.Context(), _sql)
	if err != nil {
		panic(errors.Wrap(err, "failed fetching transactions"))
	}
	defer rows.Close()

	transactions := make([]Transaction, 0, 100)
	for rows.Next() {
		var tx Transaction
		err := rows.Scan(&tx.ID, &tx.CreatedOn, &tx.UpdatedOn, &tx.Origin, &tx.SourceID, &tx.TargetID, &tx.Amount, &tx.Comments)
		if err != nil {
			panic(errors.Wrap(err, "failed scanning transaction"))
		}
		transactions = append(transactions, tx)
	}
	if err := rows.Err(); err != nil {
		panic(errors.Wrap(err, "failed iterating transaction rows"))
	}
	util.Respond(w, r, http.StatusOK, transactions)
}

func (t *Transactions) getTransaction(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	_sql := `
		SELECT
	        tx.id,
		    tx.created_on, 
		    tx.updated_on, 
		    tx.origin, 
		    tx.source_account_id,
		    tx.target_account_id,
		    tx.amount,
		    tx.comments
		FROM
	        bb.transactions AS tx 
		WHERE
	        tx.id = $1 
        ORDER BY tx.created_on DESC, tx.id `
	var tx Transaction
	err := t.db.QueryRowContext(r.Context(), _sql, ID).Scan(&tx.ID, &tx.CreatedOn, &tx.UpdatedOn, &tx.Origin, &tx.SourceID, &tx.TargetID, &tx.Amount, &tx.Comments)
	if err != nil {
		panic(errors.Wrap(err, "failed scanning transaction"))
	}
	util.Respond(w, r, http.StatusOK, tx)
}

func (t *Transactions) createTransaction(w http.ResponseWriter, r *http.Request) {
	tx := struct {
		Origin   string  `json:"origin"`
		SourceID int     `json:"sourceAccountId"`
		TargetID int     `json:"targetAccountId"`
		Amount   float64 `json:"amount"`
		Comments *string `json:"comments"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		util.Respond(w, r, http.StatusBadRequest, nil)
		return
	}

	_sql := "INSERT INTO bb.transactions (origin, source_account_id, target_account_id, amount, comments) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	var id int
	err = t.db.QueryRowContext(r.Context(), _sql, tx.Origin, tx.SourceID, tx.TargetID, tx.Amount, tx.Comments).Scan(&id)
	if err != nil {
		panic(errors.Wrapf(err, "failed creating transaction"))
	}

	w.Header().Set("Location", r.RequestURI+"/"+strconv.FormatInt(int64(id), 10))
	util.Respond(w, r, http.StatusCreated, nil)
}

func (t *Transactions) RoutesV1() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", t.getTransactions)
	router.Post("/", t.createTransaction)
	router.Get("/{id}", t.getTransaction)
	return router
}
