package rest

import (
	"net/http"
	"time"

	"github.com/golangly/errors"
	"github.com/golangly/webutil"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/bluebudgetz/gate/internal/util"
)

type (
	GetAccountTreeRequest struct{}
	GetAccountTreeData    struct {
		GetAccountData
		IncomingTx float64               `json:"incomingTx" yaml:"incomingTx"`
		OutgoingTx float64               `json:"outgoingTx" yaml:"outgoingTx"`
		Balance    float64               `json:"balance" yaml:"balance"`
		Children   []*GetAccountTreeData `json:"children" yaml:"children"`
	}
	GetAccountTreeResponse struct {
		Accounts []*GetAccountTreeData `json:"data"`
	}
)

func buildTransactionTree(result neo4j.Result) ([]*GetAccountTreeData, error) {
	nodes := make(map[string]*GetAccountTreeData)    // mapping of ID->account
	childIds := make(map[string][]string, 0)         // set of child IDs per account
	nonRoots := make(map[string]byte, 0)             // set of non-root accounts
	roots := make(map[string]*GetAccountTreeData, 0) // set of root accounts
	for result.Next() {
		rec := result.Record()

		// Extract node from record
		var node neo4j.Node
		if v, ok := rec.Get("account"); !ok {
			return nil, errors.New("account missing from result")
		} else if node, ok = v.(neo4j.Node); !ok {
			return nil, errors.New("account mismatch")
		}

		// Create account object, and register it in in id->account map
		account := GetAccountTreeData{
			GetAccountData: GetAccountData{
				ID:        node.Props()["id"].(string),
				CreatedOn: node.Props()["createdOn"].(time.Time),
				UpdatedOn: util.OptionalDateTime(node.Props()["updatedOn"]),
				Name:      node.Props()["name"].(string),
			},
			IncomingTx: node.Props()["incoming"].(float64),
			OutgoingTx: node.Props()["outgoing"].(float64),
			Balance:    node.Props()["balance"].(float64),
		}
		nodes[account.ID] = &account

		// If this account has not yet been encountered *AS A CHILD OF ANOTHER ACCOUNT*, add it as a potential root
		if _, ok := nonRoots[account.ID]; !ok {
			roots[account.ID] = &account
		}

		// Extract this account's children (direct & indirect), mark them as non-roots, remove them from roots list
		if v, ok := rec.Get("children"); !ok {
			return nil, errors.New("children missing from result")
		} else if children, ok := v.([]interface{}); !ok {
			return nil, errors.New("children mismatch")
		} else {
			for _, childId := range children {
				childIds[account.ID] = append(childIds[account.ID], childId.(string))
				nonRoots[childId.(string)] = 1
				delete(roots, childId.(string))
			}
		}
	}
	if err := result.Err(); err != nil {
		return nil, errors.New("iteration failed")
	}

	// Populate each account's "Children" array
	for accountID, ids := range childIds {
		children := make([]*GetAccountTreeData, 0)
		for _, id := range ids {
			children = append(children, nodes[id])
		}
		nodes[accountID].Children = children
	}

	// Map root accounts map to a list of root accounts
	var accounts = make([]*GetAccountTreeData, 0)
	for _, account := range roots {
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func GetAccountTree() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := GetAccountTreeRequest{}
		if err := webutil.Bind(r, &req); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)

		} else if result, err := util.GetNeo4jSession(r.Context()).Run(getAccountsTreeQuery, nil); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, errors.Wrap(err, "query failed"))

		} else if accounts, err := buildTransactionTree(result); err != nil {
			webutil.RenderWithStatusCode(w, r, http.StatusInternalServerError, err)

		} else {
			webutil.RenderWithStatusCode(w, r, http.StatusOK, GetAccountTreeResponse{Accounts: accounts})
		}
	}
}
