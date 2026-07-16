# RadarTcell — Backend

Go 1.24 HTTP API. Клин-архитектура (`ports → service → repository`), JWT-auth,
Postgres (pgx), Redis-кеш, chi router, structured logging (zap).

## Layout

```
backend/
├── cmd/                  Entry point (main.go, docs.go, performance tests)
├── internal/
│   ├── app/              Composition root: DI + router build
│   ├── domain/           Business entities, params, errors
│   ├── httpapi/          HTTP handlers, middlewares, router
│   ├── service/          Application services
│   ├── repository/       Data access (Postgres)
│   ├── cache/            Redis cache adapter
│   ├── logging/          zap wrapper
│   └── ports/            Hexagonal interfaces
├── migrations/           golang-migrate SQL migrations + seeds/
├── docs/                 OpenAPI 3 spec + swaggo output
├── scripts/              PowerShell helpers (init-db, smoke, preflight)
├── go.mod / go.sum
└── Dockerfile            Multi-stage build (non-root)
```

## Локальный запуск

```powershell
# Инфра (Postgres + Redis) поднимается из корневого docker-compose:
docker compose -f ../docker-compose.yml up -d postgres redis

# Миграции + сиды:
./scripts/init-db.ps1 -ApplySeeds

# Копируем .env (из корня репо) и запускаем:
copy ../.env.example ../.env
go run ./cmd
# API: http://localhost:8080
```

## Тесты

```powershell
go test ./...
go vet ./...
```

Perf-тесты требуют работающих Postgres + Redis:

```powershell
$env:DATABASE_URL="postgres://radar_tcell:radar_tcell_password@localhost:15433/radar_tcell?sslmode=disable"
$env:REDIS_ADDR="localhost:6379"
go test -v ./cmd -run TestGetTechnologiesLatency
```

## Сборка Docker

Из корня репо:

```powershell
docker compose build app
```

## Endpoints (кратко)

- `GET  /healthz` `/readyz`
- `GET  /swagger/*` `/openapi.yaml` (когда `SWAGGER_ENABLED=true`)
- `GET  /api/{trends,tags,sdgs,organizations,metrics,technologies,home}`
- `GET  /api/technologies/{slug}`
- `POST /api/admin/login` `GET /api/admin/me`
- `CRUD /api/admin/{trends,tags,sdgs,organizations,metrics,technologies,users}`

Полная спецификация — `docs/openapi.yaml`.
