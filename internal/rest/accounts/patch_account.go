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
	PatchAccountRequest struct {
		ID   string `path:"ID,required"`
		Body struct {
			Name     *string `json:"name" yaml:"name"`
			ParentID *string `json:"parentId" yaml:"parentId"`
		} `body:""`
	}
	PatchAccountResponse struct {
		Account Account `json:"data"`
	}
)

func PatchAccount(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := PatchAccountRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
			return
		} else if req.ID == "" {
			webutil.RenderWithStatusCode(w, r, http.StatusBadRequest, nil)
			return
		}

		// Build patch spec
		doc := bson.M{"updatedOn": time.Now()}
		if req.Body.Name != nil {
			doc["name"] = *req.Body.Name
		}
		if req.Body.ParentID != nil {
			doc["parentId"] = util.MustObjectID(*req.Body.ParentID)
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
				errors.Wrap(result.Err(), "account patching failed").AddTag("id", req.ID),
			)

		} else if err := result.Decode(&doc); err != nil {
			webutil.RenderWithStatusCode(
				w, r, http.StatusInternalServerError,
				errors.Wrap(err, "patched account decoding failed").AddTag("id", req.ID),
			)

		} else {
			webutil.RenderWithStatusCode(
				w, r, http.StatusOK,
				PatchAccountResponse{
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
