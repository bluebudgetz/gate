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
	PostAccountRequest struct {
		Body struct {
			Name     string  `json:"name" yaml:"name"`
			ParentID *string `json:"parentId" yaml:"parentId"`
		} `body:""`
	}
	PostAccountData     struct{ GetAccountData }
	PostAccountResponse struct {
		Account PostAccountData `json:"data"`
	}
)

func CreateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := PostAccountRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
		} else {
			PostNode(w, r, postAccountQuery, map[string]interface{}{"name": req.Body.Name, "parentId": req.Body.ParentID}, func(rec neo4j.Record) (interface{}, error) {
				if node, ok := rec.Get("account"); !ok {
					return nil, errors.Wrap(err, "node not found in results")
				} else if accNode, ok := node.(neo4j.Node); !ok {
					return nil, errors.Wrap(err, "node mismatch")
				} else {
					return PostAccountResponse{
						Account: PostAccountData{
							GetAccountData: GetAccountData{
								ID:        accNode.Props()["id"].(string),
								CreatedOn: accNode.Props()["createdOn"].(time.Time),
								UpdatedOn: util.OptionalDateTime(accNode.Props()["updatedOn"]),
								Name:      accNode.Props()["name"].(string),
							},
						},
					}, nil
				}
			})
		}
	}
}
