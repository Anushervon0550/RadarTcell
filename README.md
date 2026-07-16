# RadarTcell

Радар технологий: интерактивная визуализация трендов, ЦУР, организаций и метрик TRL. Продакшен-ориентированный монорепозиторий с чётким разделением бэка, публичного сайта и админ-панели.

## Структура

```
radartcell/
├── backend/                 Go API (клин-арх, chi, pgx, zap, Redis)
│   ├── cmd/
│   ├── internal/
│   ├── migrations/
│   ├── docs/                OpenAPI 3
│   ├── scripts/
│   └── Dockerfile
├── frontend/                pnpm workspace (React + Vite + TS)
│   ├── apps/
│   │   ├── web-public/      Публичный сайт    (порт 5173)
│   │   └── web-admin/       Админ-панель      (порт 5174)
│   └── packages/
│       ├── api/             Общие TS-типы + ApiClient
│       └── ui/              UI-кит (Tailwind preset + компоненты)
├── deploy/
│   └── nginx/               Prod nginx-конфиги для web-public / web-admin
├── .github/workflows/       CI: backend (Go) + frontend (Vite)
├── docker-compose.yml       Локальный/dev prod-стек
├── .env.example
└── README.md
```

Каждый слой можно разрабатывать и деплоить независимо: у бэка свой `go.mod`, у фронта — свой `pnpm-workspace`. В корне только оркестрация (Docker Compose, CI, env).

## Стек

**Backend** — Go 1.24, chi, pgx v5, zap, go-redis, JWT (bcrypt), golang-migrate, swaggo/openapi.

**Frontend** — React 18 + TypeScript (strict, `noUncheckedIndexedAccess`), Vite 6, Tailwind CSS 3, TanStack Query 5, React Router 6. Админка также: React Hook Form + Zod, Zustand (persisted auth), универсальный CRUD-фреймворк.

## Локальный запуск

```powershell
# 1. Секреты
copy .env.example .env

# 2. Инфраструктура (Postgres + Redis)
docker compose up -d postgres redis

# 3. Миграции + сиды (bash-эквивалент: см. backend/README.md)
cd backend
./scripts/init-db.ps1 -ApplySeeds
cd ..

# 4. API
cd backend
go run ./cmd
# → http://localhost:8080
cd ..

# 5. Frontend (в отдельном терминале)
cd frontend
pnpm install
pnpm dev:public   # http://localhost:5173
pnpm dev:admin    # http://localhost:5174
```

## Прод-стек в Docker

```powershell
docker compose up -d --build
```

Поднимается пять сервисов:

| Сервис        | Порт (хост) | Роль                                                  |
| ------------- | ----------- | ----------------------------------------------------- |
| `postgres`    | 15433       | БД                                                    |
| `redis`       | 6379        | Кеш                                                   |
| `app`         | 8080        | Go API                                                |
| `web-public`  | 5173        | nginx + сборка публичного сайта                       |
| `web-admin`   | 5174        | nginx + сборка админ-панели (X-Robots-Tag: noindex)   |

## Продакшен-развёртывание (краткое)

1. Отдельные поддомены:
   - `radartcell.example.com` → контейнер `web-public`
   - `admin.radartcell.example.com` → контейнер `web-admin` (за IP-allowlist или VPN)
   - `api.radartcell.example.com` → Go API
2. Публичный TLS через nginx-proxy / Traefik / cloud LB. После включения TLS
   выставите `SECURITY_HSTS=true`.
3. Строгий CORS: только те origin-ы, которые действительно нужны.
4. Секреты (`JWT_SECRET`, `ADMIN_PASSWORD`, `POSTGRES_PASSWORD`, `REDIS_PASSWORD`)
   лежат в Vault / SSM / Kubernetes Secrets — не в `.env` на сервере.
5. Логи собираются zap в JSON; настройте sink (Loki / ELK / CloudWatch).

## Что сделано по-мидловски

- **Backend отдельно от frontend** — свой `go.mod`, свой `pnpm-workspace`, свои Dockerfile'ы, свой CI-джоб.
- **Публичный сайт и админка — разные приложения**, разные бандлы, разные CSP,
  разные release-циклы, admin помечен `X-Robots-Tag: noindex` и легко закрывается
  IP-allowlist на nginx.
- Универсальный `CrudResource<T>` — новая сущность в админке добавляется одним
  файлом на ~40 строк.
- Zod-схемы валидации, `noUncheckedIndexedAccess` в TS.
- Code splitting per route (React.lazy), manual chunks в Rollup.
- Иммутабельное кеширование ассетов + `no-cache` для `index.html` в nginx.
- Security headers (CSP, HSTS, X-Frame-Options и т.д.) — и в Go, и в nginx.
- `SERVE_FRONTEND=false` по умолчанию: Go отдаёт только API, статика — отдельно.

## Roadmap (следующие итерации)

1. Redis-backed rate limiting (сейчас in-memory — не масштабируется).
2. Prometheus `/metrics` + OpenTelemetry traces.
3. Refresh tokens + переход админки на HttpOnly-cookie.
4. Audit log для админ-действий в БД.
5. Fine-grained RBAC (сейчас все админы равны).
6. i18n на фронте через `react-i18next` (бэк уже поддерживает locale).
