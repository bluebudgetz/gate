package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/docgen"

	"github.com/bluebudgetz/gate/internal/util"
)

const docProjectPath = "github.com/bluebudgetz/gate"
const docIntro = "Bluebudgetz API"

var docOfferedContentTypes = []string{"text/markdown", "application/json", "text/plain"}

func routesDoc(router chi.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch contentType := util.NegotiateContentType(r.Header.Get("Accept"), docOfferedContentTypes, ""); contentType {
		case "text/markdown", "text/plain":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", contentType)
			_, _ = w.Write([]byte(docgen.MarkdownRoutesDoc(router, docgen.MarkdownOpts{ProjectPath: docProjectPath, Intro: docIntro})))
		case "application/json":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", contentType)
			_, _ = w.Write([]byte(docgen.JSONRoutesDoc(router)))
		default:
			w.WriteHeader(http.StatusNotAcceptable)
		}
	}
}
