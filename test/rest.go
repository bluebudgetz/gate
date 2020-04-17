package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bluebudgetz/gate/internal/app"
)

var client = &http.Client{}

func GetAccount(t *testing.T, app *app.App, id string) map[string]interface{} {
	url := fmt.Sprintf("http://localhost:%d/accounts/%s", app.GetConfig().HTTP.Port, id)

	var account interface{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("request creation failed: %+v", err)
	}
	req.Header.Set("accept", "application/json")

	if resp, err := client.Do(req); err != nil {
		t.Fatalf("request sending failed: %+v", err)
	} else if resp.StatusCode == http.StatusNotFound {
		return nil
	} else if resp.StatusCode != http.StatusOK {
		t.Fatalf("wrong status code: got %d, expected %d", resp.StatusCode, http.StatusOK)
	} else if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatalf("failed reading response body: %+v", err)
	} else if err := json.Unmarshal(bodyBytes, &account); err != nil {
		t.Fatalf("failed reading response JSON: %+v", err)
	}
	return account.(map[string]interface{})
}

func MustAccount(t *testing.T, app *app.App, id string) map[string]interface{} {
	account := GetAccount(t, app, id)
	if account == nil {
		t.Fatalf("account '%s' not found", id)
	}
	return account
}

func DeleteAccount(t *testing.T, app *app.App, id string) {
	url := fmt.Sprintf("http://localhost:%d/accounts/%s", app.GetConfig().HTTP.Port, id)
	if req, err := http.NewRequest("DELETE", url, nil); err != nil {
		t.Fatalf("request creation failed: %+v", err)
	} else if resp, err := client.Do(req); err != nil {
		t.Fatalf("request sending failed: %+v", err)
	} else if resp.StatusCode == http.StatusNotFound {
		t.Fatalf("account '%s' not found", id)
	} else if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("wrong status code: got %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}

func GetTx(t *testing.T, app *app.App, id string) map[string]interface{} {
	url := fmt.Sprintf("http://localhost:%d/transactions/%s", app.GetConfig().HTTP.Port, id)

	var tx interface{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("request creation failed: %+v", err)
	}
	req.Header.Set("accept", "application/json")

	if resp, err := client.Do(req); err != nil {
		t.Fatalf("request sending failed: %+v", err)
	} else if resp.StatusCode == http.StatusNotFound {
		return nil
	} else if resp.StatusCode != http.StatusOK {
		t.Fatalf("wrong status code: got %d, expected %d", resp.StatusCode, http.StatusOK)
	} else if bodyBytes, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fatalf("failed reading response body: %+v", err)
	} else if err := json.Unmarshal(bodyBytes, &tx); err != nil {
		t.Fatalf("failed reading response JSON: %+v", err)
	}
	return tx.(map[string]interface{})
}

func MustTx(t *testing.T, app *app.App, id string) map[string]interface{} {
	tx := GetTx(t, app, id)
	if tx == nil {
		t.Fatalf("transaction '%s' not found", id)
	}
	return tx
}

func DeleteTx(t *testing.T, app *app.App, id string) {
	url := fmt.Sprintf("http://localhost:%d/transactions/%s", app.GetConfig().HTTP.Port, id)
	if req, err := http.NewRequest("DELETE", url, nil); err != nil {
		t.Fatalf("request creation failed: %+v", err)
	} else if resp, err := client.Do(req); err != nil {
		t.Fatalf("request sending failed: %+v", err)
	} else if resp.StatusCode == http.StatusNotFound {
		t.Fatalf("transaction '%s' not found", id)
	} else if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("wrong status code: got %d, expected %d", resp.StatusCode, http.StatusOK)
	}
}
