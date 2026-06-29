# Fullstack Monorepo Template (Go API)

A production-ready monorepo template with a **Go** backend:

- **`apps/api`** — Go (chi + pgx + sqlc) REST API: auth, JWT, email confirmation, password reset
- **`apps/web`** — React 19 + Vite SSR + Tailwind 4 (landing page, Terms, Privacy)
- **`apps/mobile`** — Flutter (BLoC architecture, light/dark themes, auth-ready)
- **`deploy/`** — Nginx reverse proxy + VPS Docker Compose
- **`.github/workflows/`** — GitHub Actions: test the Go API, build & push images to GHCR, deploy via Coolify

> This is the Go counterpart of the original NestJS template. The HTTP contract
> (routes, request/response shapes, error envelope) is identical, so the web and
> mobile apps work against either backend unchanged.

## The API stack

| Concern | Library |
|---------|---------|
| HTTP router | [`go-chi/chi`](https://github.com/go-chi/chi) |
| DB driver / pool | [`jackc/pgx`](https://github.com/jackc/pgx) |
| Typed queries | [`sqlc`](https://sqlc.dev) (generated from SQL) |
| Migrations | [`golang-migrate`](https://github.com/golang-migrate/migrate) (embedded, auto-run on boot) |
| Auth | [`golang-jwt`](https://github.com/golang-jwt/jwt) + [`bcrypt`](https://pkg.go.dev/golang.org/x/crypto/bcrypt) |
| Validation | [`go-playground/validator`](https://github.com/go-playground/validator) |
| Config | [`caarlos0/env`](https://github.com/caarlos0/env) |
| Logging | `log/slog` (stdlib) |
| Security / CORS / rate limit | [`unrolled/secure`](https://github.com/unrolled/secure), [`go-chi/cors`](https://github.com/go-chi/cors), [`go-chi/httprate`](https://github.com/go-chi/httprate) |
| API docs | OpenAPI 3 + Swagger UI at `/docs` |

## Getting started

```bash
# 1. Clone or use this as a GitHub template
git clone https://github.com/senorMk/go-monorepo-template my-project
cd my-project

# 2. Run the init script — it replaces all placeholders with your project details
bash scripts/init.sh

# 3. Install web dependencies
npm install

# 4. Start Postgres
docker compose up -d

# 5. Configure the API
cp apps/api/.env.example apps/api/.env   # edit secrets

# 6. Start the API and web app
npm run api   # http://localhost:3000  (migrations run automatically on boot)
npm run web   # http://localhost:4000

# 7. Start the Flutter app
cd apps/mobile
flutter pub get
flutter run
```

Requirements: **Go 1.26+**, **Node 22+**, **Docker**, and (for mobile) **Flutter**.

## Placeholders

The init script replaces these tokens across all files (including the Go module path):

| Token | Example | Used in |
|-------|---------|---------|
| `APP_NAME` | `my-app` | package.json, Go module path, Docker image tags, Nginx routes |
| `APP_DISPLAY_NAME` | `My App` | OpenAPI title, web title, Flutter app title |
| `APP_SNAKE` | `my_app` | Flutter package name, Postgres user/db |
| `APP_DOMAIN` | `myapp.com` | Nginx, email links, CORS |
| `GITHUB_USERNAME` | `johndoe` | Go module path, GHCR image paths, CI/CD |

## Project structure

```
.
├── apps/
│   ├── api/                 # Go API
│   │   ├── cmd/api/         # main entrypoint
│   │   ├── internal/
│   │   │   ├── auth/        # signup, signin, refresh, email confirm, password reset
│   │   │   ├── profile/     # GET/PATCH /v1/profile/me
│   │   │   ├── config/      # env-based configuration
│   │   │   ├── db/          # pgx pool, migrations, sqlc queries
│   │   │   │   ├── migrations/  # SQL migrations (embedded)
│   │   │   │   ├── queries/     # SQL for sqlc
│   │   │   │   └── sqlc/        # generated query layer
│   │   │   ├── docs/        # OpenAPI spec + Swagger UI
│   │   │   ├── email/       # email service stub (swap in Resend/SES/Postmark)
│   │   │   ├── httpx/       # error envelope + request validation
│   │   │   ├── server/      # router + middleware
│   │   │   └── users/       # shared user response shape
│   │   ├── sqlc.yaml
│   │   ├── Makefile
│   │   └── Dockerfile
│   ├── web/                 # React + Vite SSR
│   └── mobile/              # Flutter
├── deploy/
│   ├── nginx/app.conf
│   └── docker-compose.vps.yml
├── .github/workflows/deploy.yml
├── docker-compose.yml           # Local Postgres only
├── docker-compose.coolify.yml   # Coolify managed stack
└── scripts/init.sh
```

## API endpoints (out of the box)

```
GET   /health
POST  /v1/auth/signup
POST  /v1/auth/signin
POST  /v1/auth/signout
POST  /v1/auth/refresh
POST  /v1/auth/forgot-password
POST  /v1/auth/reset-password
POST  /v1/auth/confirm-email
GET   /v1/profile/me
PATCH /v1/profile/me
```

Swagger UI is available at `http://localhost:3000/docs` in development.

See [`apps/api/README.md`](apps/api/README.md) for API-specific docs (adding
endpoints, regenerating sqlc, migrations).

## Deploying

1. Add secrets to your GitHub repo: `COOLIFY_WEBHOOK_URL`
2. Push to `main` — GitHub Actions tests the Go API, then builds and pushes Docker images to GHCR
3. Coolify picks up the webhook and redeploys

For VPS (without Coolify), use `deploy/docker-compose.vps.yml` directly.
