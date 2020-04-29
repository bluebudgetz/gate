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

func TestGetAccount(t *testing.T) {
	app, cleanup := testfw.Run(t, func(neo4jDriver neo4j.Driver) func(chi.Router) { return NewRoutes(neo4jDriver) })
	defer cleanup()

	a1 := testfw.NewAccount(app.Neo4jDriver(), "a1", "A1", time.Now(), nil)
	resp := testfw.Request(t, app.BuildURL("localhost", "/accounts/%s", a1.ID), "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data := testfw.ReadResponseBody(t, resp, &GetAccountResponse{}).(*GetAccountResponse)
	assert.Equal(t, a1.ID, data.Account.ID)
	assert.Equal(t, a1.Name, data.Account.Name)
	assert.Equal(t, a1.CreatedOn.UTC(), data.Account.CreatedOn.UTC())
	assert.Nil(t, data.Account.UpdatedOn)
}
