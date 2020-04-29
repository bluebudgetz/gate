package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/require"

	"github.com/bluebudgetz/gate/internal/util/testfw"
)

func TestDeleteAccount(t *testing.T) {
	app, cleanup := testfw.Run(t, func(neo4jDriver neo4j.Driver) func(chi.Router) { return NewRoutes(neo4jDriver) })
	defer cleanup()

	a1 := testfw.NewAccount(app.Neo4jDriver(), "a1", "A1", time.Now(), nil)
	a1GetResp := testfw.Request(t, app.BuildURL("localhost", "/accounts/%s", a1.ID), "GET", nil, func(*http.Request) {})
	require.Equal(t, http.StatusOK, a1GetResp.StatusCode)
	data := testfw.ReadResponseBody(t, a1GetResp, &GetAccountResponse{}).(*GetAccountResponse)
	require.Equal(t, a1.ID, data.Account.ID)

	a1DeleteResp := testfw.Request(t, app.BuildURL("localhost", "/accounts/%s", a1.ID), "DELETE", nil, func(*http.Request) {})
	require.Equal(t, http.StatusNoContent, a1DeleteResp.StatusCode)

	a1GetResp2 := testfw.Request(t, app.BuildURL("localhost", "/accounts/%s", a1.ID), "GET", nil, func(*http.Request) {})
	require.Equal(t, http.StatusNotFound, a1GetResp2.StatusCode)
}
