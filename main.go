package main

//go:generate go-bindata -o ./internal/assets/bindata.go -ignore ".DS_Store" -pkg assets deployments/...
//go:generate go run cmd/gqlgen/main.go

import (
	. "github.com/bluebudgetz/common/pkg/logging"
	"github.com/bluebudgetz/gate/internal"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
)

func main() {

	// Ensure that panics are logged through our logger
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				Log.WithError(err).Error("unhandled panic")
			} else {
				Log.Errorf("unhandled panic: %+v", r)
			}
			os.Exit(1)
		}
	}()

	gate := internal.NewGate()
	err := gate.Run()
	if err != nil {
		panic(err)
	}
}
