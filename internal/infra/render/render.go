package render

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"

	"github.com/bluebudgetz/gate/internal/util"
)

var (
	OfferedContentTypes = []string{
		"application/x-yaml",
		"application/yaml",
		"text/yaml",
		"application/json",
		"text/json",
		"application/xml",
		"text/xml",
		"text/html",
		"text/plain",
	}
)

func Render(w http.ResponseWriter, r *http.Request, v interface{}) {
	switch acceptedMimeType := util.NegotiateContentType(r.Header.Get("Accept"), OfferedContentTypes, ""); acceptedMimeType {
	case "application/x-yaml", "application/yaml", "text/yaml":
		w.Header().Set("Content-Type", acceptedMimeType)
		encoder := yaml.NewEncoder(w)
		defer encoder.Close()
		if err := encoder.Encode(v); err != nil {
			log.Warn().
				Str("mimeType", acceptedMimeType).
				Interface("data", v).
				Msg("Failed encoding data to response")
		}

	case "application/json", "text/json":
		w.Header().Set("Content-Type", acceptedMimeType)
		encoder := json.NewEncoder(w)
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(v); err != nil {
			log.Warn().
				Str("mimeType", acceptedMimeType).
				Interface("data", v).
				Msg("Failed encoding data to response")
		}

	case "application/xml", "text/xml":
		w.Header().Set("Content-Type", acceptedMimeType)
		encoder := xml.NewEncoder(w)
		if err := encoder.Encode(v); err != nil {
			log.Warn().
				Str("mimeType", acceptedMimeType).
				Interface("data", v).
				Msg("Failed encoding data to response")
		} else if err := encoder.Flush(); err != nil {
			log.Warn().
				Str("mimeType", acceptedMimeType).
				Interface("data", v).
				Msg("Failed encoding data to response")
		}

	case "text/html":
		w.WriteHeader(http.StatusOK) // send HTTP 200 since this is most probably a browser
		w.Header().Set("Content-Type", acceptedMimeType)
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		if _, err := w.Write([]byte("<!DOCTYPE html><html><body><pre>")); err != nil {
			log.Warn().
				Str("mimeType", acceptedMimeType).
				Interface("data", v).
				Msg("Failed encoding data to response")
		} else if err := encoder.Encode(v); err != nil {
			log.Warn().
				Str("mimeType", acceptedMimeType).
				Interface("data", v).
				Msg("Failed encoding data to response")
		} else if _, err := w.Write([]byte("</pre></body></html>")); err != nil {
			log.Warn().
				Str("mimeType", acceptedMimeType).
				Interface("data", v).
				Msg("Failed encoding data to response")
		}

	case "text/plain":
		w.Header().Set("Content-Type", acceptedMimeType)
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(v); err != nil {
			log.Warn().
				Str("mimeType", acceptedMimeType).
				Interface("data", v).
				Msg("Failed encoding data to response")
		}

	default:
		w.WriteHeader(http.StatusNotAcceptable)
	}
}
