package handlers

import (
	"net/http"
	"time"

	"github.com/golangly/errors"
	"github.com/golangly/webutil"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/bluebudgetz/gate/internal/util"
)

type (
	GetAccountRequest struct {
		ID string `path:"ID,required"`
	}
	GetAccountData struct {
		ID        string     `json:"id" yaml:"id"`
		CreatedOn time.Time  `json:"createdOn" yaml:"createdOn"`
		UpdatedOn *time.Time `json:"updatedOn" yaml:"updatedOn"`
		Name      string     `json:"name" yaml:"name"`
	}
	GetAccountResponse struct {
		Account GetAccountData `json:"data"`
	}
)

func GetAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := GetAccountRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
		} else {
			GetNode(w, r, getAccountQuery, map[string]interface{}{"id": req.ID}, func(rec neo4j.Record) (interface{}, error) {
				if node, ok := rec.Get("account"); !ok {
					return nil, errors.New("account not found in results")
				} else if accNode, ok := node.(neo4j.Node); !ok {
					return nil, errors.New("account node mismatch")
				} else {
					return GetAccountResponse{
						Account: GetAccountData{
							ID:        accNode.Props()["id"].(string),
							CreatedOn: accNode.Props()["createdOn"].(time.Time),
							UpdatedOn: util.OptionalDateTime(accNode.Props()["updatedOn"]),
							Name:      accNode.Props()["name"].(string),
						},
					}, nil
				}
			})
		}
	}
}
