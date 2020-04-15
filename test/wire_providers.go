package test

import (
	"time"

	"github.com/golangly/errors"
	"github.com/golangly/log"
	"github.com/jessevdk/go-flags"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/bluebudgetz/gate/internal/config"
	"github.com/bluebudgetz/gate/internal/services"
)

//go:generate go-bindata -o ./assets_gen.go -ignore ".*\\.go" -pkg test ./...

func NewConfig() (config.Config, error) {
	cfg := config.Config{}
	parser := flags.NewParser(&cfg, flags.HelpFlag|flags.PrintErrors|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"
	parser.LongDescription = "Bluebudgetz API gateway. This is the API micro-service centralizing Bluebudgetz APIs."
	if _, err := parser.ParseArgs([]string{}); err != nil {
		return config.Config{}, err
	}
	cfg.HTTP.DisableAccessLog = true
	return cfg, nil
}

type testTx struct {
	srcID string
	dstID string
	id string
	amount float64
	issuedOn time.Time
	comment string
}

var txList = []testTx{
	{
		srcID:    "company",
		dstID:    "bankAccount",
		id:       "salary-2020-01",
		amount:   10000,
		issuedOn: time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
		comment:  "Salary for January 2020",
	},
}

func NewNeo4jDriver() (neo4j.Driver, func(), error) {
	driver, cleanup, err := services.NewNeo4jDriver()
	if err != nil {
		return nil, cleanup, err
	}

	session, err := services.CreateNeo4jSession(driver, neo4j.AccessModeWrite)
	if err != nil {
		return nil, cleanup, err
	}
	defer session.Close()

	// Cleanup database
	log.Info("Cleaning Neo4j database")
	if result, err := session.Run(string(MustAsset("neo4j-cleanup.cyp")), nil); err != nil {
		return nil, cleanup, errors.Wrap(err, "failed cleaning database")
	} else if err := result.Err(); err != nil {
		return nil, cleanup, errors.Wrap(err, "result error")
	}

	// Create nodes
	log.Info("Creating test nodes")
	if result, err := session.Run(string(MustAsset("neo4j-create-nodes.cyp")), nil); err != nil {
		return nil, cleanup, errors.Wrap(err, "failed populating nodes")
	} else if err := result.Err(); err != nil {
		return nil, cleanup, errors.Wrap(err, "result error")
	} else if summary, err := result.Summary(); err != nil {
		return nil, cleanup, errors.Wrap(err, "result summary error")
	} else if summary.Counters().NodesCreated() > 0 {
		return nil, cleanup, errors.Wrap(err, "nodes should have been created")
	} else if summary.Counters().RelationshipsCreated() <= 0 {
		return nil, cleanup, errors.Wrap(err, "relationships should not have been created")
	}

	// Populate new data
	log.Info("Creating transactions")
	for _, tx := range txList {
		if result, err := session.Run(string(MustAsset("neo4j-create-tx.cyp")), map[string]interface{}{
			"srcId": tx.srcID,
			"dstId": tx.dstID,
			"id": tx.id,
			"amount": tx.amount,
			"issuedOn": tx.issuedOn,
			"comment": tx.comment,
		}); err != nil {
			return nil, cleanup, errors.Wrap(err, "failed populating transactions")
		} else if err := result.Err(); err != nil {
			return nil, cleanup, errors.Wrap(err, "result error")
		} else if summary, err := result.Summary(); err != nil {
			return nil, cleanup, errors.Wrap(err, "result summary error")
		} else if summary.Counters().NodesCreated() > 0 {
			return nil, cleanup, errors.Wrap(err, "nodes should not have been created")
		} else if summary.Counters().RelationshipsCreated() <= 0 {
			return nil, cleanup, errors.Wrap(err, "no relationships created")
		}
	}

	return driver, cleanup, err
}
