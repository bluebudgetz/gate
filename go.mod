module github.com/bluebudgetz/gate

go 1.13

require (
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/cors v1.0.0
	github.com/go-chi/docgen v1.0.5
	github.com/golang/mock v1.3.1 // indirect
	github.com/golangly/errors v0.2.2
	github.com/golangly/log v0.2.0
	github.com/golangly/webutil v0.3.0
	github.com/google/wire v0.4.0
	github.com/jessevdk/go-flags v1.4.1-0.20181221193153-c0795c8afcf4
	github.com/neo4j-drivers/gobolt v1.7.4 // indirect
	github.com/neo4j/neo4j-go-driver v1.7.4
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/prometheus/client_golang v1.2.1
	gopkg.in/go-playground/validator.v9 v9.30.0
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
