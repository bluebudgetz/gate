module github.com/bluebudgetz/gate

go 1.13

require (
	cloud.google.com/go/pubsub v1.0.1
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/cors v1.0.0
	github.com/go-chi/docgen v1.0.5
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/golang/snappy v0.0.1 // indirect
	github.com/golangly/errors v0.2.2
	github.com/golangly/log v0.2.0
	github.com/golangly/webutil v0.2.0
	github.com/jessevdk/go-flags v1.4.1-0.20181221193153-c0795c8afcf4
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/prometheus/client_golang v1.2.1
	github.com/tidwall/pretty v1.0.0 // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v1.0.0 // indirect
	go.mongodb.org/mongo-driver v1.1.2
	gopkg.in/go-playground/validator.v9 v9.30.0
)

replace github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
