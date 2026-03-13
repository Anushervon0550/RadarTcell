# RadarTcell Final Handoff

## Status

Project is in a shippable state for the current scope:
- Dynamic metrics are supported (`custom_metrics`) in public and admin responses.
- API version aliases are active (`/api` + `/api/v1`, `/api/admin` + `/api/v1/admin`).
- Soft delete is implemented for technologies and catalog/admin entities.
- Admin technology restore endpoint is available.
- Cursor pagination is available for public technologies list.

## Migrations To Apply

Apply migrations in sequence (already ordered by filename):
- `000010_create_technology_metric_values.up.sql`
- `000011_soft_delete_technologies.up.sql`
- `000012_soft_delete_technologies_slug_index.up.sql`
- `000013_soft_delete_catalog_entities.up.sql`

## Verified Locally

- Full test suite passes (`go test ./...`).
- Swagger docs regenerated with `swag init` into `docs/`.

## Final Release Commands (PowerShell)

```powershell
Set-Location "C:\Users\user\TCELL\RadarTcell"
go test ./...

Set-Location "C:\Users\user\TCELL\RadarTcell"
swag init -g cmd/api/main.go -o docs

Set-Location "C:\Users\user\TCELL\RadarTcell"
powershell -ExecutionPolicy Bypass -File ".\scripts\init-db.ps1" -NoStartDb -ApplySeeds

Set-Location "C:\Users\user\TCELL\RadarTcell\deploy"
docker compose up -d --build
```

## Smoke Checklist

1. Health:
   - `GET /healthz`
   - `GET /readyz`
2. Public technologies:
   - `GET /api/technologies`
   - `GET /api/v1/technologies`
   - `GET /api/technologies?cursor=<index:id>`
3. Admin auth:
   - `POST /api/admin/login`
4. Admin technologies soft delete flow:
   - `DELETE /api/admin/technologies/{slug}`
   - `PUT /api/admin/technologies/{slug}/restore`
   - `GET /api/admin/technologies?include_deleted=true`
5. Dynamic metric value:
   - `GET /api/metrics/{id}/values?technology_id=<uuid>`

## Notes

- Legacy `custom_metric_1..4` fields are still present for backward compatibility.
- New development should prefer `custom_metrics`.
- In production keep `DATABASE_URL` with TLS (`sslmode=require` or stricter).

