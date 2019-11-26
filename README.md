# gate [![GoDoc](https://godoc.org/github.com/bluebudgetz/gate?status.svg)](http://godoc.org/github.com/bluebudgetz/gate) [![Report card](https://goreportcard.com/badge/github.com/bluebudgetz/gate)](https://goreportcard.com/report/github.com/bluebudgetz/gate) [![Sourcegraph](https://sourcegraph.com/github.com/bluebudgetz/gate/-/badge.svg)](https://sourcegraph.com/github.com/bluebudgetz/gate?badge)

Bluebudgetz API gateway.

### Tasks

- [ ] Implement authentication (JWT)
- [ ] Integrate [Cloud Code](https://cloud.google.com/code/docs/intellij/quickstart-IDEA)
- [ ] Gate should start even if MongoDB is down (just return 500 for requests)
- [ ] Thread ID implementation
    - Consider https://github.com/huandu/go-tls
    - Consider https://github.com/modern-go/gls
