generate:
	go generate ./cmd/... ./internal/...

build: generate
	go build -o bin/gate ./cmd/main.go
