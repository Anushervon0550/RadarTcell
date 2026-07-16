# RadarTcell — Frontend

Монорепо на **pnpm workspaces** с двумя SPA-приложениями и общими пакетами.

## Layout

```
frontend/
├── apps/
│   ├── web-public/       Публичный сайт (React + Vite)  — порт 5173
│   └── web-admin/        Админ-панель  (React + Vite)   — порт 5174
├── packages/
│   ├── api/              Общие TS-типы + типизированный ApiClient
│   └── ui/               UI-кит (Tailwind preset + Button/Card/Input/Modal/...)
├── package.json          Workspace root
├── pnpm-workspace.yaml
└── .npmrc
```

## Требования

- Node.js ≥ 20.11
- pnpm ≥ 9 (`corepack enable`)

## Установка

```powershell
pnpm install
```

## Разработка

```powershell
pnpm dev:public   # http://localhost:5173  → прокси /api → localhost:8080
pnpm dev:admin    # http://localhost:5174
```

## Прод-сборка

```powershell
pnpm build           # оба приложения параллельно
pnpm build:public    # только публичный сайт
pnpm build:admin     # только админка
```

Артефакты — в `apps/*/dist`. Хостятся отдельно (nginx или CDN), API живёт на другом
поддомене (см. `VITE_API_BASE`).

## Типизация и линтинг

```powershell
pnpm typecheck   # tsc --noEmit во всех пакетах
pnpm lint        # eslint во всех пакетах
```

## Как добавить новую CRUD-сущность в админке

1. `apps/web-admin/src/crud/resources/<name>.ts` — экспортировать
   `CrudResource<T>` (колонки, поля, Zod-схема, пути API).
2. Реэкспорт в `resources/index.ts`.
3. Роут в `apps/web-admin/src/App.tsx`:
   ```tsx
   <Route path="/<name>" element={<ResourceListPage resource={myResource} />} />
   ```

Форма, валидация, таблица, мутации, оптимистичный refetch — готовы.
