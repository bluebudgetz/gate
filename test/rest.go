package test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/bluebudgetz/gate/internal/app"
)

var client = &http.Client{}

func ServerURL(app *app.App, path string, pathArgs ...interface{}) string {
	return fmt.Sprintf("http://localhost:%d%s", app.GetConfig().HTTP.Port, fmt.Sprintf(path, pathArgs...))
}

func Request(t *testing.T, app *app.App, method string, path string, pathArgs []interface{}, body io.Reader, initializer func(*http.Request)) *http.Response {
	req, err := http.NewRequest(method, ServerURL(app, path, pathArgs...), body)
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
