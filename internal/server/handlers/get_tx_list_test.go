package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"

	"github.com/bluebudgetz/gate/internal/util/testfw"
)

func TestGetTxList(t *testing.T) {
	app, cleanup := testfw.Run(t, func(neo4jDriver neo4j.Driver) func(chi.Router) { return NewRoutes(neo4jDriver) })
	defer cleanup()

	assertTx := func(expected testfw.Transaction, actual *GetTransactionListItemData) {
		assert.Equal(t, expected.SourceAccount.ID, actual.SourceAccountID)
		assert.Equal(t, expected.TargetAccount.ID, actual.TargetAccountID)
		assert.Equal(t, expected.ID, actual.ID)
		assert.Equal(t, expected.IssuedOn, actual.IssuedOn)
		assert.Equal(t, expected.CreatedOn, actual.CreatedOn)
		assert.Equal(t, expected.Amount, actual.Amount)
		assert.Equal(t, expected.Origin, actual.Origin)
		assert.Equal(t, expected.Comment, actual.Comment)
		assert.Nil(t, actual.UpdatedOn)
	}

	var accounts = make([]*testfw.Account, 0)
	for i := 0; i < 10; i++ {
		id := "a" + strconv.FormatInt(int64(i), 10)
		name := "A" + strconv.FormatInt(int64(i), 10)
		accounts = append(accounts, testfw.NewAccount(app.Neo4jDriver(), id, name, time.Now(), nil))
	}
	var txs = make([]*testfw.Transaction, 0)
	for i := 0; i < 9; i++ {
		src := accounts[i]
		dst := accounts[i+1]
		id := fmt.Sprintf("%s-to-%s", src.ID, dst.ID)
		amount := float64(i * 100)
		comment := fmt.Sprintf("Transfer %f from %s to %s", amount, src.ID, dst.ID)
		txs = append(txs, src.CreateTransaction(app.Neo4jDriver(), dst, id, time.Now().Add(time.Duration(-i)*time.Hour), time.Now(), amount, comment))
	}

	resp := testfw.Request(t, app.BuildURL("localhost", "/transactions"), "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data := testfw.ReadResponseBody(t, resp, &GetTransactionListResponse{}).(*GetTransactionListResponse)
	assert.Equal(t, len(txs), len(data.Transactions))
	for i := 0; i < 9; i++ {
		assertTx(*txs[i], &data.Transactions[i])
	}
}
