param(
    [string]$ComposeFile = "deploy/docker-compose.yml",
    [string]$MigrationDir = "migrations",
    [string]$DatabaseUrl = "postgres://radar_tcell:radar_tcell_password@localhost:15433/radar_tcell?sslmode=disable",
    [switch]$ApplySeeds,
    [switch]$ForceSeed,
    [int]$WaitTimeoutSec = 90,
    [switch]$NoStartDb,
    [switch]$DryRun
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"
$RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path

function Resolve-FromRepoRoot {
    param([string]$RelativePath)

    return (Resolve-Path (Join-Path $RepoRoot $RelativePath)).Path
}

function Invoke-Step {
    param(
        [string]$Title,
        [scriptblock]$Action
    )

    Write-Host "`n== $Title ==" -ForegroundColor Cyan
    if ($DryRun) {
        Write-Host "[dry-run] step skipped" -ForegroundColor Yellow
        return
    }

    & $Action
}

function Invoke-Native {
    param(
        [string]$Command,
        [string[]]$CommandArgs = @()
    )

    & $Command @CommandArgs
    if ($LASTEXITCODE -ne 0) {
        throw "Command failed ($LASTEXITCODE): $Command $($CommandArgs -join ' ')"
    }
}

function Get-ComposeCommand {
    if (Get-Command docker -ErrorAction SilentlyContinue) {
        try {
            docker compose version | Out-Null
            return @("docker", "compose")
        }
        catch {}
    }

    if (Get-Command docker-compose -ErrorAction SilentlyContinue) {
        return @("docker-compose")
    }

    throw "docker compose not found. Install Docker Desktop (docker compose) or docker-compose."
}

function Wait-PostgresReady {
    param(
        [string]$ContainerName,
        [int]$TimeoutSec
    )

    $deadline = (Get-Date).AddSeconds($TimeoutSec)
    while ((Get-Date) -lt $deadline) {
        $status = ""
        try {
            $status = (docker inspect -f "{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}" $ContainerName 2>$null).Trim()
        }
        catch {
            $status = ""
        }

        if ($status -eq "healthy" -or $status -eq "running") {
            try {
                docker exec $ContainerName pg_isready -U radar_tcell -d radar_tcell | Out-Null
                Write-Host "Postgres is ready." -ForegroundColor Green
                return
            }
            catch {}
        }

        Start-Sleep -Seconds 2
    }

    throw "Postgres did not become ready in $TimeoutSec sec. Check: docker logs $ContainerName"
}

$composePath = Resolve-FromRepoRoot $ComposeFile
$migrationPath = Resolve-FromRepoRoot $MigrationDir

if (-not (Test-Path $composePath)) {
    throw "Compose file not found: $ComposeFile"
}
if (-not (Test-Path $migrationPath)) {
    throw "Migration directory not found: $MigrationDir"
}

$containerName = "radartcell-postgres"

Push-Location $RepoRoot
try {
    if (-not $NoStartDb) {
        $composeCmd = Get-ComposeCommand

        Invoke-Step "Start Postgres" {
            if ($composeCmd.Count -eq 2) {
                Invoke-Native -Command $composeCmd[0] -CommandArgs @($composeCmd[1], "-f", $composePath, "up", "-d", "postgres")
            }
            else {
                Invoke-Native -Command $composeCmd[0] -CommandArgs @("-f", $composePath, "up", "-d", "postgres")
            }
        }

        Invoke-Step "Wait for Postgres" {
            Wait-PostgresReady -ContainerName $containerName -TimeoutSec $WaitTimeoutSec
        }
    }

    Invoke-Step "Apply migrations (up)" {
        Invoke-Native -Command "migrate" -CommandArgs @("-path", $MigrationDir, "-database", $DatabaseUrl, "up")
    }

    if ($ApplySeeds) {
        $seedFiles = Get-ChildItem -Path $migrationPath -Filter "seed_*.sql" | Sort-Object Name
        if ($seedFiles.Count -eq 0) {
            Write-Host "No seed_*.sql files found. Skipping." -ForegroundColor Yellow
        }
        else {
            # Check if data already exists (technologies table has rows)
            $count = 0
            try {
                $raw = & psql $DatabaseUrl -t -c "SELECT COUNT(*) FROM technologies;" 2>$null
                $countStr = ($raw | Where-Object { $_ -match '\S' } | Select-Object -First 1).Trim()
                if ($countStr -match '^\d+$') { $count = [int]$countStr }
            } catch {}

            if (-not $ForceSeed -and $count -gt 0) {
                Write-Host "`n== Seeds ==" -ForegroundColor Cyan
                Write-Host "Data already exists ($count rows in technologies). Skipping seeds. Use -ForceSeed to override." -ForegroundColor Yellow            }
            else {
                foreach ($seed in $seedFiles) {
                    Invoke-Step "Apply seed: $($seed.Name)" {
                        Invoke-Native -Command "psql" -CommandArgs @($DatabaseUrl, "-f", $seed.FullName)
                    }
                }
            }
        }
    }
}
finally {
    Pop-Location
}

Write-Host "`nDone." -ForegroundColor Green
