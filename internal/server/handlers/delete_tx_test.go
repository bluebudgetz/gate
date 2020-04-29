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

func TestDeleteTx(t *testing.T) {
	app, cleanup := testfw.Run(t, func(neo4jDriver neo4j.Driver) func(chi.Router) { return NewRoutes(neo4jDriver) })
	defer cleanup()

	a1 := testfw.NewAccount(app.Neo4jDriver(), "a1", "A1", time.Now(), nil)
	a2 := testfw.NewAccount(app.Neo4jDriver(), "a2", "A2", time.Now(), nil)
	t1 := a1.CreateTransaction(app.Neo4jDriver(), a2, "a1-to-a2", time.Now(), time.Now(), 19, "Testing transaction")

	firstGetResp := testfw.Request(t, app.BuildURL("localhost", "/transactions/%s", t1.ID), "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, firstGetResp.StatusCode)

	deleteResp := testfw.Request(t, app.BuildURL("localhost", "/transactions/%s", t1.ID), "DELETE", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	secondGetResp := testfw.Request(t, app.BuildURL("localhost", "/transactions/%s", t1.ID), "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNotFound, secondGetResp.StatusCode)
	assert.Equal(t, int64(0), secondGetResp.ContentLength)
}
