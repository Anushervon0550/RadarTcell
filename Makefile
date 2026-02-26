.PHONY: test run smoke

test:
	go test ./cmd/... ./internal/...

run:
	go run ./cmd/api

smoke:
	powershell -ExecutionPolicy Bypass -File .\scripts\smoke.ps1