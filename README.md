# RadarTcell

Backend API for RadarTcell (Go + PostgreSQL + Redis + optional MinIO).

## Quick start (Docker)

```powershell
Copy-Item .env.local.example .env
docker compose -f deploy/docker-compose.yml --env-file .env up --build
Invoke-RestMethod "http://localhost:8080/healthz"
Invoke-RestMethod "http://localhost:8080/readyz"
powershell -ExecutionPolicy Bypass -File .\scripts\smoke.ps1 -BaseUrl "http://localhost:8080" -AdminUser "admin" -AdminPassword "admin123" -TrustedOrigin "http://localhost:3000"
```

`deploy/.env.dev-example` is for local reference only.
Do not use `deploy/.env` or any `sslmode=disable` settings in production.

For Docker local run, compose uses `DATABASE_URL_DOCKER` (or its internal default with host `postgres`).
Regular `DATABASE_URL` in `.env` can stay host-oriented for `go run`.

API endpoints:
- `http://localhost:8080/healthz`
- `http://localhost:8080/readyz`

API versioning:
- Current routes are available under both `/api/...` and `/api/v1/...`.
- Admin routes are available under both `/api/admin/...` and `/api/v1/admin/...`.

Dynamic metrics:
- `GET /api/technologies` and `GET /api/technologies/{slug}` include `custom_metrics` (dynamic metric values).
- Legacy fields `custom_metric_1..4` are kept for backward compatibility.

Admin soft-delete notes:
- `DELETE /api/admin/technologies/{slug}` performs soft delete (sets `deleted_at`).
- `PUT /api/admin/technologies/{slug}/restore` restores a soft-deleted technology.
- `GET /api/admin/technologies?include_deleted=true` includes deleted technologies in admin list.
- Admin deletes for trends/tags/sdgs/organizations/metrics use soft delete (`deleted_at`) and are hidden from public catalog endpoints.

## Local run (without Docker)

1. Copy local env template and fill secrets:

```powershell
Copy-Item .env.local.example .env
```

2. Run API:

```powershell
go run ./cmd/api
```

> For local Postgres from `deploy/docker-compose.yml`, `DATABASE_URL` must use `sslmode=disable`.

## Production env template

- Use `.env.prod.example` as a base for production values.
- Keep production secrets out of git (for example, in `.env.prod` on server/CI secret store).
- For production DB, use `sslmode=require` (or stricter).

## Required environment variables

Core:
- `DATABASE_URL`
- `APP_PORT` (default `8080`)
- `ENV` (default `production`)

Auth:
- `ADMIN_USER`
- `ADMIN_PASSWORD`
- `ADMIN_AUTH_MODE` (optional: `db_then_env` default, `db_only`, `env_only`)
- `JWT_SECRET`
- `JWT_TTL_HOURS` (default `8`)

Admin auth source:
- `db_then_env` (default): first `admin_users`, fallback to `ADMIN_USER/ADMIN_PASSWORD`.
- `db_only`: only `admin_users`.
- `env_only`: only `ADMIN_USER/ADMIN_PASSWORD`.

Database pool:
- `DATABASE_MAX_CONNS` (default `20`)

Redis/cache:
- `REDIS_ADDR`
- `REDIS_PASSWORD`
- `REDIS_DB`
- `CATALOG_CACHE_TTL_SECONDS`
- `TECHNOLOGY_CACHE_TTL_SECONDS`

MinIO/S3 (optional):
- `MINIO_ENDPOINT`
- `MINIO_ACCESS_KEY`
- `MINIO_SECRET_KEY`
- `MINIO_BUCKET`
- `MINIO_PUBLIC_URL`
- `MINIO_USE_SSL`
- `MINIO_PUBLIC_READ` (default `false`)

## Security notes

- Rate limiting uses `RemoteAddr` as key source and does not trust spoofable forwarded headers.
- CSRF middleware blocks state-changing requests when both `Origin` and `Referer` are missing.
- Slugs for admin CRUD are validated by regex: `^[a-z0-9]+(?:-[a-z0-9]+)*$`.
- New MinIO buckets are private by default. Public read policy is applied only if `MINIO_PUBLIC_READ=true`.

## CI

GitHub Actions workflow includes:
- `go vet ./...`
- `go test ./...`
- `golangci-lint`

