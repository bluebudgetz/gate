package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bluebudgetz/gate/internal/server/handlers"
	"github.com/bluebudgetz/gate/test"
)

func TestDeleteTx(t *testing.T) {
	app, cleanup := test.Run(t)
	defer cleanup()

	firstGetResp := test.Request(t, app, "GET", "/transactions/%s", []interface{}{"salary-2020-01-01"}, nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, firstGetResp.StatusCode)
	acc := test.ResponseBodyObject(t, firstGetResp, &handlers.GetTransactionResponse{}).(*handlers.GetTransactionResponse)
	assert.Equal(t, "salary-2020-01-01", acc.Tx.ID)

	deleteResp := test.Request(t, app, "DELETE", "/transactions/%s", []interface{}{"salary-2020-01-01"}, nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	secondGetResp := test.Request(t, app, "GET", "/transactions/%s", []interface{}{"salary-2020-01-01"}, nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNotFound, secondGetResp.StatusCode)
	assert.Equal(t, int64(0), secondGetResp.ContentLength)
}
