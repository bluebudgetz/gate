FROM ubuntu:18.04 AS builder
RUN apt-get update && \
    apt-get install -y curl pkg-config && \
    cd /tmp && curl -O https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf /tmp/go1.14.2.linux-amd64.tar.gz && rm -fv /tmp/go1.14.2.linux-amd64.tar.gz && \
    curl -sSL -O https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
    dpkg -i seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
    rm -fv seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
    rm -rf /var/lib/apt/lists/*
ENV GOOS="linux"
ENV GOARCH="amd64"
ENV PATH=/usr/local/go/bin:/go/bin:$PATH
WORKDIR /go/src/
COPY go.mod go.sum ./
RUN go mod download
COPY ./cmd ./cmd/
COPY ./internal ./internal/
RUN go build -o ../bin/gate ./cmd/main.go

FROM ubuntu:18.04
RUN apt-get update && apt-get install -y curl && \
    curl -sSL -O https://github.com/neo4j-drivers/seabolt/releases/download/v1.7.4/seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
    dpkg -i seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
    rm -fv seabolt-1.7.4-Linux-ubuntu-18.04.deb && \
    rm -rf /var/lib/apt/lists/*
COPY --from=builder /go/bin/gate /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/gate"]
