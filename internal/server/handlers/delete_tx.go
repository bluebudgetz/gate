package handlers

import (
	"net/http"

	"github.com/golangly/webutil"
)

type (
	DeleteTransactionRequest struct {
		ID string `path:"ID,required"`
	}
)

func DeleteTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := DeleteTransactionRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
		} else {
			DeleteRelationship(w, r, deleteTxQuery, map[string]interface{}{"id": req.ID})
		}
	}
}
