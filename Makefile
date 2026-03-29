.PHONY: run build migrate seed bootstrap reset-db tidy test

run:
	go run ./cmd/api/main.go

build:
	go build -o bin/e-proc-api ./cmd/api/main.go

migrate:
	DB_MIGRATE=true go run ./cmd/api/main.go

seed:
	DB_SEED=true go run ./cmd/api/main.go

bootstrap:
	DB_MIGRATE=true DB_SEED=true go run ./cmd/api/main.go

reset-db:
	DB_RESET=true DB_MIGRATE=true DB_SEED=true go run ./cmd/api/main.go

tidy:
	go mod tidy

test:
	go test ./...
