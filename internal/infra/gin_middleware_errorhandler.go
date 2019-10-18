package infra

import (
	"github.com/bluebudgetz/gate/internal/infra/render"
	"github.com/gin-gonic/gin"
)

const ErrorTypeHTTP gin.ErrorType = 1 << 61

func ginHandlerForLastErrorOfType(defaultStatusCode int, errorType gin.ErrorType, staticMessage string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if !c.Writer.Written() {
			e := c.Errors.ByType(errorType).Last()
			if e != nil {
				statusCode := defaultStatusCode
				if e.IsType(ErrorTypeHTTP) {
					statusCode = int(e.Type &^ ErrorTypeHTTP)
				}
				if staticMessage != "" {
					render.Render(c, statusCode, []string{staticMessage})
				} else {
					render.Render(c, statusCode, e.Error())
				}
			}
		}
	}
}

func ginHandlerForErrorsOfType(statusCode int, errorType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if !c.Writer.Written() {
			ginErrors := c.Errors.ByType(errorType)
			if len(ginErrors) > 0 {
				render.Render(c, statusCode, ginErrors.Errors())
			}
		}
	}
}
