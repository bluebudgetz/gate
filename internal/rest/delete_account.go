package rest

import (
	"net/http"

	"github.com/golangly/webutil"
)

type (
	DeleteAccountRequest struct {
		ID string `path:"ID,required"`
	}
)

func DeleteAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := DeleteAccountRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
		} else {
			DeleteNode(w, r, deleteAccountQuery, map[string]interface{}{"id": req.ID})
		}
	}
}
