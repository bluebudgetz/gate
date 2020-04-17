FROM golang:1.13.3 AS builder
RUN apt-get update && apt-get install -y libssl-dev && \
    wget https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
    dpkg -i seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
    rm seabolt-1.7.4-Linux-ubuntu-18.04.deb
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd/
COPY ./internal ./internal/
RUN GOOS=linux GOARCH=amd64 go build --tags seabolt_static -o ./bin/gate ./cmd/main.go

FROM ubuntu
RUN apt-get update && apt-get install -y wget libssl-dev && \
    wget https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
    dpkg -i seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
    rm seabolt-1.7.4-Linux-ubuntu-18.04.deb
WORKDIR /tmp
COPY --from=builder /app/bin/gate /
EXPOSE 3000
WORKDIR /
ENTRYPOINT ["/gate"]
