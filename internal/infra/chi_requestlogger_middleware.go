package infra

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func RequestLogger(next http.Handler) http.Handler {
	zerologDictFromHeader := func(headers http.Header) *zerolog.Event {
		dict := zerolog.Dict()
		for name, values := range headers {
			dict.Strs(name, values)
		}
		return dict
	}

	requestLoggerInvoke := func(next http.Handler, r *http.Request, w http.ResponseWriter) (recovered interface{}, stack []byte) {
		defer func() {
			if recovered = recover(); recovered != nil {
				stack = debug.Stack()
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
		return nil, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the response to enable access to the final status code & response contents
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Tee the response to an additional buffer so we can look inside & log it
		buf := bytes.Buffer{}
		ww.Tee(&buf)

		// Invoke target handler, but time the request
		start := time.Now()
		recovered, stack := requestLoggerInvoke(next, r, ww)
		targetStopTime := time.Now()

		// Create the logging event; an Error event for panics, and Info event otherwise
		var event *zerolog.Event
		var message string
		if recovered != nil {
			if _, ok := recovered.(error); !ok {
				recovered = errors.Errorf("%s", recovered)
			}
			event = log.Error().
				Err(recovered.(error)).
				Str(zerolog.ErrorStackFieldName, fmt.Sprintf("%s\n%s", recovered.(error).Error(), string(stack))).
				Bytes("out", buf.Bytes())
			message = "HTTP request failed"
		} else {
			event = log.Info()
			message = "HTTP request completed"
		}

		// Log it!
		event.
			Str("rid", w.Header().Get(RequestIDHeaderName)).
			Str("remoteAddr", r.RemoteAddr).
			Str("proto", r.Proto).
			Str("method", r.Method).
			Str("uri", r.RequestURI).
			Str("host", r.Host).
			Dict("requestHeaders", zerologDictFromHeader(r.Header)).
			Dict("responseHeaders", zerologDictFromHeader(w.Header())).
			TimeDiff("elapsed", targetStopTime, start).
			Int("bytesWritten", ww.BytesWritten()).
			Int("status", ww.Status()).
			Msg(message)
	})
}
