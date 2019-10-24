package infra

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type statusCodeSavingResponseWriter struct {
	w          http.ResponseWriter
	statusCode int
}

func (s *statusCodeSavingResponseWriter) Header() http.Header {
	return s.w.Header()
}

func (s *statusCodeSavingResponseWriter) Write(bytes []byte) (int, error) {
	return s.w.Write(bytes)
}

func (s *statusCodeSavingResponseWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.w.WriteHeader(statusCode)
}

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
		// Use custom response writer that saves the status code (so we can log it)
		ww := &statusCodeSavingResponseWriter{statusCode: -1, w: w}

		// Invoke target handler, but time the request
		start := time.Now()
		recovered, stack := requestLoggerInvoke(next, r, ww)
		targetStopTime := time.Now()

		// TODO: support logging actual request/response body
		var event *zerolog.Event
		var message string
		if recovered != nil {
			event = log.Error()
			if err, ok := recovered.(error); ok {
				event = event.Err(err)
			} else {
				event = event.Interface(zerolog.ErrorFieldName, recovered)
			}
			event = event.Str(zerolog.ErrorStackFieldName, string(stack))
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
			Int("status", ww.statusCode).
			Msg(message)
	})
}
