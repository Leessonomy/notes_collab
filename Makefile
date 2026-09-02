.PHONY: dev dev-web dev-api build-web build-api test-web test-api down migrate-up migrate-down migrate-status migrate-create

-include .env
export

DATABASE_URL ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable

dev:
	docker compose up --build

dev-web:
	cd apps/web && npm start

dev-api:
	cd apps/api && go run ./cmd/server

build-web:
	cd apps/web && npm run build

build-api:
	cd apps/api && go build -o bin/server ./cmd/server

test-web:
	cd apps/web && npm test

test-api:
	cd apps/api && go test ./...

down:
	docker compose down

migrate-up:
	cd apps/api && go tool goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	cd apps/api && go tool goose -dir migrations postgres "$(DATABASE_URL)" down

migrate-status:
	cd apps/api && go tool goose -dir migrations postgres "$(DATABASE_URL)" status

migrate-create:
	cd apps/api && go tool goose -dir migrations create $(name) sql
