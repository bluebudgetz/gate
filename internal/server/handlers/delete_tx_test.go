package handlers

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"

	"github.com/bluebudgetz/gate/internal/util"
)

func TestDeleteTx(t *testing.T) {
	app, cleanup := util.RunTestApp(t, func(neo4jDriver neo4j.Driver) func(chi.Router) { return NewRoutes(neo4jDriver) })
	defer cleanup()
	url := app.BuildURL("localhost", "/transactions/%s", "salary-2020-01-01")

	firstGetResp := util.Request(t, url, "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, firstGetResp.StatusCode)
	acc := util.ResponseBodyObject(t, firstGetResp, &GetTransactionResponse{}).(*GetTransactionResponse)
	assert.Equal(t, "salary-2020-01-01", acc.Tx.ID)

	deleteResp := util.Request(t, url, "DELETE", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	secondGetResp := util.Request(t, url, "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNotFound, secondGetResp.StatusCode)
	assert.Equal(t, int64(0), secondGetResp.ContentLength)
}
