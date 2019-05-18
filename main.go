package main

//go:generate go-bindata -o ./internal/assets/bindata.go -ignore ".DS_Store" -pkg assets deployments/...

import (
	. "github.com/bluebudgetz/common/pkg/logging"
	"github.com/bluebudgetz/gate/internal"
	"net/http"
	"os"
)

func main() {

	// Ensure that panics are logged through our logger
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				Log.WithError(err).Error("unhandled error (raised via panic)")
			} else {
				Log.WithField("panic", r).Errorf("unhandled panic")
			}
			os.Exit(2)
		}
	}()

	app, err := internal.New()
	if err != nil {
		panic(err)
	}

	err = app.Run()
	if err != nil && err != http.ErrServerClosed {
		Log.WithError(err).Error("server failed")
		os.Exit(1)
	}
}
