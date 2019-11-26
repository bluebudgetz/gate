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
	GetAccountRequest struct {
		ID string `path:"ID,required"`
	}
	GetAccountResponse struct {
		Account Account `json:"data"`
	}
)

func GetAccount(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := GetAccountRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)

		} else if req.ID == "" {
			webutil.RenderWithStatusCode(w, r, http.StatusBadRequest, err)

		} else if result := coll(mongoClient).FindOne(r.Context(), bson.M{"_id": util.MustObjectID(req.ID)}); result.Err() == mongo.ErrNoDocuments {
			webutil.RenderWithStatusCode(w, r, http.StatusNotFound, nil)

		} else if result.Err() != nil {
			webutil.RenderWithStatusCode(
				w, r, http.StatusInternalServerError,
				errors.Wrap(result.Err(), "account lookup failed").AddTag("id", req.ID),
			)

		} else {
			doc := bson.M{}
			if err := result.Decode(&doc); err != nil {
				webutil.RenderWithStatusCode(
					w, r, http.StatusInternalServerError,
					errors.Wrap(err, "account decoding failed").AddTag("id", req.ID),
				)
			} else {
				webutil.RenderWithStatusCode(
					w, r, http.StatusOK,
					GetAccountResponse{
						Account: Account{
							ID:        util.MustObjectIDHex(doc["_id"]),
							CreatedOn: util.MustDateTime(doc["createdOn"]),
							UpdatedOn: util.OptionalDateTime(doc["updatedOn"]),
							Name:      doc["name"].(string),
							ParentID:  util.OptionalObjectIDHex(doc["parentId"]),
						},
					},
				)
			}
		}
	}
}
