# RadarTcell Frontend

Single-page frontend serves from `web/` and is mounted by the Go backend via
`withFrontend(...)` in `cmd/api/main.go`. Files in this folder:

- `index.html` — shell разметка приложения и сайдбар.
- `styles.css` — стили (тёмная тема, радар, карточки, модалка).
- `app.js` — SPA логика: роутинг по hash, радар, каталог, админ-CRUD.

## Что реализовано

- `#/explore` — интерактивный SVG-радар с секторами по трендам, кольцами по
  TRL и подсветкой по выбранному тренду через легенду.
- `#/catalog` — каталог карточек технологий с фильтрами (поиск, тренд, тег,
  ЦУР, организация, диапазон TRL).
- Карточка технологии открывается модальным окном из любого места:
  - hero-картинка из `image_url`,
  - стадия (Идея / Прототип / Продукт) и TRL,
  - подробное описание (`description_full`),
  - 4 кастомные метрики с прогресс-барами,
  - связи: теги, ЦУР, организации с логотипами.
- Списки сущностей: `#/trends`, `#/tags`, `#/sdgs`, `#/organizations`,
  `#/metrics`.
- Перекрёстные выборки: `#/trend/:slug`, `#/tag/:slug`, `#/sdg/:code`,
  `#/organization/:slug` — все используют те же красивые карточки.
- Админ-CRUD через `#/admin/...` (login через `/api/admin/login`).

## Демо-данные

Добавлен файл `migrations/seed_0003_rich_demo_data.sql` — 30 технологий по
6 трендам с описаниями на русском, картинками с Unsplash и логотипами
организаций через clearbit.

Применить:

```powershell
psql "postgres://radar_tcell:radar_tcell_password@localhost:15433/radar_tcell?sslmode=disable" `
  -f migrations/seed_0003_rich_demo_data.sql
```

Или через скрипт инициализации:

```powershell
.\scripts\init-db.ps1 -ApplySeeds -ForceSeed
```

## Запуск

1. Поднять Postgres и Redis (см. `deploy/docker-compose.yml`).
2. Применить миграции и сиды.
3. Запустить API:

```powershell
go run ./cmd/api
```

4. Открыть `http://localhost:8080/`.

## Заметки

- `API_BASE` берётся из `window.__API_BASE__` или из `window.location.origin`.
- Токен админки хранится в `localStorage` под ключом `rt_admin_token`.
- Картинки трендов и технологий приходят из БД (`image_url`). Если поле
  пустое, используется фолбэк.
