package transactions

import (
	"net/http"
	"time"

	"github.com/golangly/webutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bluebudgetz/gate/internal/util"
)

type (
	PostTransactionRequest struct {
		Body struct {
			IssuedOn        time.Time `json:"issuedOn" yaml:"issuedOn"`
			Origin          string    `json:"origin" yaml:"origin"`
			SourceAccountID string    `json:"sourceAccountId" yaml:"sourceAccountId"`
			TargetAccountID string    `json:"targetAccountId" yaml:"targetAccountId"`
			Amount          float64   `json:"amount" yaml:"amount"`
			Comment         string    `json:"comment" yaml:"comment"`
		} `body:""`
	}
	PostTransactionResponse struct {
		Transaction Transaction
	}
)

func CreateTransaction(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := PostTransactionRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
			return
		}

		doc := bson.M{
			"createdOn":       time.Now(),
			"updatedOn":       nil,
			"issuedOn":        req.Body.IssuedOn,
			"origin":          req.Body.Origin,
			"sourceAccountId": util.MustObjectID(req.Body.SourceAccountID),
			"targetAccountId": util.MustObjectID(req.Body.TargetAccountID),
			"amount":          req.Body.Amount,
			"comment":         req.Body.Comment,
		}
		if result, err := coll(mongoClient).InsertOne(r.Context(), &doc); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)

		} else {
			webutil.RenderWithStatusCode(
				w, r, http.StatusOK,
				PostTransactionResponse{
					Transaction: Transaction{
						ID:              result.InsertedID.(primitive.ObjectID).Hex(),
						CreatedOn:       util.MustDateTime(doc["createdOn"]),
						UpdatedOn:       nil,
						IssuedOn:        util.MustDateTime(doc["issuedOn"]),
						Origin:          req.Body.Origin,
						SourceAccountID: req.Body.SourceAccountID,
						TargetAccountID: req.Body.TargetAccountID,
						Amount:          req.Body.Amount,
						Comment:         req.Body.Comment,
					},
				},
			)
		}
	}
}
