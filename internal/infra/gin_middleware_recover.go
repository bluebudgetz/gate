package infra

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func ginRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				if panicError, ok := r.(error); ok {
					// TODO: stacktrace is incorrect
					if _, ok := panicError.(stackTracer); ok {
						c.Error(gin.Error{Err: panicError, Type: gin.ErrorTypePrivate})
					} else {
						c.Error(gin.Error{Err: errors.WithStack(panicError), Type: gin.ErrorTypePrivate})
					}
				} else {
					// TODO: test this provide good stacktrace
					c.Error(gin.Error{Err: errors.WithStack(fmt.Errorf("%v", r)), Type: gin.ErrorTypePrivate})
				}
			}
		}()
		c.Next()
	}
}
