$ErrorActionPreference = 'Continue'
& go test ./internal/... -count=1
exit $LASTEXITCODE
