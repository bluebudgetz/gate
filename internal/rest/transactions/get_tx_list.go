package transactions

import (
	"net/http"

	"github.com/golangly/webutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/bluebudgetz/gate/internal/util"
)

type (
	GetTransactionListRequest struct {
		util.Paging
	}
	GetTransactionListResponse struct {
		Transactions []Transaction `json:"transactions"`
	}
)

func GetTransactionList(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := GetTransactionListRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
			return
		}

		// Fetch cursor
		limit := int64(req.Paging.PageSize)
		skip := int64((req.Paging.Page - 1) * req.Paging.PageSize)
		opts := options.FindOptions{
			Limit: &limit,
			Skip:  &skip,
			Sort:  bson.M{"issuedOn": 1},
		}
		cur, err := coll(mongoClient).Find(r.Context(), bson.M{}, &opts)
		if err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
			return
		}
		defer cur.Close(r.Context())

		// Load transactions from cursor
		var transactions = make([]Transaction, 0)
		for cur.Next(r.Context()) {
			var doc bson.M
			if err := cur.Decode(&doc); err != nil {
				webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
				return
			}
			transactions = append(transactions, Transaction{
				ID:              util.MustObjectIDHex(doc["_id"]),
				CreatedOn:       util.MustDateTime(doc["createdOn"]),
				UpdatedOn:       util.OptionalDateTime(doc["updatedOn"]),
				IssuedOn:        util.MustDateTime(doc["issuedOn"]),
				Origin:          doc["origin"].(string),
				SourceAccountID: doc["sourceAccountId"].(primitive.ObjectID).Hex(),
				TargetAccountID: doc["targetAccountId"].(primitive.ObjectID).Hex(),
				Amount:          doc["amount"].(float64),
				Comment:         doc["comment"].(string),
			})
		}

		// If cursor failed, fail
		if err := cur.Err(); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
		} else {
			webutil.RenderWithStatusCode(w, r, http.StatusOK, GetTransactionListResponse{Transactions: transactions})
		}
	}
}
