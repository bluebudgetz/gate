# gate ![Build master](https://github.com/bluebudgetz/gate/workflows/Build%20master/badge.svg) [![GoDoc](https://godoc.org/github.com/bluebudgetz/gate?status.svg)](http://godoc.org/github.com/bluebudgetz/gate) [![Report card](https://goreportcard.com/badge/github.com/bluebudgetz/gate)](https://goreportcard.com/report/github.com/bluebudgetz/gate) [![Sourcegraph](https://sourcegraph.com/github.com/bluebudgetz/gate/-/badge.svg)](https://sourcegraph.com/github.com/bluebudgetz/gate?badge) [![Coverage Status](https://coveralls.io/repos/github/bluebudgetz/gate/badge.svg?branch=master)](https://coveralls.io/github/bluebudgetz/gate?branch=master)

Bluebudgetz API gateway.

## Tasks

- [ ] Implement authentication (JWT)
- [ ] Increase test coverage
- [ ] Upgrade Golang
- [ ] Publish test coverage as comments in PRs

## Contributing

First, please read the [contribution guide](.github/CONTRIBUTING.md).

### Development environment

#### Linux

TBD.

#### MacOS

Do not use the `seabolt_static` tag (via `--tags seabolt_static`) as it seems not to work at the moment.

Instead, simply install `michael-simons/homebrew-seabolt/seabolt` brew package, which will make `Seabolt` propertly available on your MacOS workstation. Then simply use a standard `go build ..` command, like so:
 
```bash
# Install Seabolt:
$ brew install michael-simons/homebrew-seabolt/seabolt

# Verify it works:
$ BOLT_USER="" seabolt-cli run "UNWIND range(1, 23) AS n RETURN n"

# Build and run:
$ go build ...
```

## License

[GNUv3](./LICENSE)
