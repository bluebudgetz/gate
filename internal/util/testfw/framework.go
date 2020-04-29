package testfw

import (
	"time"

	"github.com/golangly/errors"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Account struct {
	ID        string
	Name      string
	CreatedOn time.Time
	Parent    *Account
}

type Transaction struct {
	SourceAccount *Account
	TargetAccount *Account
	ID            string
	IssuedOn      time.Time
	CreatedOn     time.Time
	Amount        float64
	Origin        string
	Comment       string
}

func NewAccount(drv neo4j.Driver, id string, name string, createdOn time.Time, parent *Account) *Account {
	session, err := drv.Session(neo4j.AccessModeWrite)
	if err != nil {
		panic(errors.Wrap(err, "create session error"))
	}
	defer session.Close()

	params := map[string]interface{}{
		"id":        id,
		"name":      name,
		"createdOn": createdOn.UTC(),
	}

	query := "CREATE (:Account {id: $id, name: $name, createdOn: $createdOn})"
	if result, err := session.Run(query, params); err != nil {
		panic(errors.Wrapf(err, "failed populating node '%s'", id))
	} else if err := result.Err(); err != nil {
		panic(errors.Wrapf(err, "failed populating node '%s'", id))
	} else if summary, err := result.Summary(); err != nil {
		panic(errors.Wrapf(err, "failed populating node '%s'", id))
	} else if summary.Counters().NodesCreated() != 1 {
		panic(errors.Newf("failed populating node '%s'", id))
	}
	if parent != nil {
		params := map[string]interface{}{"parentId": parent.ID, "childId": id}
		query := "MERGE (p:Account {id: $parentId}) WITH p MATCH (c:Account {id: $childId}) MERGE (c)-[:childOf]->(p)"
		if result, err := session.Run(query, params); err != nil {
			panic(errors.Wrapf(err, "failed populating node '%s'", id))
		} else if err := result.Err(); err != nil {
			panic(errors.Wrapf(err, "failed populating node '%s'", id))
		} else if _, err := result.Summary(); err != nil {
			panic(errors.Wrapf(err, "failed populating node '%s'", id))
		}
	}
	return &Account{
		ID:        id,
		Name:      name,
		CreatedOn: createdOn,
		Parent:    parent,
	}
}

func (a *Account) CreateTransaction(drv neo4j.Driver, target *Account, id string, issuedOn time.Time, createdOn time.Time, amount float64, comment string) *Transaction {
	session, err := drv.Session(neo4j.AccessModeWrite)
	if err != nil {
		panic(errors.Wrap(err, "create session error"))
	}
	defer session.Close()

	origin := "tests"
	params := map[string]interface{}{
		"srcId":     a.ID,
		"dstId":     target.ID,
		"id":        id,
		"amount":    amount,
		"issuedOn":  issuedOn.UTC(),
		"createdOn": createdOn.UTC(),
		"origin":    origin,
		"comment":   comment,
	}

	query := "" +
		"MATCH (src:Account {id: $srcId}) " +
		"MATCH (dst:Account {id: $dstId}) " +
		"MERGE (src)-[tx:Paid {" +
		"	id: $id, " +
		"	amount: $amount, " +
		"	issuedOn: $issuedOn, " +
		"	createdOn: $createdOn, " +
		"	comment: $comment, " +
		"	origin: $origin" +
		"}]->(dst)" +
		"RETURN tx"
	if result, err := session.Run(query, params); err != nil {
		panic(errors.Wrapf(err, "failed populating transaction '%s'", id))
	} else if err := result.Err(); err != nil {
		panic(errors.Wrapf(err, "failed populating transaction '%s'", id))
	} else if _, err := result.Summary(); err != nil {
		panic(errors.Wrapf(err, "failed populating transaction '%s'", id))
	} else if !result.Next() {
		panic(errors.Newf("failed populating transaction '%s'", id))
	} else if tx, ok := result.Record().Get("tx"); !ok {
		panic(errors.Newf("failed populating transaction '%s'", id))
	} else if txNode, ok := tx.(neo4j.Relationship); !ok {
		panic(errors.Newf("failed populating transaction '%s'", id))
	} else {
		return &Transaction{
			SourceAccount: a,
			TargetAccount: target,
			ID:            txNode.Props()["id"].(string),
			IssuedOn:      txNode.Props()["issuedOn"].(time.Time),
			CreatedOn:     txNode.Props()["createdOn"].(time.Time),
			Amount:        txNode.Props()["amount"].(float64),
			Origin:        txNode.Props()["origin"].(string),
			Comment:       txNode.Props()["comment"].(string),
		}
	}
}
