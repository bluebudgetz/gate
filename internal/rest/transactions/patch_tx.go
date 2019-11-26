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
	PatchTransactionRequest struct {
		ID   string `path:"ID,required"`
		Body struct {
			IssuedOn        *time.Time `json:"issuedOn" yaml:"issuedOn"`
			Origin          *string    `json:"origin" yaml:"origin"`
			SourceAccountID *string    `json:"sourceAccountId" yaml:"sourceAccountId"`
			TargetAccountID *string    `json:"targetAccountId" yaml:"targetAccountId"`
			Amount          *float64   `json:"amount" yaml:"amount"`
			Comment         *string    `json:"comment" yaml:"comment"`
		} `body:""`
	}
	PatchTransactionResponse struct {
		Transaction Transaction
	}
)

func PatchTransaction(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := PatchTransactionRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
			return
		} else if req.ID == "" {
			webutil.RenderWithStatusCode(w, r, http.StatusBadRequest, nil)
			return
		}

		// Build patch spec
		doc := bson.M{"updatedOn": time.Now()}
		if req.Body.IssuedOn != nil {
			doc["issuedOn"] = *req.Body.IssuedOn
		}
		if req.Body.Origin != nil {
			doc["origin"] = *req.Body.Origin
		}
		if req.Body.SourceAccountID != nil {
			doc["sourceAccountId"] = util.MustObjectID(*req.Body.SourceAccountID)
		}
		if req.Body.TargetAccountID != nil {
			doc["targetAccountID"] = util.MustObjectID(*req.Body.TargetAccountID)
		}
		if req.Body.Amount != nil {
			doc["amount"] = *req.Body.Amount
		}
		if req.Body.Comment != nil {
			doc["comment"] = *req.Body.Comment
		}

		// Patch it
		after := options.After
		opts := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
		updateSpec := bson.M{"$set": doc}
		if result := coll(mongoClient).FindOneAndUpdate(r.Context(), bson.M{"_id": util.MustObjectID(req.ID)}, updateSpec, opts); result.Err() == mongo.ErrNoDocuments {
			webutil.RenderWithStatusCode(w, r, http.StatusNotFound, nil)

		} else if result.Err() != nil {
			webutil.RenderWithStatusCode(
				w, r, http.StatusInternalServerError,
				errors.Wrap(result.Err(), "transaction patching failed").AddTag("id", req.ID),
			)

		} else if err := result.Decode(&doc); err != nil {
			webutil.RenderWithStatusCode(
				w, r, http.StatusInternalServerError,
				errors.Wrap(err, "patched transaction decoding failed").AddTag("id", req.ID),
			)

		} else {
			webutil.RenderWithStatusCode(
				w, r, http.StatusOK,
				PatchTransactionResponse{
					Transaction: Transaction{
						ID:              util.MustObjectIDHex(doc["_id"]),
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
