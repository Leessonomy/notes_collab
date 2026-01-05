.PHONY: dev dev-web dev-api build-web build-api test-web test-api down

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
