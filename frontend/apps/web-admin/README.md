# @radartcell/web-admin

Админ-панель RadarTcell — отдельное SPA-приложение, задеплоенное на своём поддомене
(`admin.radartcell.example.com`).

## Стек
- Vite + React 18 + TypeScript (strict)
- TanStack Query 5 (кеш + оптимистичные апдейты)
- React Hook Form + Zod (типизированные формы с валидацией)
- Zustand (хранилище JWT-токена, persisted)
- Универсальный CRUD-фреймворк на `CrudResource<T>`

## Разработка

```powershell
pnpm install
pnpm --filter @radartcell/web-admin dev
```

http://localhost:5174, прокси на API уже настроен.

## Роутинг

- `/login` — вход
- `/` — дашборд
- `/technologies` `/trends` `/tags` `/sdgs` `/organizations` `/metrics` — CRUD-разделы
- `/users` — управление админ-аккаунтами

## Как добавить новый CRUD

1. Создайте файл в `src/crud/resources/<name>.ts`, экспортируйте `CrudResource<T>` с полями (`fields`), колонками таблицы (`columns`), и Zod-схемой.
2. Реэкспортируйте в `resources/index.ts`.
3. Добавьте маршрут в `App.tsx`.

Ничего больше — форма, валидация, таблица и мутации уже готовы.
