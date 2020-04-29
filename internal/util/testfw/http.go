package testfw

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
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

func ReadResponseBody(t *testing.T, resp *http.Response, target interface{}) interface{} {
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
