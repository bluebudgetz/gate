package transactions

import (
	"context"
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
		CreateTransaction(ctx context.Context, issuedOn time.Time, origin string, sourceAccountID string,
			targetAccountID string, amount float64, comment string) (*Transaction, error)
		DeleteTransaction(ctx context.Context, id string) error
		GetTransaction(ctx context.Context, id string) (*Transaction, error)
		GetTransactionsList(ctx context.Context, page uint, pageSize uint) ([]Transaction, error)
		PatchTransaction(ctx context.Context, id string, issuedOn *time.Time, origin *string,
			sourceAccountID *string, targetAccountID *string, amount *float64, comment *string) (*Transaction, error)
		UpdateTransaction(ctx context.Context, id string, issuedOn time.Time, origin string,
			sourceAccountID string, targetAccountID string, amount float64, comment string) (*Transaction, error)
	}

	manager struct {
		mongoClient *mongo.Client
	}

	Transaction struct {
		ID              string     `json:"id" yaml:"id"`
		CreatedOn       time.Time  `json:"createdOn" yaml:"createdOn"`
		UpdatedOn       *time.Time `json:"updatedOn" yaml:"updatedOn"`
		IssuedOn        time.Time  `json:"issuedOn" yaml:"issuedOn"`
		Origin          string     `json:"origin" yaml:"origin"`
		SourceAccountID string     `json:"sourceAccountId" yaml:"sourceAccountId"`
		TargetAccountID string     `json:"targetAccountId" yaml:"targetAccountId"`
		Amount          float64    `json:"amount" yaml:"amount"`
		Comment         string     `json:"comment" yaml:"comment"`
	}
)

var (
	ErrInvalidID     = fmt.Errorf("invalid transaction ID")
	ErrNotFound      = fmt.Errorf("transaction not found")
	ErrInternalError = fmt.Errorf("internal error")
)

func init() {
}

func NewManager(mongoClient *mongo.Client) Manager {
	return &manager{mongoClient}
}

func (m *manager) coll() *mongo.Collection {
	return m.mongoClient.Database("bluebudgetz").Collection("transactions")
}

func (m *manager) CreateTransaction(ctx context.Context, issuedOn time.Time, origin string, sourceAccountID string,
	targetAccountID string, amount float64, comment string) (*Transaction, error) {

	doc := bson.M{
		"createdOn":       time.Now(),
		"updatedOn":       nil,
		"issuedOn":        issuedOn,
		"origin":          origin,
		"sourceAccountId": util.MustObjectID(sourceAccountID),
		"targetAccountId": util.MustObjectID(targetAccountID),
		"amount":          amount,
		"comment":         comment,
	}
	if result, err := m.coll().InsertOne(ctx, &doc); err != nil {
		log.Error().Err(err).Msg("Failed persisting new transaction")
		return nil, ErrInternalError
	} else {
		return &Transaction{
			ID:              result.InsertedID.(primitive.ObjectID).Hex(),
			CreatedOn:       util.MustDateTime(doc["createdOn"]),
			UpdatedOn:       nil,
			IssuedOn:        util.MustDateTime(doc["issuedOn"]),
			Origin:          origin,
			SourceAccountID: sourceAccountID,
			TargetAccountID: targetAccountID,
			Amount:          amount,
			Comment:         comment,
		}, nil
	}
}

func (m *manager) DeleteTransaction(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidID

	} else if result, err := m.coll().DeleteOne(ctx, bson.M{"_id": util.MustObjectID(id)}); err == mongo.ErrNoDocuments {
		return ErrNotFound

	} else if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed deleting transaction")
		return ErrInternalError

	} else if result != nil && result.DeletedCount == 0 {
		return ErrNotFound

	} else {
		return nil
	}
}

func (m *manager) GetTransaction(ctx context.Context, id string) (*Transaction, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	doc := bson.M{}

	if result := m.coll().FindOne(ctx, bson.M{"_id": util.MustObjectID(id)}); result.Err() == mongo.ErrNoDocuments {
		return nil, ErrNotFound

	} else if result.Err() != nil {
		log.Error().Err(result.Err()).Str("id", id).Msg("Failed looking up transaction")
		return nil, ErrInternalError

	} else if err := result.Decode(&doc); err == mongo.ErrNoDocuments {
		return nil, ErrNotFound

	} else if err != nil {
		log.Error().Err(result.Err()).Str("id", id).Msg("Failed looking up transaction")
		return nil, ErrInternalError

	} else {
		return &Transaction{
			ID:              util.MustObjectIDHex(doc["_id"]),
			CreatedOn:       util.MustDateTime(doc["createdOn"]),
			UpdatedOn:       util.OptionalDateTime(doc["updatedOn"]),
			IssuedOn:        util.MustDateTime(doc["issuedOn"]),
			Origin:          doc["origin"].(string),
			SourceAccountID: doc["sourceAccountId"].(primitive.ObjectID).Hex(),
			TargetAccountID: doc["targetAccountId"].(primitive.ObjectID).Hex(),
			Amount:          doc["amount"].(float64),
			Comment:         doc["comment"].(string),
		}, nil
	}
}

func (m *manager) GetTransactionsList(ctx context.Context, page uint, pageSize uint) ([]Transaction, error) {

	// Fetch cursor
	limit := int64(pageSize)
	skip := int64((page - 1) * pageSize)
	opts := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  bson.M{"issuedOn": 1},
	}
	cur, err := m.coll().Find(ctx, bson.M{}, &opts)
	if err != nil {
		log.Error().Err(err).Msg("Failed fetching transactions from MongoDB")
		return nil, ErrInternalError
	}
	defer cur.Close(ctx)

	// Load transactions from cursor
	var transactions = make([]Transaction, 0)
	for cur.Next(ctx) {
		var doc bson.M
		if err := cur.Decode(&doc); err != nil {
			return nil, ErrInternalError
		}
		transactions = append(transactions, Transaction{
			ID:              util.MustObjectIDHex(doc["_id"]),
			CreatedOn:       util.MustDateTime(doc["createdOn"]),
			UpdatedOn:       util.OptionalDateTime(doc["updatedOn"]),
			IssuedOn:        util.MustDateTime(doc["issuedOn"]),
			Origin:          doc["origin"].(string),
			SourceAccountID: doc["sourceAccountId"].(primitive.ObjectID).Hex(),
			TargetAccountID: doc["targetAccountId"].(primitive.ObjectID).Hex(),
			Amount:          doc["amount"].(float64),
			Comment:         doc["comment"].(string),
		})
	}

	// If cursor failed, fail
	if err := cur.Err(); err != nil {
		log.Error().Err(err).Msg("Failed fetching transactions from MongoDB")
		return nil, ErrInternalError
	} else {
		return transactions, nil
	}
}

func (m *manager) PatchTransaction(ctx context.Context, id string, issuedOn *time.Time, origin *string,
	sourceAccountID *string, targetAccountID *string, amount *float64, comment *string) (*Transaction, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	// Build patch spec
	doc := bson.M{"updatedOn": time.Now()}
	if issuedOn != nil {
		doc["issuedOn"] = *issuedOn
	}
	if origin != nil {
		doc["origin"] = *origin
	}
	if sourceAccountID != nil {
		doc["sourceAccountId"] = util.MustObjectID(*sourceAccountID)
	}
	if targetAccountID != nil {
		doc["targetAccountID"] = util.MustObjectID(*targetAccountID)
	}
	if amount != nil {
		doc["amount"] = *amount
	}
	if comment != nil {
		doc["comment"] = *comment
	}

	// Patch it
	filter := bson.M{"_id": util.MustObjectID(id)}
	after := options.After
	opts := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	updateSpec := bson.M{"$set": doc}
	if result := m.coll().FindOneAndUpdate(ctx, filter, updateSpec, opts); result.Err() == mongo.ErrNoDocuments {
		return nil, ErrNotFound

	} else if result.Err() != nil {
		log.Error().Err(result.Err()).Msg("Failed patching transaction")
		return nil, ErrInternalError

	} else if err := result.Decode(&doc); err == mongo.ErrNoDocuments {
		return nil, ErrNotFound

	} else if err != nil {
		log.Error().Err(result.Err()).Msg("Failed patching account")
		return nil, ErrInternalError

	} else {
		return &Transaction{
			ID:              util.MustObjectIDHex(doc["_id"]),
			CreatedOn:       util.MustDateTime(doc["createdOn"]),
			UpdatedOn:       util.OptionalDateTime(doc["updatedOn"]),
			IssuedOn:        util.MustDateTime(doc["issuedOn"]),
			Origin:          doc["origin"].(string),
			SourceAccountID: util.MustObjectIDHex(doc["sourceAccountId"]),
			TargetAccountID: util.MustObjectIDHex(doc["targetAccountId"]),
			Amount:          doc["amount"].(float64),
			Comment:         doc["comment"].(string),
		}, nil
	}
}

func (m *manager) UpdateTransaction(ctx context.Context, id string, issuedOn time.Time, origin string,
	sourceAccountID string, targetAccountID string, amount float64, comment string) (*Transaction, error) {

	if id == "" {
		return nil, ErrInvalidID
	}
	t := time.Now()

	// Update spec
	filter := bson.M{"_id": util.MustObjectID(id)}
	doc := bson.M{
		"updatedOn":       t,
		"issuedOn":        issuedOn,
		"origin":          origin,
		"sourceAccountId": util.MustObjectID(sourceAccountID),
		"targetAccountId": util.MustObjectID(targetAccountID),
		"amount":          amount,
		"comment":         comment,
	}

	// Check if the document exists or not; we're not simply using "FindOneAnd{Update|Replace}" + upsert, because
	// this would create an inconsistency wrt "createdOn" and "updatedOn" when new document would be created.
	after := options.After
	opts := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	if result := m.coll().FindOneAndUpdate(ctx, filter, &bson.M{"$set": doc}, opts); result.Err() == mongo.ErrNoDocuments {

		// Create it
		doc["_id"] = util.MustObjectID(id)
		doc["createdOn"] = t
		doc["updatedOn"] = nil
		if _, err := m.coll().InsertOne(ctx, &doc); err != nil {
			log.Error().Err(err).Msg("Failed persisting new transaction")
			return nil, ErrInternalError
		} else {
			return &Transaction{
				ID:              util.MustObjectID(id).Hex(),
				CreatedOn:       util.MustDateTime(doc["createdOn"]),
				UpdatedOn:       util.OptionalDateTime(doc["updatedOn"]),
				IssuedOn:        util.MustDateTime(doc["issuedOn"]),
				Origin:          doc["origin"].(string),
				SourceAccountID: util.MustObjectIDHex(doc["sourceAccountId"]),
				TargetAccountID: util.MustObjectIDHex(doc["targetAccountId"]),
				Amount:          doc["amount"].(float64),
				Comment:         doc["comment"].(string),
			}, nil
		}

	} else if result.Err() != nil {
		log.Error().Err(result.Err()).Msg("Failed looking up & updating new transaction")
		return nil, ErrInternalError

	} else if err := result.Decode(&doc); err != nil {
		log.Error().Err(err).Msg("Failed reading updated transaction")
		return nil, ErrInternalError

	} else {
		return &Transaction{
			ID:              util.MustObjectIDHex(doc["_id"]),
			CreatedOn:       util.MustDateTime(doc["createdOn"]),
			UpdatedOn:       util.OptionalDateTime(doc["updatedOn"]),
			IssuedOn:        util.MustDateTime(doc["issuedOn"]),
			Origin:          doc["origin"].(string),
			SourceAccountID: util.MustObjectIDHex(doc["sourceAccountId"]),
			TargetAccountID: util.MustObjectIDHex(doc["targetAccountId"]),
			Amount:          doc["amount"].(float64),
			Comment:         doc["comment"].(string),
		}, nil
	}
}
