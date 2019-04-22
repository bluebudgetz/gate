FROM golang:1.12.1 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./main.go ./
COPY ./internal ./internal/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./gate ./main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/gate /
EXPOSE 3000
ENTRYPOINT ["/gate"]
