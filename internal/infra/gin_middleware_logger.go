package infra

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

func ginLogger(disableRequestLog bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Invoke target handler, but time the request
		start := time.Now()
		c.Next()
		targetStopTime := time.Now()

		requestHeaders := zerolog.Dict()
		for name, values := range c.Request.Header {
			requestHeaders.Strs(name, values)
		}

		responseHeaders := zerolog.Dict()
		for name, values := range c.Writer.Header() {
			responseHeaders.Strs(name, values)
		}

		var event *zerolog.Event
		var message string

		if len(c.Errors) > 0 {
			var errors []error
			for _, e := range c.Errors {
				errors = append(errors, e)
			}
			event = log.Error().Stack().Errs("errors", errors)
			message = "HTTP request failed"
		} else if !disableRequestLog {
			event = log.Info()
			message = "HTTP request completed"
		} else {
			return // No error, and request log is disabled
		}

		// TODO: support logging actual request/response body

		event.
			Str("id", c.GetString("X-Request-ID")).
			Str("remoteAddr", c.ClientIP()).
			Str("proto", c.Request.Proto).
			Str("method", c.Request.Method).
			Str("uri", c.Request.RequestURI).
			Str("host", c.Request.Host).
			Dict("requestHeaders", requestHeaders).
			Dict("responseHeaders", responseHeaders).
			TimeDiff("elapsed", targetStopTime, start).
			Int("status", c.Writer.Status()).
			Int("bytesWritten", c.Writer.Size()).
			Msg(message)
	}
}
