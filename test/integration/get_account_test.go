package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bluebudgetz/gate/internal/server/handlers"
	"github.com/bluebudgetz/gate/internal/util"
	"github.com/bluebudgetz/gate/test"
)

func TestGetAccount(t *testing.T) {
	app, cleanup := test.Run(t)
	defer cleanup()

	resp := test.Request(t, app, "GET", "/accounts/%s", []interface{}{"company"}, nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	acc := test.ResponseBodyObject(t, resp, &handlers.GetAccountResponse{}).(*handlers.GetAccountResponse)
	assert.Equal(t, "company", acc.Account.ID)
	assert.Equal(t, "Big Company", acc.Account.Name)
	assert.Equal(t, util.MustParseTimeRFC3339("2007-06-05T04:03:02Z"), acc.Account.CreatedOn)
	assert.Nil(t, acc.Account.UpdatedOn)
}
