package main

import (
	"context"
	"encoding/json"
	mongoutil "github.com/bluebudgetz/gate/internal/util/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"log"
	"os"
)

var (
	stdout = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	stderr = log.New(os.Stderr, "", log.Ldate|log.Ltime)
)

func main() {
	ctx := context.Background()

	mongo, err := mongoutil.CreateMongoClient("mongodb://localhost:27017")
	if err != nil {
		stderr.Fatalf("Failed connecting to MongoDB %v\n", err)
	}

	var aggDoc bson.A
	if queryBytes, err := ioutil.ReadAll(os.Stdin); err != nil {
		stderr.Fatalf("Failed reading query from stdin %v\n", err)
	} else if err := json.Unmarshal(queryBytes, &aggDoc); err != nil {
		stderr.Fatalf("Failed decoding query JSON %v\n", err)
	}

	cur, err := mongo.Database("bluebudgetz").Collection(os.Args[1]).Aggregate(ctx, aggDoc)
	if err != nil {
		stderr.Fatalf("Failed executing aggregation %v\n", err)
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var doc bson.M
		if err := cur.Decode(&doc); err != nil {
			stderr.Fatalf("Failed decoding results %v\n", err)
		}

		if jsonBytes, err := json.MarshalIndent(doc, "", "  "); err != nil {
			stderr.Fatalf("Failed marshalling JSON %v\n", err)
		} else if _, err := stdout.Writer().Write(jsonBytes); err != nil {
			stderr.Fatalf("Failed writing JSON %v\n", err)
		}
	}
	if err := cur.Err(); err != nil {
		stderr.Fatalf("Failed decoding results %v\n", err)
	}
}
