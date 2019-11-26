package accounts

import (
	"net/http"

	"github.com/golangly/errors"
	"github.com/golangly/webutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bluebudgetz/gate/internal/util"
)

type (
	GetAccountListRequest struct {
		util.Paging
	}
	GetAccountListResponse struct {
		Accounts []AccountWithBalance `json:"data"`
	}
)

func GetAccountList(mongoClient *mongo.Client) http.HandlerFunc {
	newQuery := func(paging util.Paging) []interface{} {
		stages := make([]interface{}, len(getAccountsListQueryDoc))
		copy(stages, getAccountsListQueryDoc)
		stages = append(stages, bson.M{
			"$facet": bson.M{
				"metadata": []bson.M{
					{"$count": "total"},
				},
				"data": []bson.M{
					{"$skip": (paging.Page - 1) * paging.PageSize},
					{"$limit": paging.PageSize},
				},
			},
		})
		return stages
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := GetAccountListRequest{
			Paging: util.Paging{Page: 1, PageSize: 10},
		}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
			return
		}

		query := newQuery(req.Paging)
		cur, err := coll(mongoClient).Aggregate(r.Context(), query)
		if err != nil {
			webutil.RenderWithStatusCode(
				w, r, http.StatusInternalServerError,
				errors.Wrap(err, "accounts query failed").AddTag("query", query),
			)
			return
		}
		defer cur.Close(r.Context())

		// Load accounts from cursor
		var accounts = make([]AccountWithBalance, 0)
		if cur.Next(r.Context()) {
			var doc bson.M
			if err := cur.Decode(&doc); err != nil {
				webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
				return
			}
			for _, dataArrayItem := range doc["data"].(bson.A) {
				accDoc := dataArrayItem.(bson.M)
				accounts = append(accounts, AccountWithBalance{
					Account: Account{
						ID:        util.MustObjectIDHex(accDoc["_id"]),
						CreatedOn: util.MustDateTime(accDoc["createdOn"]),
						UpdatedOn: util.OptionalDateTime(accDoc["updatedOn"]),
						Name:      accDoc["name"].(string),
						ParentID:  util.OptionalObjectIDHex(accDoc["parentId"]),
					},
					TotalIncomingAmount: accDoc["incoming"].(float64),
					TotalOutgoingAmount: accDoc["outgoing"].(float64),
					Balance:             accDoc["balance"].(float64),
				})
			}
		}

		// If cursor failed, fail
		if err := cur.Err(); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
		} else {
			webutil.RenderWithStatusCode(w, r, http.StatusOK, GetAccountListResponse{Accounts: accounts})
		}
	}
}
