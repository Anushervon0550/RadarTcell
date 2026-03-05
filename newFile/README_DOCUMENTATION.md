# 📖 ДОКУМЕНТАЦИЯ ПРОЕКТА RadarTcell

## ✅ БЫСТРЫЙ СТАРТ (5 минут)

```powershell
cd C:\PROJECT\RadarTcell\deploy
docker compose up -d

cd C:\PROJECT\RadarTcell
migrate -path "./migrations" -database $env:DATABASE_URL up

go run .\cmd\api\main.go
```

Swagger: `http://localhost:8080/swagger/`

---

## 🧭 КАРТА ДОКУМЕНТАЦИИ

- `ПОЛНОЕ_ОБЪЯСНЕНИЕ_ПРОЕКТА.md` — основной понятный обзор проекта.
- `ШПАРГАЛКА.md` — команды, эндпоинты, быстрые подсказки.
- `АРХИТЕКТУРА_ВИЗУАЛЬНО.md` — схемы и потоки.
- `ПРАКТИЧЕСКОЕ_РУКОВОДСТВО.md` — как добавлять фичи.
- `FAQ.md` — частые вопросы.

---

## 🎯 ЧТО ЭТО

RadarTcell — REST API для каталога технологий (тренды, SDG, теги, организации, метрики). Сервер:
- считает координаты радара,
- делает фильтрацию и сортировку,
- отдает готовые данные фронтенду.

Важно:
- Админские эндпоинты: `/api/admin/*`.
- Локализация: `?locale=ru`.
- Системные: `/healthz`, `/readyz`, `/metrics`, `/openapi.yaml`, `/swagger/`.

---

## 🔑 ПОЛЕЗНЫЕ ССЫЛКИ

- OpenAPI: `http://localhost:8080/openapi.yaml`
- Metrics: `http://localhost:8080/metrics`
- Grafana: `http://localhost:3000`

---

## ✅ ПОСЛЕДНИЕ ОБНОВЛЕНИЯ

Добавлено:
- Admin i18n API для переводов трендов/технологий/метрик.
- Таблицы локализации (trend_i18n, technology_i18n, metric_definition_i18n).
- Мониторинг: /metrics, Prometheus + Grafana.
- Структурные логи (zap).
- Тесты для координат/фильтрации/сортировки.

Изменено:
- Публичные списки и технологии поддерживают `?locale=ru`.
- OpenAPI обновлен под новые эндпоинты.

Удалено:
- Ничего.

---

## 🔒 ПРОД‑МИНИМУМ

- Отключите Swagger в проде: `SWAGGER_ENABLED=false`.
- Оставьте доступ только к `/healthz`, `/readyz`, `/metrics`.

---
