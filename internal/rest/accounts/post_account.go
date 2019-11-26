package accounts

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
	PostAccountRequest struct {
		Body struct {
			Name     string  `json:"name" yaml:"name"`
			ParentID *string `json:"parentId" yaml:"parentId"`
		} `body:""`
	}
	PostAccountResponse struct {
		Account Account `json:"data"`
	}
)

func CreateAccount(mongoClient *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := PostAccountRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)
			return
		}

		doc := bson.M{
			"createdOn": time.Now(),
			"updatedOn": nil,
			"name":      req.Body.Name,
			"parentId":  util.OptionalObjectID(req.Body.ParentID),
		}
		if result, err := coll(mongoClient).InsertOne(r.Context(), &doc); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)

		} else {
			webutil.RenderWithStatusCode(
				w, r, http.StatusOK,
				PostAccountResponse{
					Account: Account{
						ID:        result.InsertedID.(primitive.ObjectID).Hex(),
						CreatedOn: util.MustDateTime(doc["createdOn"]),
						UpdatedOn: nil,
						Name:      req.Body.Name,
						ParentID:  req.Body.ParentID,
					},
				},
			)
		}
	}
}
