package accounts

import (
	"database/sql"
	"encoding/json"
	"github.com/bluebudgetz/gate/internal/api/util"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type Accounts struct {
	db   *sql.DB
	port int
}

func New(db *sql.DB, port int) *Accounts {
	return &Accounts{db: db, port: port}
}

func (a *Accounts) getRootAccounts(w http.ResponseWriter, r *http.Request) {
	_sql := `
		WITH RECURSIVE cte (id, child_id, created_on, updated_on, name) AS (
	        SELECT acc.id, acc.id, acc.created_on, acc.updated_on, acc.name
	        FROM bb.accounts AS acc
	        WHERE acc.parent_id IS NULL
	        UNION ALL
	        SELECT h.id, acc.id, h.created_on, h.updated_on, h.name
	        FROM bb.accounts AS acc
	            INNER JOIN cte AS h ON acc.parent_id = h.child_id
	    )
		SELECT
		    cte.id                                                                                          AS id,
       		cte.created_on                                                                                  AS created_on,
       		cte.updated_on                                                                                  AS updated_on,
       		cte.name                                                                                        AS name,
        	COALESCE((SELECT COUNT(*) FROM bb.accounts AS children WHERE children.parent_id = cte.id), 0)   AS child_count,
        	COALESCE(SUM(outgoing_tx.amount), 0)                                                            AS outgoing,
        	COALESCE(SUM(incoming_tx.amount), 0)                                                            AS incoming,
        	COALESCE(SUM(incoming_tx.amount), 0) - COALESCE(SUM(outgoing_tx.amount), 0)                     AS balance
		FROM cte
         	LEFT JOIN bb.transactions AS outgoing_tx ON outgoing_tx.source_account_id = cte.child_id
         	LEFT JOIN bb.transactions AS incoming_tx ON incoming_tx.target_account_id = cte.child_id
		GROUP BY 
			cte.id, 
		    cte.created_on, 
	        cte.updated_on,
	        cte.name 
        ORDER BY cte.name, cte.id `
	rows, err := a.db.QueryContext(r.Context(), _sql)
	if err != nil {
		panic(errors.Wrap(err, "failed fetching root accounts"))
	}
	defer rows.Close()

	accounts := make([]Account, 0, 100)
	for rows.Next() {
		var acc Account
		err := rows.Scan(&acc.ID, &acc.CreatedOn, &acc.UpdatedOn, &acc.Name, &acc.ChildCount, &acc.Outgoing, &acc.Incoming, &acc.Balance)
		if err != nil {
			panic(errors.Wrap(err, "failed scanning account"))
		}
		accounts = append(accounts, acc)
	}
	if err := rows.Err(); err != nil {
		panic(errors.Wrap(err, "failed iterating root account rows"))
	}

	util.Respond(w, r, http.StatusOK, accounts)
}

func (a *Accounts) getChildAccounts(w http.ResponseWriter, r *http.Request) {
	parentID := chi.URLParam(r, "id")
	_sql := `
		WITH RECURSIVE cte (id, child_id, created_on, updated_on, name) AS (
	        SELECT acc.id, acc.id, acc.created_on, acc.updated_on, acc.name
	        FROM bb.accounts AS acc
	        WHERE acc.parent_id = $1
	        UNION ALL
	        SELECT h.id, acc.id, h.created_on, h.updated_on, h.name
	        FROM bb.accounts AS acc
	            INNER JOIN cte AS h ON acc.parent_id = h.child_id
	    )
		SELECT
		    cte.id                                                                                          AS id,
       		cte.created_on                                                                                  AS created_on,
       		cte.updated_on                                                                                  AS updated_on,
       		cte.name                                                                                        AS name,
        	COALESCE((SELECT COUNT(*) FROM bb.accounts AS children WHERE children.parent_id = cte.id), 0)   AS child_count,
        	COALESCE(SUM(outgoing_tx.amount), 0)                                                            AS outgoing,
        	COALESCE(SUM(incoming_tx.amount), 0)                                                            AS incoming,
        	COALESCE(SUM(incoming_tx.amount), 0) - COALESCE(SUM(outgoing_tx.amount), 0)                     AS balance
		FROM cte
         	LEFT JOIN bb.transactions AS outgoing_tx ON outgoing_tx.source_account_id = cte.child_id
         	LEFT JOIN bb.transactions AS incoming_tx ON incoming_tx.target_account_id = cte.child_id
		GROUP BY 
			cte.id, 
		    cte.created_on, 
	        cte.updated_on,
	        cte.name 
        ORDER BY cte.name, cte.id `
	rows, err := a.db.QueryContext(r.Context(), _sql, parentID)
	if err != nil {
		panic(errors.Wrap(err, "failed fetching root accounts"))
	}
	defer rows.Close()

	accounts := make([]Account, 0, 100)
	for rows.Next() {
		var acc Account
		err := rows.Scan(&acc.ID, &acc.CreatedOn, &acc.UpdatedOn, &acc.Name, &acc.ChildCount, &acc.Outgoing, &acc.Incoming, &acc.Balance)
		if err != nil {
			panic(errors.Wrap(err, "failed scanning account"))
		}
		accounts = append(accounts, acc)
	}
	if err := rows.Err(); err != nil {
		panic(errors.Wrap(err, "failed iterating root account rows"))
	}

	util.Respond(w, r, http.StatusOK, accounts)
}

func (a *Accounts) getAccount(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	_sql := `
		WITH RECURSIVE cte (id, child_id, created_on, updated_on, name) AS (
	        SELECT acc.id, acc.id, acc.created_on, acc.updated_on, acc.name
	        FROM bb.accounts AS acc
	        WHERE acc.id = $1
	        UNION ALL
	        SELECT h.id, acc.id, h.created_on, h.updated_on, h.name
	        FROM bb.accounts AS acc
	            INNER JOIN cte AS h ON acc.parent_id = h.child_id)
		SELECT
		    cte.id                                                                                          AS id,
       		cte.created_on                                                                                  AS created_on,
       		cte.updated_on                                                                                  AS updated_on,
       		cte.name                                                                                        AS name,
        	COALESCE((SELECT COUNT(*) FROM bb.accounts AS children WHERE children.parent_id = cte.id), 0)   AS child_count,
        	COALESCE(SUM(outgoing_tx.amount), 0)                                                            AS outgoing,
        	COALESCE(SUM(incoming_tx.amount), 0)                                                            AS incoming,
        	COALESCE(SUM(incoming_tx.amount), 0) - COALESCE(SUM(outgoing_tx.amount), 0)                     AS balance
		FROM cte
         	LEFT JOIN bb.transactions AS outgoing_tx ON outgoing_tx.source_account_id = cte.child_id
         	LEFT JOIN bb.transactions AS incoming_tx ON incoming_tx.target_account_id = cte.child_id
		GROUP BY 
			cte.id, 
		    cte.created_on, 
	        cte.updated_on,
	        cte.name`
	var acc Account
	err := a.db.QueryRowContext(r.Context(), _sql, ID).Scan(&acc.ID, &acc.CreatedOn, &acc.UpdatedOn, &acc.Name, &acc.ChildCount, &acc.Outgoing, &acc.Incoming, &acc.Balance)
	if err != nil {
		panic(errors.Wrap(err, "failed scanning account"))
	}

	util.Respond(w, r, http.StatusOK, acc)
}

func (a *Accounts) createAccount(w http.ResponseWriter, r *http.Request) {
	account := struct {
		Name     string `json:"name"`
		ParentID *int   `json:"parentId"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		util.Respond(w, r, http.StatusBadRequest, nil)
		return
	}

	_sql := `
		INSERT INTO bb.accounts (name, parent_id) 
		VALUES ($1, $2) 
		RETURNING id
	`

	var id int
	err = a.db.QueryRowContext(r.Context(), _sql, account.Name, account.ParentID).Scan(&id)
	if err != nil {
		panic(errors.Wrapf(err, "failed creating account '%s'", account.Name))
	}

	w.Header().Set("Location", r.RequestURI+"/"+strconv.FormatInt(int64(id), 10))
	util.Respond(w, r, http.StatusCreated, nil)
}

func (a *Accounts) putAccount(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")

	account := struct {
		Name     string `json:"name"`
		ParentID *int   `json:"parentId"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		util.Respond(w, r, http.StatusBadRequest, nil)
		return
	}

	_sql := `
		UPDATE bb.accounts 
		SET name = $2, parent_id = $3
		WHERE id = $1
	`
	result, err := a.db.ExecContext(r.Context(), _sql, ID, account.Name, account.ParentID)
	if err != nil {
		panic(errors.Wrapf(err, "failed updating account '%s'", account.Name))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(errors.Wrapf(err, "failed fetching number of affected rows"))
	} else if rowsAffected <= 0 {
		panic(errors.Wrapf(err, "could not find any account with ID '%d'", ID))
	}

	util.Respond(w, r, http.StatusNoContent, nil)
}

func (a *Accounts) patchAccount(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")

	_selectSql := `SELECT name, parent_id FROM bb.accounts WHERE id = $1`
	var name string
	var parentId *int
	if err := a.db.QueryRowContext(r.Context(), _selectSql, ID).Scan(&name, &parentId); err != nil {
		panic(errors.Wrapf(err, "failed fetching current values of account '%d'", ID))
	}

	account := struct {
		Name     *string `json:"name"`
		ParentID *int    `json:"parentId"`
	}{Name: &name, ParentID: parentId}
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		util.Respond(w, r, http.StatusBadRequest, nil)
		return
	}

	_sql := `UPDATE bb.accounts SET name = $2, parent_id = $3 WHERE id = $1`
	result, err := a.db.ExecContext(r.Context(), _sql, ID, account.Name, account.ParentID)
	if err != nil {
		panic(errors.Wrapf(err, "failed patching account '%s'", account.Name))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(errors.Wrapf(err, "failed fetching number of affected rows"))
	} else if rowsAffected <= 0 {
		panic(errors.Wrapf(err, "could not find any account with ID '%d'", ID))
	}

	util.Respond(w, r, http.StatusNoContent, nil)
}

func (a *Accounts) RoutesV1() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", a.getRootAccounts)
	router.Post("/", a.createAccount)
	router.Get("/{id}", a.getAccount)
	router.Put("/{id}", a.putAccount)
	router.Patch("/{id}", a.patchAccount)
	router.Get("/{id}/children", a.getChildAccounts)
	return router
}
