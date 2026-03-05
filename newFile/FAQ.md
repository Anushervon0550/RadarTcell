# ❓ FAQ: RadarTcell

## ✅ Частые вопросы

**Q: Как запустить проект?**

```powershell
cd C:\PROJECT\RadarTcell\deploy
docker compose up -d

cd C:\PROJECT\RadarTcell
migrate -path "./migrations" -database $env:DATABASE_URL up

go run .\cmd\api\main.go
```

**Q: Где Swagger?**

`http://localhost:8080/swagger/`

**Q: Почему 401?**

Нужен JWT: `POST /api/admin/login`, затем `Authorization: Bearer <token>`.

**Q: Почему 404?**

Проверьте slug/code и префикс `/api` или `/api/admin`.

**Q: Как включить локализацию?**

Добавьте `?locale=ru` к публичным спискам/технологиям.

**Q: Где метрики?**

`http://localhost:8080/metrics`.

**Q: Какие админские эндпоинты?**

Все начинаются с `/api/admin`.

---

## ✅ Ключевые изменения

- Админские эндпоинты: `/api/admin/*`.
- Публичные списки/технологии поддерживают `?locale=ru`.
- Есть admin i18n API: `/api/admin/i18n/*`.
- Системные: `/healthz`, `/readyz`, `/metrics`, `/openapi.yaml`, `/swagger/`.
- Миграции: `000006` (field_key), `000007` (i18n).
