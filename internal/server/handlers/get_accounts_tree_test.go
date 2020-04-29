package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"

	"github.com/bluebudgetz/gate/internal/util/testfw"
)

func TestGetAccountsTree(t *testing.T) {
	app, cleanup := testfw.Run(t, func(neo4jDriver neo4j.Driver) func(chi.Router) { return NewRoutes(neo4jDriver) })
	defer cleanup()

	assertAccount := func(expected testfw.Account, actual *GetAccountTreeData) {
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.CreatedOn.UTC(), actual.CreatedOn.UTC())
		assert.Nil(t, actual.UpdatedOn)
		assert.Equal(t, float64(0), actual.Balance)
		assert.Equal(t, float64(0), actual.IncomingTx)
		assert.Equal(t, float64(0), actual.OutgoingTx)
	}

	a01 := testfw.NewAccount(app.Neo4jDriver(), "a01", "A01", time.Now(), nil)
	a11 := testfw.NewAccount(app.Neo4jDriver(), "a11", "A11", time.Now(), a01)
	a12 := testfw.NewAccount(app.Neo4jDriver(), "a12", "A12", time.Now(), a01)
	a02 := testfw.NewAccount(app.Neo4jDriver(), "a02", "A02", time.Now(), nil)
	a21 := testfw.NewAccount(app.Neo4jDriver(), "a23", "A21", time.Now(), a02)
	fmt.Printf("%+v\n", a01)
	fmt.Printf("%+v\n", a11)
	fmt.Printf("%+v\n", a12)
	fmt.Printf("%+v\n", a02)
	fmt.Printf("%+v\n", a21)

	resp := testfw.Request(t, app.BuildURL("localhost", "/accounts"), "GET", nil, func(r *http.Request) {})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data := testfw.ReadResponseBody(t, resp, &GetAccountTreeResponse{}).(*GetAccountTreeResponse)
	assert.Equal(t, 2, len(data.Accounts))
	assertAccount(*a01, data.Accounts[0])
	assertAccount(*a11, data.Accounts[0].Children[0])
	assertAccount(*a12, data.Accounts[0].Children[1])
	assertAccount(*a02, data.Accounts[1])
	assertAccount(*a21, data.Accounts[1].Children[0])
}
