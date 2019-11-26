package accounts

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
	PutAccountRequest struct {
		ID   string `path:"ID,required"`
		Body struct {
			Name     string  `json:"name" yaml:"name"`
			ParentID *string `json:"parentId" yaml:"parentId"`
		} `body:""`
	}
	PutAccountResponse struct {
		Account Account `json:"data"`
	}
)

func UpdateAccount(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := PutAccountRequest{}
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
			"updatedOn": t,
			"name":      req.Body.Name,
			"parentId":  util.OptionalObjectID(req.Body.ParentID),
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
					PutAccountResponse{
						Account: Account{
							ID:        req.ID,
							CreatedOn: t,
							UpdatedOn: nil,
							Name:      req.Body.Name,
							ParentID:  req.Body.ParentID,
						},
					},
				)
			}

		} else if result.Err() != nil {
			webutil.RenderWithStatusCode(
				w, r, http.StatusInternalServerError,
				errors.Wrap(result.Err(), "account update failed").AddTag("id", req.ID),
			)

		} else if err := result.Decode(&doc); err != nil {
			webutil.RenderWithStatusCode(
				w, r, http.StatusInternalServerError,
				errors.Wrap(err, "updated account decoding failed").AddTag("id", req.ID),
			)

		} else {
			webutil.RenderWithStatusCode(
				w, r, http.StatusOK,
				PutAccountResponse{
					Account: Account{
						ID:        req.ID,
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
