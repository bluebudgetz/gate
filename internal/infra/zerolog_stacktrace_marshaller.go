package infra

import (
	"bytes"
	"fmt"
	"github.com/kr/text"
)

type causer interface {
	Cause() error
}

func stackTraceMarshaller(err error) interface{} {
	buf := bytes.NewBuffer(make([]byte, 0, 100))
	for err != nil {
		buf.WriteString(fmt.Sprintf("%s", err))
		if stErr, ok := err.(stackTracer); ok {
			buf.WriteString("\n")
			st := stErr.StackTrace()
			for _, frame := range st {
				buf.WriteString(text.Indent(fmt.Sprintf("%+v\n", frame), "  "))
			}
		}
		if causer, ok := err.(causer); ok {
			err = causer.Cause()
		} else {
			err = nil
		}
	}
	return buf.String()
}
