# NotesCollab

A monorepo containing the Angular frontend (`apps/web`) and the Go backend (`apps/api`).

## Structure

```
apps/
  web/   Angular app (UI), see apps/web/README.md
  api/   Go backend (stdlib net/http), serves /api/*
```

## Development

Run everything together with Docker Compose:

```bash
docker compose up --build
```

- Frontend: http://localhost:4200
- Backend: http://localhost:8080 (proxied under `/api` by the frontend)

Or run each service locally (requires Node 22+ and Go 1.26+):

```bash
make dev-web   # cd apps/web && npm start
make dev-api   # cd apps/api && go run ./cmd/server
```

## Build & test

```bash
make build-web
make build-api
make test-web
make test-api
```

See `apps/web/README.md` for Angular-specific commands (scaffolding, e2e, etc).
