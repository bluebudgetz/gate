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
	GetTransactionRequest struct {
		ID string `path:"ID,required"`
	}
	GetTransactionResponse struct {
		Tx Transaction `json:"data"`
	}
)

func GetTransaction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := GetTransactionRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
		} else {
			GetNode(w, r, getTxQuery, map[string]interface{}{"id": req.ID}, func(rec neo4j.Record) (interface{}, error) {
				if tx, ok := rec.Get("tx"); !ok {
					return nil, errors.New("tx node not found in results")
				} else if txRel, ok := tx.(neo4j.Relationship); !ok {
					return nil, errors.New("tx node mismatch")
				} else if src, ok := rec.Get("src"); !ok {
					return nil, errors.New("src node not found in results")
				} else if srcNode, ok := src.(neo4j.Node); !ok {
					return nil, errors.New("src node mismatch")
				} else if dst, ok := rec.Get("dst"); !ok {
					return nil, errors.New("dst node not found in results")
				} else if dstNode, ok := dst.(neo4j.Node); !ok {
					return nil, errors.New("dst node mismatch")
				} else {
					return GetTransactionResponse{
						Tx: Transaction{
							ID:              txRel.Props()["id"].(string),
							CreatedOn:       txRel.Props()["createdOn"].(time.Time),
							UpdatedOn:       util.OptionalDateTime(txRel.Props()["updatedOn"]),
							IssuedOn:        txRel.Props()["issuedOn"].(time.Time),
							Origin:          txRel.Props()["origin"].(string),
							SourceAccountID: srcNode.Props()["id"].(string),
							TargetAccountID: dstNode.Props()["id"].(string),
							Amount:          txRel.Props()["amount"].(float64),
							Comment:         txRel.Props()["comment"].(string),
						},
					}, nil
				}
			})
		}
	}
}
