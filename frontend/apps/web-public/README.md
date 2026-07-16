# @radartcell/web-public

Публичный сайт RadarTcell — интерактивный радар технологий с каталогом.

## Стек
- Vite + React 18 + TypeScript (strict)
- TanStack Query 5 (кеш + автообновление)
- React Router 6 (code splitting per route)
- Tailwind CSS + общий пакет `@radartcell/ui`
- Общие типы API: `@radartcell/api`

## Разработка

```powershell
pnpm install
pnpm --filter @radartcell/web-public dev
```

Откройте http://localhost:5173. Прокси `/api → localhost:8080` уже настроен в `vite.config.ts`.

## Сборка

```powershell
pnpm --filter @radartcell/web-public build
```

Готовые артефакты — в `dist/`. Хостятся отдельно (nginx / CDN); API живёт на другом домене (см. `VITE_API_BASE`).

## Env

Смотрите `.env.example`. Скопируйте в `.env.local` и адаптируйте.
