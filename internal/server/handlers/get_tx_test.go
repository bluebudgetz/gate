package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"

	"github.com/bluebudgetz/gate/internal/util/testfw"
)

func TestGetTx(t *testing.T) {
	app, cleanup := testfw.Run(t, func(neo4jDriver neo4j.Driver) func(chi.Router) { return NewRoutes(neo4jDriver) })
	defer cleanup()

	a1 := testfw.NewAccount(app.Neo4jDriver(), "a1", "A1", time.Now(), nil)
	a2 := testfw.NewAccount(app.Neo4jDriver(), "a2", "A2", time.Now(), nil)
	t1 := a1.CreateTransaction(app.Neo4jDriver(), a2, "a1-to-a2", time.Now(), time.Now(), 19, "Testing transaction")

	firstGetResp := testfw.Request(t, app.BuildURL("localhost", "/transactions/%s", t1.ID), "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, firstGetResp.StatusCode)
	data := testfw.ReadResponseBody(t, firstGetResp, &GetTransactionResponse{}).(*GetTransactionResponse)
	assert.Equal(t, t1.ID, data.Tx.ID)
	assert.Equal(t, t1.CreatedOn, data.Tx.CreatedOn)
	assert.Equal(t, t1.Comment, data.Tx.Comment)
	assert.Equal(t, t1.Amount, data.Tx.Amount)
	assert.Equal(t, t1.Origin, data.Tx.Origin)
	assert.Equal(t, t1.IssuedOn, data.Tx.IssuedOn)
	assert.Nil(t, data.Tx.UpdatedOn)
	assert.Equal(t, a1.ID, data.Tx.SourceAccountID)
	assert.Equal(t, a2.ID, data.Tx.TargetAccountID)
}
