param(
    [Parameter(Mandatory = $true)]
    [string]$BaseUrl,
    [string]$EnvFile = ".env"
)

$ErrorActionPreference = "Stop"

function Assert-ValidBaseUrl {
    param([string]$Url)

    if ([string]::IsNullOrWhiteSpace($Url)) {
        throw "BaseUrl is required. Example: -BaseUrl 'https://api.example.com'"
    }

    if ($Url -match "[<>]" -or $Url -match "real-prod-host" -or $Url -match "your-host") {
        throw "BaseUrl looks like a placeholder ('$Url'). Use a real host, e.g. https://api.example.com"
    }

    $parsed = $null
    if (-not [System.Uri]::TryCreate($Url, [System.UriKind]::Absolute, [ref]$parsed)) {
        throw "BaseUrl is not a valid absolute URI: '$Url'"
    }

    if ($parsed.Scheme -notin @("http", "https")) {
        throw "BaseUrl must use http or https scheme"
    }

    if ([string]::IsNullOrWhiteSpace($parsed.Host)) {
        throw "BaseUrl host is empty"
    }

    $localHosts = @("localhost", "127.0.0.1", "::1")
    if ($localHosts -notcontains $parsed.Host.ToLower()) {
        try {
            Resolve-DnsName $parsed.Host -ErrorAction Stop | Out-Null
        }
        catch {
            throw "DNS resolution failed for '$($parsed.Host)'. Check DNS/ingress and pass a reachable -BaseUrl"
        }
    }
}

function Read-EnvFile {
    param([string]$Path)

    if (-not (Test-Path $Path)) {
        throw "Env file not found: $Path"
    }

    $map = @{}
    Get-Content $Path | ForEach-Object {
        $line = $_.Trim()
        if ($line -eq "" -or $line.StartsWith("#")) { return }
        $idx = $line.IndexOf("=")
        if ($idx -le 0) { return }
        $k = $line.Substring(0, $idx).Trim()
        $v = $line.Substring($idx + 1).Trim()
        $map[$k] = $v
    }
    return $map
}

function Assert-NonPlaceholder {
    param(
        [hashtable]$Env,
        [string]$Key
    )

    if (-not $Env.ContainsKey($Key) -or [string]::IsNullOrWhiteSpace($Env[$Key])) {
        throw "Missing env var: $Key"
    }

    $v = $Env[$Key]
    $bad = @("change_me", "REPLACE_", "your-secret", "admin123")
    foreach ($token in $bad) {
        if ($v -like "*$token*") {
            throw "Unsafe value for ${Key}: contains '$token'"
        }
    }
}

function Check-Url {
    param([string]$Url)

    try {
        $resp = Invoke-RestMethod -Method Get -Uri $Url -TimeoutSec 8
        return $resp
    }
    catch {
        throw "HTTP check failed for ${Url}: $($_.Exception.Message)"
    }
}

Assert-ValidBaseUrl -Url $BaseUrl

Write-Host "== Preflight: env checks ==" -ForegroundColor Cyan
$envMap = Read-EnvFile -Path $EnvFile

Assert-NonPlaceholder -Env $envMap -Key "JWT_SECRET"
Assert-NonPlaceholder -Env $envMap -Key "ADMIN_PASSWORD"
Assert-NonPlaceholder -Env $envMap -Key "POSTGRES_PASSWORD"
Assert-NonPlaceholder -Env $envMap -Key "REDIS_PASSWORD"
Assert-NonPlaceholder -Env $envMap -Key "GRAFANA_ADMIN_PASSWORD"

if (-not $envMap.ContainsKey("SWAGGER_ENABLED") -or $envMap["SWAGGER_ENABLED"].ToLower() -ne "false") {
    throw "SWAGGER_ENABLED must be false in prod"
}

if (-not $envMap.ContainsKey("DATABASE_URL") -or $envMap["DATABASE_URL"] -notmatch "sslmode=require") {
    throw "DATABASE_URL must contain sslmode=require"
}

Write-Host "Env checks: OK" -ForegroundColor Green

Write-Host "== Preflight: health checks ==" -ForegroundColor Cyan
$health = Check-Url -Url "$BaseUrl/healthz"
$ready = Check-Url -Url "$BaseUrl/readyz"

if ($health.status -ne "ok") { throw "healthz status is not ok" }
if ($ready.status -ne "ready") { throw "readyz status is not ready" }

Write-Host "Health checks: OK" -ForegroundColor Green
Write-Host "Preflight OK" -ForegroundColor Green

