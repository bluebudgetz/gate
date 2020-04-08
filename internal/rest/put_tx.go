package rest

import (
	"net/http"
	"time"

	"github.com/golangly/errors"
	"github.com/golangly/webutil"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type (
	PutTransactionRequest struct {
		ID   string `path:"ID,required"`
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
		Transaction Transaction `json:"data"`
	}
)

func UpdateTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := PutTransactionRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
		} else {
			PutNode(w, r, putTxQuery, map[string]interface{}{
				"id":       req.ID,
				"issuedOn": req.Body.IssuedOn,
				"origin":   req.Body.Origin,
				"amount":   req.Body.Amount,
				"comment":  req.Body.Comment,
			}, func(rec neo4j.Record) (interface{}, error) {
				if src, ok := rec.Get("src"); !ok {
					return nil, errors.Wrap(err, "src not found in results")
				} else if srcNode, ok := src.(neo4j.Node); !ok {
					return nil, errors.Wrap(err, "src node mismatch")
				} else if tx, ok := rec.Get("tx"); !ok {
					return nil, errors.Wrap(err, "tx not found in results")
				} else if txNode, ok := tx.(neo4j.Node); !ok {
					return nil, errors.Wrap(err, "tx node mismatch")
				} else if dst, ok := rec.Get("dst"); !ok {
					return nil, errors.Wrap(err, "dst not found in results")
				} else if dstNode, ok := dst.(neo4j.Node); !ok {
					return nil, errors.Wrap(err, "dst node mismatch")
				} else {
					return PutTransactionResponse{
						Transaction: Transaction{
							ID:              txNode.Props()["id"].(string),
							CreatedOn:       txNode.Props()["createdOn"].(time.Time),
							UpdatedOn:       txNode.Props()["updatedOn"].(*time.Time),
							IssuedOn:        txNode.Props()["issuedOn"].(time.Time),
							Origin:          txNode.Props()["origin"].(string),
							SourceAccountID: srcNode.Props()["id"].(string),
							TargetAccountID: dstNode.Props()["id"].(string),
							Amount:          txNode.Props()["amount"].(float64),
							Comment:         txNode.Props()["comment"].(string),
						},
					}, nil
				}
			})
		}
	}
}
