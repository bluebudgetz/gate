package handlers

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"

	"github.com/bluebudgetz/gate/internal/util"
)

func TestGetAccount(t *testing.T) {
	app, cleanup := util.RunTestApp(t, func(neo4jDriver neo4j.Driver) func(chi.Router) { return NewRoutes(neo4jDriver) })
	defer cleanup()
	url := app.BuildURL("localhost", "/accounts/%s", "company")

	resp := util.Request(t, url, "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	acc := util.ResponseBodyObject(t, resp, &GetAccountResponse{}).(*GetAccountResponse)
	assert.Equal(t, "company", acc.Account.ID)
	assert.Equal(t, "Big Company", acc.Account.Name)
	assert.Equal(t, util.MustParseTimeRFC3339("2007-06-05T04:03:02Z"), acc.Account.CreatedOn)
	assert.Nil(t, acc.Account.UpdatedOn)
}
