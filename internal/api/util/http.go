package util

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func Respond(w http.ResponseWriter, r *http.Request, statusCode int, value interface{}) {
	if value == nil {
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	} else {
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		if err := enc.Encode(value); err != nil {
			panic(errors.Wrap(err, "failed marshalling JSON"))
		}
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if _, err := w.Write(buf.Bytes()); err != nil {
			panic(errors.Wrap(err, "failed writing JSON"))
		}
	}
}

func ExternalURL(r *http.Request) string {
	// TODO: support "x-forwarded-proto" value for the scheme here
	return "http://" + r.Host;
}