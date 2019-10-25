package transactions

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/bluebudgetz/gate/internal/infra/bind"
	"github.com/bluebudgetz/gate/internal/infra/render"
	"github.com/bluebudgetz/gate/internal/util"
)

type (
	GetTransactionResponse struct {
		Tx Transaction
	}

	PostTransactionRequest struct {
		Body struct {
			IssuedOn        time.Time `json:"issuedOn" yaml:"issuedOn"`
			Origin          string    `json:"origin" yaml:"origin"`
			SourceAccountID string    `json:"sourceAccountId" yaml:"sourceAccountId"`
			TargetAccountID string    `json:"targetAccountId" yaml:"targetAccountId"`
			Amount          float64   `json:"amount" yaml:"amount"`
			Comment         string    `json:"comment" yaml:"comment"`
		} `body:""`
	}

	PostTransactionResponse struct {
		Tx Transaction
	}

	PutTransactionRequest struct {
		Body struct {
			IssuedOn        time.Time `json:"issuedOn" yaml:"issuedOn"`
			Origin          string    `json:"origin" yaml:"origin"`
			SourceAccountID string    `json:"sourceAccountId" yaml:"sourceAccountId"`
			TargetAccountID string    `json:"targetAccountId" yaml:"targetAccountId"`
			Amount          float64   `json:"amount" yaml:"amount"`
			Comment         string    `json:"comment" yaml:"comment"`
		} `body:""`
	}

	PutTransactionResponse struct {
		Tx Transaction
	}

	PatchTransactionRequest struct {
		Body struct {
			IssuedOn        *time.Time `json:"issuedOn" yaml:"issuedOn"`
			Origin          *string    `json:"origin" yaml:"origin"`
			SourceAccountID *string    `json:"sourceAccountId" yaml:"sourceAccountId"`
			TargetAccountID *string    `json:"targetAccountId" yaml:"targetAccountId"`
			Amount          *float64   `json:"amount" yaml:"amount"`
			Comment         *string    `json:"comment" yaml:"comment"`
		} `body:""`
	}

	PatchTransactionResponse struct {
		Tx Transaction
	}
)

func Delete(mgr Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := mgr.DeleteTransaction(r.Context(), chi.URLParam(r, "ID")); err == ErrInvalidID {
			w.WriteHeader(http.StatusBadRequest)

		} else if err == ErrNotFound {
			w.WriteHeader(http.StatusNotFound)

		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func Get(mgr Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tx, err := mgr.GetTransaction(r.Context(), chi.URLParam(r, "ID")); err == ErrInvalidID {
			w.WriteHeader(http.StatusBadRequest)

		} else if err == ErrNotFound || tx == nil {
			w.WriteHeader(http.StatusNotFound)

		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			render.Render(w, r, GetTransactionResponse{Tx: *tx})
		}
	}
}

func List(mgr Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paging := util.Paging{Page: 1, PageSize: 10}
		if !bind.Bind(r, w, &paging) {
			return
		} else if transactions, err := mgr.GetTransactionsList(r.Context(), paging.Page, paging.PageSize); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			render.Render(w, r, transactions)
		}
	}
}

func Patch(mgr Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PatchTransactionRequest
		if !bind.Bind(r, w, &req) {
			return

		} else if tx, err := mgr.PatchTransaction(r.Context(), chi.URLParam(r, "ID"),
			req.Body.IssuedOn, req.Body.Origin, req.Body.SourceAccountID, req.Body.TargetAccountID,
			req.Body.Amount, req.Body.Comment); err == ErrInvalidID {

			w.WriteHeader(http.StatusBadRequest)

		} else if err == ErrNotFound || tx == nil {
			w.WriteHeader(http.StatusNotFound)

		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			render.Render(w, r, PatchTransactionResponse{Tx: *tx})
		}
	}
}

func Post(mgr Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostTransactionRequest
		if !bind.Bind(r, w, &req) {
			return

		} else if tx, err := mgr.CreateTransaction(r.Context(), req.Body.IssuedOn, req.Body.Origin,
			req.Body.SourceAccountID, req.Body.TargetAccountID, req.Body.Amount, req.Body.Comment); err != nil {

			w.WriteHeader(http.StatusInternalServerError)

		} else {
			render.Render(w, r, PostTransactionResponse{Tx: *tx})
		}
	}
}

func Put(mgr Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PutTransactionRequest
		if !bind.Bind(r, w, &req) {
			return

		} else if tx, err := mgr.UpdateTransaction(r.Context(), chi.URLParam(r, "ID"), req.Body.IssuedOn,
			req.Body.Origin, req.Body.SourceAccountID, req.Body.TargetAccountID, req.Body.Amount, req.Body.Comment); err == ErrInvalidID {

			w.WriteHeader(http.StatusBadRequest)

		} else if err == ErrInternalError || tx == nil {
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			w.WriteHeader(http.StatusOK)
			render.Render(w, r, PutTransactionResponse{Tx: *tx})
		}
	}
}
