package handlers

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/require"

	"github.com/bluebudgetz/gate/internal/util"
)

func TestDeleteAccount(t *testing.T) {
	app, cleanup := util.RunTestApp(t, func(neo4jDriver neo4j.Driver) func(chi.Router) { return NewRoutes(neo4jDriver) })
	defer cleanup()
	url := app.BuildURL("localhost", "/accounts/%s", "company")

	firstGetResp := util.Request(t, url, "GET", nil, func(*http.Request) {})
	require.Equal(t, http.StatusOK, firstGetResp.StatusCode)
	acc := util.ResponseBodyObject(t, firstGetResp, &GetAccountResponse{}).(*GetAccountResponse)
	require.Equal(t, "company", acc.Account.ID)

	deleteResp := util.Request(t, url, "DELETE", nil, func(*http.Request) {})
	require.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	secondGetResp := util.Request(t, url, "GET", nil, func(*http.Request) {})
	require.Equal(t, http.StatusNotFound, secondGetResp.StatusCode)
	require.Equal(t, int64(0), secondGetResp.ContentLength)
}
