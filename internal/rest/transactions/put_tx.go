package transactions

import (
	"net/http"
	"time"

	"github.com/golangly/errors"
	"github.com/golangly/webutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/bluebudgetz/gate/internal/util"
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
		Transaction Transaction
	}
)

func UpdateTransaction(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := PutTransactionRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
			return
		} else if req.ID == "" {
			webutil.RenderWithStatusCode(w, r, http.StatusBadRequest, nil)
			return
		}

		t := time.Now()

		// Update spec
		filter := bson.M{"_id": util.MustObjectID(req.ID)}
		doc := bson.M{
			"updatedOn":       t,
			"issuedOn":        req.Body.IssuedOn,
			"origin":          req.Body.Origin,
			"sourceAccountId": util.MustObjectID(req.Body.SourceAccountID),
			"targetAccountId": util.MustObjectID(req.Body.TargetAccountID),
			"amount":          req.Body.Amount,
			"comment":         req.Body.Comment,
		}

		// Check if the document exists or not; we're not simply using "FindOneAnd{Update|Replace}" + upsert, because
		// this would create an inconsistency wrt "createdOn" and "updatedOn" when new document would be created.
		after := options.After
		opts := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
		if result := coll(mongoClient).FindOneAndUpdate(r.Context(), filter, &bson.M{"$set": doc}, opts); result.Err() == mongo.ErrNoDocuments {

			// Create it
			doc["_id"] = util.MustObjectID(req.ID)
			doc["createdOn"] = t
			doc["updatedOn"] = nil
			if _, err := coll(mongoClient).InsertOne(r.Context(), &doc); err != nil {
				webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
			} else {
				webutil.RenderWithStatusCode(
					w, r, http.StatusOK,
					PutTransactionResponse{
						Transaction: Transaction{
							ID:              req.ID,
							CreatedOn:       util.MustDateTime(doc["createdOn"]),
							UpdatedOn:       util.OptionalDateTime(doc["updatedOn"]),
							IssuedOn:        util.MustDateTime(doc["issuedOn"]),
							Origin:          doc["origin"].(string),
							SourceAccountID: util.MustObjectIDHex(doc["sourceAccountId"]),
							TargetAccountID: util.MustObjectIDHex(doc["targetAccountId"]),
							Amount:          doc["amount"].(float64),
							Comment:         doc["comment"].(string),
						},
					},
				)
			}

		} else if result.Err() != nil {
			webutil.RenderWithStatusCode(
				w, r, http.StatusInternalServerError,
				errors.Wrap(result.Err(), "transaction update failed").AddTag("id", req.ID),
			)

		} else if err := result.Decode(&doc); err != nil {
			webutil.RenderWithStatusCode(
				w, r, http.StatusInternalServerError,
				errors.Wrap(err, "updated transaction decoding failed").AddTag("id", req.ID),
			)

		} else {
			webutil.RenderWithStatusCode(
				w, r, http.StatusOK,
				PutTransactionResponse{
					Transaction: Transaction{
						ID:              req.ID,
						CreatedOn:       util.MustDateTime(doc["createdOn"]),
						UpdatedOn:       util.OptionalDateTime(doc["updatedOn"]),
						IssuedOn:        util.MustDateTime(doc["issuedOn"]),
						Origin:          doc["origin"].(string),
						SourceAccountID: util.MustObjectIDHex(doc["sourceAccountId"]),
						TargetAccountID: util.MustObjectIDHex(doc["targetAccountId"]),
						Amount:          doc["amount"].(float64),
						Comment:         doc["comment"].(string),
					},
				},
			)
		}
	}
}
