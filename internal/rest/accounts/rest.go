package accounts

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/bluebudgetz/gate/internal/infra/bind"
	"github.com/bluebudgetz/gate/internal/infra/render"
	"github.com/bluebudgetz/gate/internal/util"
)

type (
	GetAccountResponse struct {
		ID        string     `json:"id" yaml:"id"`
		CreatedOn time.Time  `json:"createdOn" yaml:"createdOn"`
		UpdatedOn *time.Time `json:"updatedOn" yaml:"updatedOn"`
		Name      string     `json:"name" yaml:"name"`
		ParentID  *string    `json:"parentId" yaml:"parentId"`
	}

	PostAccountRequest struct {
		Body struct {
			Name     string  `json:"name" yaml:"name"`
			ParentID *string `json:"parentId" yaml:"parentId"`
		} `body:""`
	}

	PostAccountResponse struct {
		ID        string    `json:"id" yaml:"id"`
		CreatedOn time.Time `json:"createdOn" yaml:"createdOn"`
		Name      string    `json:"name" yaml:"name"`
		ParentID  *string   `json:"parentId" yaml:"parentId"`
	}

	PutAccountRequest struct {
		Body struct {
			Name     string  `json:"name" yaml:"name"`
			ParentID *string `json:"parentId" yaml:"parentId"`
		} `body:""`
	}

	PutAccountResponse struct {
		ID        string     `json:"id" yaml:"id"`
		CreatedOn time.Time  `json:"createdOn" yaml:"createdOn"`
		UpdatedOn *time.Time `json:"updatedOn" yaml:"updatedOn"`
		Name      string     `json:"name" yaml:"name"`
		ParentID  *string    `json:"parentId" yaml:"parentId"`
	}

	PatchAccountRequest struct {
		Body struct {
			Name     *string `json:"name" yaml:"name"`
			ParentID *string `json:"parentId" yaml:"parentId"`
		} `body:""`
	}

	PatchAccountResponse struct {
		ID        string    `json:"id" yaml:"id"`
		CreatedOn time.Time `json:"createdOn" yaml:"createdOn"`
		UpdatedOn time.Time `json:"updatedOn" yaml:"updatedOn"`
		Name      string    `json:"name" yaml:"name"`
		ParentID  *string   `json:"parentId" yaml:"parentId"`
	}
)

func Delete(mgr Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := mgr.DeleteAccount(r.Context(), chi.URLParam(r, "ID")); err == ErrInvalidID {
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
		if acc, err := mgr.GetAccount(r.Context(), chi.URLParam(r, "ID")); err == ErrInvalidID {
			w.WriteHeader(http.StatusBadRequest)

		} else if err == ErrNotFound || acc == nil {
			w.WriteHeader(http.StatusNotFound)

		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			render.Render(w, r, GetAccountResponse{
				ID:        acc.ID,
				CreatedOn: acc.CreatedOn,
				UpdatedOn: acc.UpdatedOn,
				Name:      acc.Name,
				ParentID:  acc.ParentID,
			})
		}
	}
}

func List(mgr Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		paging := util.Paging{Page: 1, PageSize: 10}
		if !bind.Bind(r, w, &paging) {
			return
		} else if accounts, err := mgr.GetAccountsList(r.Context(), paging.Page, paging.PageSize); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			render.Render(w, r, accounts)
		}
	}
}

func Patch(mgr Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PatchAccountRequest
		if !bind.Bind(r, w, &req) {
			return

		} else if acc, err := mgr.PatchAccount(r.Context(), chi.URLParam(r, "ID"), req.Body.Name, req.Body.ParentID); err == ErrInvalidID {
			w.WriteHeader(http.StatusBadRequest)

		} else if err == ErrNotFound || acc == nil {
			w.WriteHeader(http.StatusNotFound)

		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			render.Render(w, r, PatchAccountResponse{
				ID:        acc.ID,
				CreatedOn: acc.CreatedOn,
				UpdatedOn: *acc.UpdatedOn,
				Name:      acc.Name,
				ParentID:  acc.ParentID,
			})
		}
	}
}

func Post(mgr Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostAccountRequest
		if !bind.Bind(r, w, &req) {
			return

		} else if acc, err := mgr.CreateAccount(r.Context(), req.Body.Name, req.Body.ParentID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			render.Render(w, r, PostAccountResponse{
				ID:        acc.ID,
				CreatedOn: acc.CreatedOn,
				Name:      acc.Name,
				ParentID:  acc.ParentID,
			})
		}
	}
}

func Put(mgr Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PutAccountRequest
		if !bind.Bind(r, w, &req) {
			return

		} else if account, err := mgr.UpdateAccount(r.Context(), chi.URLParam(r, "ID"), req.Body.Name, req.Body.ParentID); err == ErrInvalidID {
			w.WriteHeader(http.StatusBadRequest)

		} else if err == ErrInternalError || account == nil {
			w.WriteHeader(http.StatusInternalServerError)

		} else {
			w.WriteHeader(http.StatusOK)
			render.Render(w, r, PutAccountResponse{
				ID:        account.ID,
				CreatedOn: account.CreatedOn,
				UpdatedOn: account.UpdatedOn,
				Name:      account.Name,
				ParentID:  account.ParentID,
			})
		}
	}
}
