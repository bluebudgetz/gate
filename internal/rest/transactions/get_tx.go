package transactions

import (
	"net/http"

	"github.com/golangly/errors"
	"github.com/golangly/webutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bluebudgetz/gate/internal/util"
)

type (
	GetTransactionRequest struct {
		ID string `path:"ID,required"`
	}
	GetTransactionResponse struct {
		Tx Transaction
	}
)

func GetTransaction(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := GetTransactionRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)

		} else if req.ID == "" {
			webutil.RenderWithStatusCode(w, r, http.StatusBadRequest, err)

		} else if result := coll(mongoClient).FindOne(r.Context(), bson.M{"_id": util.MustObjectID(req.ID)}); result.Err() == mongo.ErrNoDocuments {
			webutil.RenderWithStatusCode(w, r, http.StatusNotFound, nil)

		} else if result.Err() != nil {
			webutil.RenderWithStatusCode(
				w, r, http.StatusInternalServerError,
				errors.Wrap(result.Err(), "transaction lookup failed").AddTag("id", req.ID),
			)

		} else {
			doc := bson.M{}
			if err := result.Decode(&doc); err != nil {
				webutil.RenderWithStatusCode(
					w, r, http.StatusInternalServerError,
					errors.Wrap(err, "transaction decoding failed").AddTag("id", req.ID),
				)
			} else {
				webutil.RenderWithStatusCode(
					w, r, http.StatusOK,
					GetTransactionResponse{
						Tx: Transaction{
							ID:              util.MustObjectIDHex(doc["_id"]),
							CreatedOn:       util.MustDateTime(doc["createdOn"]),
							UpdatedOn:       util.OptionalDateTime(doc["updatedOn"]),
							IssuedOn:        util.MustDateTime(doc["issuedOn"]),
							Origin:          doc["origin"].(string),
							SourceAccountID: doc["sourceAccountId"].(primitive.ObjectID).Hex(),
							TargetAccountID: doc["targetAccountId"].(primitive.ObjectID).Hex(),
							Amount:          doc["amount"].(float64),
							Comment:         doc["comment"].(string),
						},
					},
				)
			}
		}
	}
}
