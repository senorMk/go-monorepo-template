# APP_DISPLAY_NAME API

Go REST API for the APP_DISPLAY_NAME monorepo.

Stack: **chi** (router) · **pgx** + **sqlc** (Postgres) · **golang-migrate**
(migrations) · **golang-jwt** + **bcrypt** (auth) · **go-playground/validator**
(validation) · **caarlos0/env** (config) · **log/slog** (logging).

## Run locally

```bash
cp .env.example .env          # set DATABASE_URL + JWT secrets
go run ./cmd/api              # http://localhost:3000
```

Migrations in `internal/db/migrations` run automatically on boot, so the API is
ready as soon as it starts. Swagger UI is served at `/docs` (development only).

Common tasks are in the `Makefile`:

```bash
make dev          # run the API
make build        # compile to bin/api
make test         # go test ./...
make vet          # go vet ./...
make sqlc         # regenerate the query layer
make migrate-up   # apply migrations with the CLI (DB_URL=...)
```

## Layout

```
cmd/api/                 main(): config → migrate → connect → serve
internal/
  config/                env-based configuration (caarlos0/env)
  db/
    db.go                pgx pool + embedded migration runner
    migrations/          *.up.sql / *.down.sql (embedded via go:embed)
    queries/             *.sql consumed by sqlc
    sqlc/                generated, type-safe query layer (do not edit)
  auth/                  JWT, opaque tokens, signup/signin/refresh/reset/confirm
  profile/               /v1/profile/me (GET, PATCH)
  email/                 stub email service (swap for a real provider)
  httpx/                 JSON error envelope + request decode/validate
  users/                 shared user response shape
  docs/                  embedded OpenAPI spec + Swagger UI
  server/                router, middleware stack, /health
```

## Configuration

| Variable | Default | Notes |
|----------|---------|-------|
| `APP_ENV` | `development` | `production` disables `/docs` and switches to JSON logs |
| `PORT` | `3000` | |
| `DATABASE_URL` | — (required) | `postgresql://user:pass@host:5432/db` |
| `JWT_SECRET` | — (required) | HS256 signing secret for access tokens |
| `JWT_REFRESH_SECRET` | — | reserved for future use |
| `JWT_ACCESS_EXPIRES_IN` | `15m` | Go duration |
| `JWT_REFRESH_EXPIRES_IN` | `720h` | Go duration (30 days) |
| `CORS_ORIGINS` | `http://localhost:4000` | comma-separated |
| `EMAIL_FROM` | `info@APP_DOMAIN` | |
| `APP_SCHEME` | `APP_NAME` | deep-link scheme for confirm/reset emails |

## Error envelope

Every error response has the same shape:

```json
{ "statusCode": 401, "error": "UNAUTHORIZED", "message": "Invalid credentials" }
```

`error` codes: `BAD_REQUEST`, `UNAUTHORIZED`, `FORBIDDEN`, `NOT_FOUND`,
`CONFLICT`, `UNPROCESSABLE`, `INTERNAL_ERROR`.

## Adding an endpoint

1. **Schema** — add a migration:
   `make migrate-create name=add_widgets` then edit the generated `*.up.sql` /
   `*.down.sql` in `internal/db/migrations`.
2. **Queries** — write SQL in `internal/db/queries/*.sql` with sqlc annotations
   (`-- name: ListWidgets :many`), then run `make sqlc`.
3. **Feature package** — add `internal/widgets/` with a `Service` (business logic
   over `*sqlc.Queries`) and a `Handler` exposing `RegisterRoutes(r chi.Router,
   authMW)`.
4. **Wire it** — construct the handler in `internal/server/server.go` and call
   `RegisterRoutes`. Protect routes by passing `authMW` (see `profile`).
5. **Document it** — add the path/schema to `internal/docs/openapi.yaml`.

## Notes

- Opaque tokens (refresh / reset / confirm) are stored as SHA-256 hashes; a
  database leak does not expose usable tokens.
- Passwords are hashed with bcrypt.
- Refresh tokens rotate on every `/v1/auth/refresh` (the old one is revoked).
- The `email` package only logs links — wire a real provider (Resend, SES,
  Postmark) before going live.
