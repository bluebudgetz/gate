package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bluebudgetz/gate/internal/server/handlers"
	"github.com/bluebudgetz/gate/test"
)

func TestDeleteAccount(t *testing.T) {
	app, cleanup := test.Run(t)
	defer cleanup()

	firstGetResp := test.Request(t, app, "GET", "/accounts/%s", []interface{}{"company"}, nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, firstGetResp.StatusCode)
	acc := test.ResponseBodyObject(t, firstGetResp, &handlers.GetAccountResponse{}).(*handlers.GetAccountResponse)
	assert.Equal(t, "company", acc.Account.ID)

	deleteResp := test.Request(t, app, "DELETE", "/accounts/%s", []interface{}{"company"}, nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	secondGetResp := test.Request(t, app, "GET", "/accounts/%s", []interface{}{"company"}, nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNotFound, secondGetResp.StatusCode)
	assert.Equal(t, int64(0), secondGetResp.ContentLength)
}
