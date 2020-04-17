FROM golang:1.13.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd/
COPY ./internal ./internal/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build --tags seabolt_static -o ./bin/gate ./cmd/main.go

FROM scratch
WORKDIR /tmp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/bin/gate /
EXPOSE 3000
WORKDIR /
ENTRYPOINT ["/gate"]
