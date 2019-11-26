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
	DeleteAccountRequest struct {
		ID string `path:"ID,required"`
	}
)

func DeleteAccount(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := DeleteAccountRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)

		} else if req.ID == "" {
			webutil.RenderWithStatusCode(w, r, http.StatusBadRequest, err)

		} else if result, err := coll(mongoClient).DeleteOne(r.Context(), bson.M{"_id": util.MustObjectID(req.ID)}); err != nil {
			webutil.RenderWithStatusCode(
				w, r, http.StatusInternalServerError,
				errors.Wrap(err, "account deletion failed").AddTag("id", req.ID),
			)

		} else if result.DeletedCount == 0 {
			webutil.RenderWithStatusCode(w, r, http.StatusNotFound, nil)

		} else if result.DeletedCount == 1 {
			webutil.RenderWithStatusCode(w, r, http.StatusNoContent, nil)

		} else {
			webutil.RenderWithStatusCode(
				w, r,
				http.StatusInternalServerError,
				errors.New("too many accounts deleted").
					AddTag("deleted", result.DeletedCount).
					AddTag("id", req.ID),
			)
		}
	}
}
