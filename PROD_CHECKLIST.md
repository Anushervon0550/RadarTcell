# RadarTcell Production Checklist

## 1) Secrets and environment

- [ ] Use a dedicated prod `.env` outside git.
- [ ] Set strong values for:
  - `JWT_SECRET`
  - `ADMIN_PASSWORD`
  - `POSTGRES_PASSWORD`
  - `REDIS_PASSWORD`
  - `GRAFANA_ADMIN_PASSWORD`
- [ ] Ensure `SWAGGER_ENABLED=false`.
- [ ] Ensure `DATABASE_URL` uses TLS (`sslmode=require` or stricter).

## 2) Infrastructure

- [ ] DNS resolves public API host.
- [ ] Ingress/Firewall exposes only required ports.
- [ ] `/metrics` is reachable only from trusted network.
- [ ] Redis is protected by password/auth.

## 3) Deploy

```powershell
Set-Location "C:\PROJECT\RadarTcell\deploy"
docker compose pull
docker compose up -d --build
```

## 4) Migrations and seed

```powershell
Set-Location "C:\PROJECT\RadarTcell"
powershell -ExecutionPolicy Bypass -File ".\scripts\init-db.ps1" -NoStartDb -ApplySeeds
```

## 5) Health checks

```powershell
Invoke-RestMethod "https://<your-host>/healthz"
Invoke-RestMethod "https://<your-host>/readyz"
```

## 6) Smoke test

```powershell
Set-Location "C:\PROJECT\RadarTcell"
powershell -ExecutionPolicy Bypass -File ".\scripts\smoke.ps1" -BaseUrl "https://<your-host>" -AdminUser "admin" -AdminPassword "<real-password>"
```

## 7) Post-deploy quick checks

- [ ] `go test ./...` is green in CI.
- [ ] Prometheus scrapes expected metrics.
- [ ] Grafana dashboard has fresh data.
- [ ] Login and one read-only public endpoint work.
- [ ] Admin CRUD and cleanup workflow work.

