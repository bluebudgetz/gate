package infra

import (
	"bytes"
	"fmt"
	"runtime/debug"

	"github.com/kr/text"
	"github.com/pkg/errors"
)

type causer interface {
	Cause() error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func stackTraceMarshaller(err error) interface{} {
	buf := bytes.NewBuffer(make([]byte, 0, 100))
	firstError := true
	for err != nil {
		buf.WriteString(fmt.Sprintf("%s", err))
		if stErr, ok := err.(stackTracer); ok {
			buf.WriteString("\n")
			st := stErr.StackTrace()
			for _, frame := range st {
				buf.WriteString(text.Indent(fmt.Sprintf("%+v\n", frame), "  "))
			}
		} else if firstError {
			buf.WriteString("\n")
			buf.WriteString(string(debug.Stack()))
		}
		firstError = false

		if causer, ok := err.(causer); ok {
			err = causer.Cause()
		} else {
			err = nil
		}
	}
	return buf.String()
}
