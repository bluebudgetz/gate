package accounts

import (
	"encoding/json"
	"fmt"
	"github.com/bluebudgetz/gate/internal/infra/render"
	"github.com/bluebudgetz/gate/internal/rest/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

//go:generate go-bindata -o ./assets.go -ignore ".*\\.go" -pkg accounts ./...

var getAccountsQueryDoc bson.A

func init() {
	if err := json.Unmarshal(MustAsset("get_accounts_query.json"), &getAccountsQueryDoc); err != nil {
		log.Fatal().Err(err).Msg("Failed loading 'get_accounts_query.json'")
	}
}

type GetInput struct {
	Page     uint `form:"_page"`
	PageSize uint `form:"_pageSize"`
}

type FetchedAccount struct {
	ID                  string      `json:"id"`
	CreatedOn           time.Time   `json:"createdOn"`
	UpdatedOn           *time.Time  `json:"updatedOn"`
	Name                string      `json:"name"`
	ParentID            *string     `json:"parentId"`
	TotalIncomingAmount json.Number `json:"totalIncomingAmount"`
	TotalOutgoingAmount json.Number `json:"totalOutgoingAmount"`
	Balance             json.Number `json:"balance"`
}

func GET(mongo *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := GetInput{Page: 1, PageSize: 10}
		if err := c.Bind(&input); err != nil {
			return
		}

		// Build query
		stages := make([]interface{}, len(getAccountsQueryDoc))
		copy(stages, getAccountsQueryDoc)
		stages = append(stages, buildAggregateFacetStage(input.Page, input.PageSize))

		// Fetch cursor
		cur, err := mongo.Database("bluebudgetz").Collection("accounts").Aggregate(c, stages)
		if err != nil {
			c.Error(err)
			return
		}
		defer cur.Close(c)

		// Load accounts from cursor
		var accounts []FetchedAccount
		if cur.Next(c) {
			var doc bson.M
			if err := cur.Decode(&doc); err != nil {
				c.Error(err)
				return
			}

			metadata := doc["metadata"].(bson.A)[0].(bson.M)
			c.Header("X-Total", fmt.Sprintf("%d", metadata["total"].(int32)))
			accounts = make([]FetchedAccount, 0)
			for _, dataArrayItem := range doc["data"].(bson.A) {
				accDoc := dataArrayItem.(bson.M)
				accounts = append(accounts, FetchedAccount{
					ID:                  accDoc["_id"].(primitive.ObjectID).Hex(),
					CreatedOn:           util.MustDateTime(accDoc["createdOn"]),
					UpdatedOn:           util.OptionalDateTime(accDoc["updatedOn"]),
					Name:                accDoc["name"].(string),
					ParentID:            util.OptionalObjectID(accDoc["parentId"]),
					TotalIncomingAmount: util.MustJsonNumber(accDoc["incoming"]),
					TotalOutgoingAmount: util.MustJsonNumber(accDoc["outgoing"]),
					Balance:             util.MustJsonNumber(accDoc["balance"]),
				})
			}
		}

		// If cursor failed, fail
		if err := cur.Err(); err != nil {
			c.Error(err)
			return
		}

		// Render
		render.Render(c, http.StatusOK, accounts)
	}
}

func buildAggregateFacetStage(page uint, pageSize uint) bson.M {
	return bson.M{
		"$facet": bson.M{
			"metadata": []bson.M{
				{"$count": "total"},
			},
			"data": []bson.M{
				{"$skip": (page - 1) * pageSize},
				{"$limit": pageSize},
			},
		},
	}
}
