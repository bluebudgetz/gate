package handlers

import (
	"net/http"
	"time"

	"github.com/golangly/errors"
	"github.com/golangly/webutil"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/bluebudgetz/gate/internal/services"
	"github.com/bluebudgetz/gate/internal/util"
)

type (
	GetTransactionListRequest struct {
		util.Paging
	}
	GetTransactionListItemData struct {
		ID              string     `json:"id" yaml:"id"`
		CreatedOn       time.Time  `json:"createdOn" yaml:"createdOn"`
		UpdatedOn       *time.Time `json:"updatedOn" yaml:"updatedOn"`
		IssuedOn        time.Time  `json:"issuedOn" yaml:"issuedOn"`
		Origin          string     `json:"origin" yaml:"origin"`
		SourceAccountID string     `json:"sourceAccountId" yaml:"sourceAccountId"`
		TargetAccountID string     `json:"targetAccountId" yaml:"targetAccountId"`
		Amount          float64    `json:"amount" yaml:"amount"`
		Comment         string     `json:"comment" yaml:"comment"`
	}
	GetTransactionListResponse struct {
		Transactions []GetTransactionListItemData `json:"data"`
	}
)

func GetTransactionList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := GetTransactionListRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
			return
		}

		// Query
		result, err := services.GetNeo4jSession(r.Context()).Run(getTxListQuery, map[string]interface{}{
			"skip":  (req.Paging.Page - 1) * req.Paging.PageSize,
			"limit": req.Paging.PageSize,
		})
		if err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, errors.Wrap(err, "query failed"))
			return
		}

		// Read transactions
		var transactions = make([]GetTransactionListItemData, 0)
		for result.Next() {
			rec := result.Record()
			if tx, ok := rec.Get("tx"); !ok {
				webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, errors.New("failed getting tx node"))
				return
			} else if src, ok := rec.Get("src"); !ok {
				webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, errors.New("failed getting src node"))
				return
			} else if dst, ok := rec.Get("src"); !ok {
				webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, errors.New("failed getting dst node"))
				return
			} else {
				txNode := tx.(neo4j.Node).Props()
				srcNode := src.(neo4j.Node).Props()
				dstNode := dst.(neo4j.Node).Props()
				transactions = append(transactions, GetTransactionListItemData{
					ID:              txNode["id"].(string),
					CreatedOn:       txNode["createdOn"].(time.Time),
					UpdatedOn:       txNode["updatedOn"].(*time.Time),
					IssuedOn:        txNode["issuedOn"].(time.Time),
					Origin:          txNode["origin"].(string),
					SourceAccountID: srcNode["id"].(string),
					TargetAccountID: dstNode["id"].(string),
					Amount:          txNode["amount"].(float64),
					Comment:         txNode["comment"].(string),
				})
			}
		}
		if err := result.Err(); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, errors.Wrap(err, "query failed"))
			return
		}
		webutil.RenderWithStatusCode(w, r, http.StatusOK, GetTransactionListResponse{Transactions: transactions})
	}
}
