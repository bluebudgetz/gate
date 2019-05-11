module github.com/bluebudgetz/gate

require (
	github.com/99designs/gqlgen v0.8.3
	github.com/bluebudgetz/common v0.0.1
	github.com/davecgh/go-spew v1.1.1
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/golang-migrate/migrate/v4 v4.3.0
	github.com/jackc/pgx v3.2.0+incompatible
	github.com/kr/pty v1.1.3 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rs/cors v1.6.0
	github.com/vektah/gqlparser v1.1.2
)

//replace github.com/bluebudgetz/common => ../common
