# Миграции базы данных

Схема управляется через [golang-migrate](https://github.com/golang-migrate/migrate).

## Структура

```
migrations/
  000001_extensions.up.sql        # расширения + общая функция set_updated_at()
  000001_extensions.down.sql
  000002_create_trends.*.sql
  ...
  000013_soft_delete_catalog_entities.*.sql
  seeds/                          # демо-данные (НЕ управляются migrate)
    0001_demo_data.sql
    0002_more_technologies.sql
    0003_rich_demo_data.sql
  README.md
```

Каждая миграция — это пара файлов:

- `NNNNNN_name.up.sql` — применение изменения,
- `NNNNNN_name.down.sql` — откат изменения.

Версия (`NNNNNN`) монотонно возрастает и не переиспользуется. Новые
изменения добавляются как новая миграция, существующие файлы не редактируются.

## Соглашения

- Все DDL идемпотентны (`IF NOT EXISTS` / `IF EXISTS`), чтобы повторный
  прогон не падал.
- Поля `created_at` / `updated_at` есть у всех основных таблиц; `updated_at`
  обновляется триггером `set_updated_at()` (создаётся в `000001`).
- Удаление — мягкое (`deleted_at`), уникальность по `slug` действует только
  для «живых» строк (частичные уникальные индексы).
- Каждый `up` имеет симметричный `down`.

## Применение

Накатить все миграции:

```powershell
migrate -path migrations -database "postgres://radar_tcell:radar_tcell_password@localhost:15433/radar_tcell?sslmode=disable" up
```

Откатить последнюю:

```powershell
migrate -path migrations -database "$DATABASE_URL" down 1
```

## Сиды (демо-данные)

Сиды лежат в `migrations/seeds/` и применяются отдельно через `psql`
(golang-migrate их игнорирует). Применяются по возрастанию имени:

```powershell
psql "$DATABASE_URL" -f migrations/seeds/0001_demo_data.sql
psql "$DATABASE_URL" -f migrations/seeds/0002_more_technologies.sql
psql "$DATABASE_URL" -f migrations/seeds/0003_rich_demo_data.sql
```

Или одной командой через скрипт:

```powershell
./scripts/init-db.ps1 -ApplySeeds
```

## Создание новой миграции

```powershell
migrate create -ext sql -dir migrations -seq <короткое_имя>
```

После генерации заполните `*.up.sql` и обязательно — `*.down.sql`.
