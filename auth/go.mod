module auth

go 1.23.1

require github.com/google/uuid v1.6.0

require (
	github.com/gorilla/mux v1.8.1
	golang.org/x/crypto v0.27.0
)

require (
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lib/pq v1.10.9
)

require (
	github.com/BurntSushi/toml v1.4.0
	go.uber.org/zap v1.27.0
)

require go.uber.org/multierr v1.10.0 // indirect
