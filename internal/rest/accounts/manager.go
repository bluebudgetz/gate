package accounts

//go:generate go-bindata -o ./assets.go -ignore ".*\\.go" -pkg accounts ./...

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/bluebudgetz/gate/internal/util"
)

type (
	Manager interface {
		CreateAccount(ctx context.Context, name string, parentID *string) (*Account, error)
		DeleteAccount(ctx context.Context, id string) error
		GetAccount(ctx context.Context, id string) (*Account, error)
		GetAccountsList(ctx context.Context, page uint, pageSize uint) ([]AccountWithBalance, error)
		PatchAccount(ctx context.Context, id string, name *string, parentID *string) (*Account, error)
		UpdateAccount(ctx context.Context, id string, name string, parentID *string) (*Account, error)
	}

	manager struct {
		mongoClient *mongo.Client
	}

	Account struct {
		ID        string     `json:"id" yaml:"id"`
		CreatedOn time.Time  `json:"createdOn" yaml:"createdOn"`
		UpdatedOn *time.Time `json:"updatedOn" yaml:"updatedOn"`
		Name      string     `json:"name" yaml:"name"`
		ParentID  *string    `json:"parentId" yaml:"parentId"`
	}

	AccountWithBalance struct {
		Account
		TotalIncomingAmount float64 `json:"totalIncomingAmount" yaml:"totalIncomingAmount"`
		TotalOutgoingAmount float64 `json:"totalOutgoingAmount" yaml:"totalOutgoingAmount"`
		Balance             float64 `json:"balance" yaml:"balance"`
	}
)

var (
	getAccountsListQueryDoc bson.A

	ErrInvalidID     = fmt.Errorf("invalid account ID")
	ErrNotFound      = fmt.Errorf("account not found")
	ErrInternalError = fmt.Errorf("internal error")
)

func init() {
	if err := json.Unmarshal(MustAsset("get_accounts_list_query.json"), &getAccountsListQueryDoc); err != nil {
		log.Fatal().Err(err).Msg("Failed loading 'get_accounts_query.json'")
	}
}

func NewManager(mongoClient *mongo.Client) Manager {
	return &manager{mongoClient}
}

func (m *manager) coll() *mongo.Collection {
	return m.mongoClient.Database("bluebudgetz").Collection("accounts")
}

func (m *manager) DeleteAccount(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidID

	} else if result, err := m.coll().DeleteOne(ctx, bson.M{"_id": util.OptionalObjectID(id)}); err == mongo.ErrNoDocuments {
		return ErrNotFound

	} else if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed deleting account")
		return ErrInternalError

	} else if result != nil && result.DeletedCount == 0 {
		return ErrNotFound

	} else {
		return nil
	}
}

func (m *manager) GetAccount(ctx context.Context, id string) (*Account, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	doc := bson.M{}

	if result := m.coll().FindOne(ctx, bson.M{"_id": util.OptionalObjectID(id)}); result.Err() == mongo.ErrNoDocuments {
		return nil, ErrNotFound

	} else if result.Err() != nil {
		log.Error().Err(result.Err()).Str("id", id).Msg("Failed looking up account")
		return nil, ErrInternalError

	} else if err := result.Decode(&doc); err == mongo.ErrNoDocuments {
		return nil, ErrNotFound

	} else if err != nil {
		log.Error().Err(result.Err()).Str("id", id).Msg("Failed looking up account")
		return nil, ErrInternalError

	} else {
		return &Account{
			ID:        *util.OptionalObjectIDHex(doc["_id"]),
			CreatedOn: util.MustDateTime(doc["createdOn"]),
			UpdatedOn: util.OptionalDateTime(doc["updatedOn"]),
			Name:      doc["name"].(string),
			ParentID:  util.OptionalObjectIDHex(doc["parentId"]),
		}, nil
	}
}

func (m *manager) GetAccountsList(ctx context.Context, page uint, pageSize uint) ([]AccountWithBalance, error) {

	// Build query
	stages := make([]interface{}, len(getAccountsListQueryDoc))
	copy(stages, getAccountsListQueryDoc)
	stages = append(stages, bson.M{
		"$facet": bson.M{
			"metadata": []bson.M{
				{"$count": "total"},
			},
			"data": []bson.M{
				{"$skip": (page - 1) * pageSize},
				{"$limit": pageSize},
			},
		},
	})

	// Fetch cursor
	cur, err := m.coll().Aggregate(ctx, stages)
	if err != nil {
		log.Error().Err(err).Msg("Failed fetching accounts from MongoDB")
		return nil, ErrInternalError
	}
	defer cur.Close(ctx)

	// Load accounts from cursor
	var accounts = make([]AccountWithBalance, 0)
	if cur.Next(ctx) {
		var doc bson.M
		if err := cur.Decode(&doc); err != nil {
			log.Error().Err(err).Interface("doc", doc).Msg("Failed decoding accounts from MongoDB")
			return nil, ErrInternalError
		}
		for _, dataArrayItem := range doc["data"].(bson.A) {
			accDoc := dataArrayItem.(bson.M)
			accounts = append(accounts, AccountWithBalance{
				Account: Account{
					ID:        *util.OptionalObjectIDHex(accDoc["_id"]),
					CreatedOn: util.MustDateTime(accDoc["createdOn"]),
					UpdatedOn: util.OptionalDateTime(accDoc["updatedOn"]),
					Name:      accDoc["name"].(string),
					ParentID:  util.OptionalObjectIDHex(accDoc["parentId"]),
				},
				TotalIncomingAmount: accDoc["incoming"].(float64),
				TotalOutgoingAmount: accDoc["outgoing"].(float64),
				Balance:             accDoc["balance"].(float64),
			})
		}
	}

	// If cursor failed, fail
	if err := cur.Err(); err != nil {
		log.Error().Err(err).Msg("Failed fetching accounts from MongoDB")
		return nil, ErrInternalError
	} else {
		return accounts, nil
	}
}

func (m *manager) PatchAccount(ctx context.Context, id string, name *string, parentID *string) (*Account, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	// Build patch spec
	doc := bson.M{"updatedOn": time.Now()}
	if name != nil {
		doc["name"] = *name
	}
	if parentID != nil {
		doc["parentId"] = util.OptionalObjectID(*parentID)
	}

	// Patch it
	after := options.After
	opts := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	updateSpec := bson.M{"$set": doc}
	if result := m.coll().FindOneAndUpdate(ctx, bson.M{"_id": util.OptionalObjectID(id)}, updateSpec, opts); result.Err() == mongo.ErrNoDocuments {
		return nil, ErrNotFound

	} else if result.Err() != nil {
		log.Error().Err(result.Err()).Interface("name", name).Interface("parentID", parentID).Msg("Failed patching account")
		return nil, ErrInternalError

	} else if err := result.Decode(&doc); err == mongo.ErrNoDocuments {
		return nil, ErrNotFound

	} else if err != nil {
		log.Error().Err(result.Err()).Interface("name", name).Interface("parentID", parentID).Msg("Failed patching account")
		return nil, ErrInternalError

	} else {
		return &Account{
			ID:        *util.OptionalObjectIDHex(doc["_id"]),
			CreatedOn: util.MustDateTime(doc["createdOn"]),
			UpdatedOn: util.OptionalDateTime(doc["updatedOn"]),
			Name:      doc["name"].(string),
			ParentID:  util.OptionalObjectIDHex(doc["parentId"]),
		}, nil
	}
}

func (m *manager) CreateAccount(ctx context.Context, name string, parentID *string) (*Account, error) {
	doc := bson.M{
		"createdOn": time.Now(),
		"updatedOn": nil,
		"name":      name,
		"parentId":  util.OptionalObjectID(parentID),
	}
	if result, err := m.coll().InsertOne(ctx, &doc); err != nil {
		log.Error().Err(err).Msg("Failed persisting new account")
		return nil, ErrInternalError
	} else {
		return &Account{
			ID:        result.InsertedID.(primitive.ObjectID).Hex(),
			CreatedOn: util.MustDateTime(doc["createdOn"]),
			UpdatedOn: nil,
			Name:      name,
			ParentID:  parentID,
		}, nil
	}
}

func (m *manager) UpdateAccount(ctx context.Context, id string, name string, parentID *string) (*Account, error) {
	if id == "" {
		return nil, ErrInvalidID
	}
	t := time.Now()

	// Update spec
	filter := bson.M{"_id": util.OptionalObjectID(id)}
	doc := bson.M{
		"updatedOn": t,
		"name":      name,
		"parentId":  util.OptionalObjectID(parentID),
	}

	// Check if the document exists or not; we're not simply using "FindOneAnd{Update|Replace}" + upsert, because
	// this would create an inconsistency wrt "createdOn" and "updatedOn" when new document would be created.
	after := options.After
	opts := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	if result := m.coll().FindOneAndUpdate(ctx, filter, &bson.M{"$set": doc}, opts); result.Err() == mongo.ErrNoDocuments {

		// Create it
		doc["createdOn"] = t
		doc["updatedOn"] = nil
		if result, err := m.coll().InsertOne(ctx, &doc); err != nil {
			log.Error().Err(err).Msg("Failed persisting new account")
			return nil, ErrInternalError
		} else {
			return &Account{
				ID:        result.InsertedID.(primitive.ObjectID).Hex(),
				CreatedOn: util.MustDateTime(doc["createdOn"]),
				UpdatedOn: nil,
				Name:      name,
				ParentID:  parentID,
			}, nil
		}

	} else if result.Err() != nil {
		log.Error().Err(result.Err()).Msg("Failed looking up & updating new account")
		return nil, ErrInternalError

	} else if err := result.Decode(&doc); err != nil {
		log.Error().Err(err).Msg("Failed reading updated account")
		return nil, ErrInternalError

	} else {
		return &Account{
			ID:        *util.OptionalObjectIDHex(doc["_id"]),
			CreatedOn: util.MustDateTime(doc["createdOn"]),
			UpdatedOn: util.OptionalDateTime(doc["updatedOn"]),
			Name:      doc["name"].(string),
			ParentID:  util.OptionalObjectIDHex(doc["parentId"]),
		}, nil
	}
}
