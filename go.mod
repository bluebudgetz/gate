module github.com/bluebudgetz/gate

require (
	github.com/99designs/gqlgen v0.8.3
	github.com/bluebudgetz/common v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang-migrate/migrate/v4 v4.3.0
	github.com/kr/pty v1.1.3 // indirect
	github.com/pkg/errors v0.8.1
	github.com/vektah/gqlparser v1.1.2
)

replace github.com/bluebudgetz/common => ../common
