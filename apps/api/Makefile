# APP_DISPLAY_NAME API — common developer tasks.
# Tools used by some targets (install once):
#   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
#   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

DB_URL ?= $(DATABASE_URL)

.PHONY: dev run build test vet tidy sqlc migrate-up migrate-down migrate-create

## dev / run: start the API (auto-runs migrations on boot)
dev run:
	go run ./cmd/api

## build: compile the API binary to bin/api
build:
	go build -ldflags="-s -w" -o bin/api ./cmd/api

## test: run the test suite
test:
	go test ./...

## vet: run go vet
vet:
	go vet ./...

## tidy: tidy go.mod / go.sum
tidy:
	go mod tidy

## sqlc: regenerate the typed query layer from internal/db/queries
sqlc:
	sqlc generate

## migrate-up: apply all migrations (the API also auto-migrates on boot)
migrate-up:
	migrate -path internal/db/migrations -database "$(DB_URL)" up

## migrate-down: roll back the last migration
migrate-down:
	migrate -path internal/db/migrations -database "$(DB_URL)" down 1

## migrate-create name=add_something: scaffold a new migration pair
migrate-create:
	migrate create -ext sql -dir internal/db/migrations -seq $(name)
