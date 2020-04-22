package app

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bluebudgetz/gate/internal/server/handlers"
	"github.com/bluebudgetz/gate/internal/util"
)

var client = &http.Client{}

func Request(t *testing.T, url string, method string, body io.Reader, initializer func(*http.Request)) *http.Response {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("request creation failed: %+v", err)
		return nil
	}
	req.Header.Set("accept", "application/json")
	initializer(req)
	if resp, err := client.Do(req); err != nil {
		t.Fatalf("request failed: %+v", err)
		return nil
	} else {
		return resp
	}
}

func ResponseBodyObject(t *testing.T, resp *http.Response, target interface{}) interface{} {
	if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatalf("failed reading response body: %+v", err)
		return nil
	} else if len(bodyBytes) == 0 {
		return nil
	} else if err := json.Unmarshal(bodyBytes, target); err != nil {
		t.Fatalf("failed reading response JSON object: %+v", err)
		return nil
	} else {
		return target
	}
}

func RunTestApp(t *testing.T) (*App, func()) {

	// Create application
	application, cleanup, err := InitializeTestApp()
	if err != nil {
		t.Fatalf("failed creating application: %+v", err)
	}

	// Run application
	quitChan := make(chan error, 10)
	go func() {
		if err := application.Run(quitChan); err != nil {
			t.Fatalf("failed running application: %+v", err)
		}
	}()

	// Sleep to give the app a chance to fully load
	time.Sleep(500 * time.Millisecond)

	// Return application & composite cleanup function
	return application, func() {
		quitChan <- nil
		cleanup()
	}
}

func TestDeleteAccount(t *testing.T) {
	app, cleanup := RunTestApp(t)
	defer cleanup()
	url := app.BuildURL("localhost", "/accounts/%s", "company")

	firstGetResp := Request(t, url, "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, firstGetResp.StatusCode)
	acc := ResponseBodyObject(t, firstGetResp, &handlers.GetAccountResponse{}).(*handlers.GetAccountResponse)
	assert.Equal(t, "company", acc.Account.ID)

	deleteResp := Request(t, url, "DELETE", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	secondGetResp := Request(t, url, "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNotFound, secondGetResp.StatusCode)
	assert.Equal(t, int64(0), secondGetResp.ContentLength)
}

func TestDeleteTx(t *testing.T) {
	app, cleanup := RunTestApp(t)
	defer cleanup()
	url := app.BuildURL("localhost", "/transactions/%s", "salary-2020-01-01")

	firstGetResp := Request(t, url, "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, firstGetResp.StatusCode)
	acc := ResponseBodyObject(t, firstGetResp, &handlers.GetTransactionResponse{}).(*handlers.GetTransactionResponse)
	assert.Equal(t, "salary-2020-01-01", acc.Tx.ID)

	deleteResp := Request(t, url, "DELETE", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)

	secondGetResp := Request(t, url, "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusNotFound, secondGetResp.StatusCode)
	assert.Equal(t, int64(0), secondGetResp.ContentLength)
}

func TestGetAccount(t *testing.T) {
	app, cleanup := RunTestApp(t)
	defer cleanup()
	url := app.BuildURL("localhost", "/accounts/%s", "company")

	resp := Request(t, url, "GET", nil, func(*http.Request) {})
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	acc := ResponseBodyObject(t, resp, &handlers.GetAccountResponse{}).(*handlers.GetAccountResponse)
	assert.Equal(t, "company", acc.Account.ID)
	assert.Equal(t, "Big Company", acc.Account.Name)
	assert.Equal(t, util.MustParseTimeRFC3339("2007-06-05T04:03:02Z"), acc.Account.CreatedOn)
	assert.Nil(t, acc.Account.UpdatedOn)
}
