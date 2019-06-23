package accounts

import (
	"encoding/json"
	"github.com/bluebudgetz/gate/internal/util"
	mongoutil "github.com/bluebudgetz/gate/internal/util/mongo"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

func (module *Module) getAccountsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cur, err := module.mongo.Database("bluebudgetz").Collection("accounts").Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed fetching accounts from MongoDB")
		return
	}
	defer cur.Close(ctx)

	list := make([]AccountDocument, 0)
	for cur.Next(ctx) {
		var doc AccountDocument
		if err := cur.Decode(&doc); err != nil {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(err).Msg("Failed decoding accounts from MongoDB")
			return
		}
		list = append(list, doc)
	}
	if err := cur.Err(); err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed decoding accounts from MongoDB")
		return
	}

	// Write response back to client
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if module.config.Environment != util.EnvProduction {
		encoder.SetIndent("", "  ")
	}
	w.WriteHeader(http.StatusOK)
	if err := encoder.Encode(list); err != nil {
		log.Error().Err(err).Msg("Failed encoding accounts to client")
		return
	}
}

func (module *Module) getAccountsTree(w http.ResponseWriter, r *http.Request) {
	type AccountTreeDTO struct {
		ID        *primitive.ObjectID `json:"id,omitempty"`
		CreatedOn *time.Time          `json:"createdOn,omitempty"`
		UpdatedOn *time.Time          `json:"updatedOn,omitempty"`
		Name      *string             `json:"name,omitempty"`
		Children  *[]*AccountTreeDTO  `json:"children,omitempty"`
	}

	cur, err := module.mongo.Database("bluebudgetz").Collection("accounts").Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed fetching accounts from MongoDB")
		return
	}
	defer cur.Close(r.Context())

	// fetch all account documents, and mark which ones are child accounts
	accounts := make(map[string]*AccountTreeDTO, 100)
	parents := make(map[string]string, 100)
	roots := make([]*AccountTreeDTO, 0)
	for cur.Next(r.Context()) {
		var account AccountDocument
		if err := cur.Decode(&account); err != nil {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(err).Msg("Failed decoding accounts from MongoDB")
			return
		}

		children := make([]*AccountTreeDTO, 0)
		accounts[account.ID.Hex()] = &AccountTreeDTO{
			ID:        account.ID,
			CreatedOn: account.CreatedOn,
			UpdatedOn: account.UpdatedOn,
			Name:      account.Name,
			Children:  &children,
		}
		if account.ParentID != nil {
			parents[account.ID.Hex()] = account.ParentID.Hex()
		} else {
			roots = append(roots, accounts[account.ID.Hex()])
		}
	}
	if err := cur.Err(); err != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed decoding accounts from MongoDB")
		return
	}

	// find all roots, and populate "Children" arrays
	for _, account := range accounts {
		if parentId, ok := parents[account.ID.Hex()]; ok {
			children := append(*accounts[parentId].Children, account)
			accounts[parentId].Children = &children
		}
	}

	// Write response back to client
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if module.config.Environment != util.EnvProduction {
		encoder.SetIndent("", "  ")
	}
	w.WriteHeader(http.StatusOK)
	if err := encoder.Encode(roots); err != nil {
		log.Error().Err(err).Msg("Failed encoding accounts to client")
		return
	}
}

func (module *Module) getAccount(w http.ResponseWriter, r *http.Request) {
	ID := mongoutil.ObjectIdFromNillableString(chi.URLParam(r, "id"))
	if ID == nil {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	result := module.mongo.Database("bluebudgetz").Collection("accounts").FindOne(r.Context(), bson.M{"_id": ID})
	if result.Err() != nil {
		http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
		log.Error().Err(result.Err()).Str("id", ID.Hex()).Msg("Failed fetching account from MongoDB")
		return
	}

	var account AccountDocument
	if err := result.Decode(&account); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Account could not be found.", http.StatusNotFound)
			return
		} else {
			http.Error(w, "Internal error occurred.", http.StatusInternalServerError)
			log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed fetching account from MongoDB")
			return
		}
	}

	// Write response back to client
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if module.config.Environment != util.EnvProduction {
		encoder.SetIndent("", "  ")
	}
	w.WriteHeader(http.StatusOK)
	if err := encoder.Encode(account); err != nil {
		log.Error().Err(err).Str("id", ID.Hex()).Msg("Failed encoding account to client")
		return
	}
}
