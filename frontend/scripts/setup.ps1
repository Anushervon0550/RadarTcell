#requires -Version 5.1
# ---------------------------------------------------------------------------
# Bootstrap for the frontend workspace behind a corporate HTTP proxy that
# performs TLS interception (Squid + MITM CA).
#
# What it does:
#   1. Refreshes PATH so a freshly installed Node is visible.
#   2. Detects the system proxy from IE settings.
#   3. Exports the Windows Root CA store to a PEM bundle Node can trust.
#   4. Configures npm + pnpm with proxy, timeouts, and CA file.
#   5. Installs pnpm via npm (Corepack ignores env proxy on Node 20+).
#   6. Runs pnpm install for the workspace.
# ---------------------------------------------------------------------------

function Fail([string]$msg) { Write-Host $msg -ForegroundColor Red; exit 1 }

# 1. Refresh PATH
$env:Path = [Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [Environment]::GetEnvironmentVariable('Path','User')

# 2. Proxy detection
$proxy = $env:HTTPS_PROXY
if (-not $proxy) { $proxy = $env:HTTP_PROXY }
if (-not $proxy) {
    try {
        $inet = Get-ItemProperty 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings' -ErrorAction SilentlyContinue
        if ($inet -and $inet.ProxyEnable -eq 1 -and $inet.ProxyServer) {
            $server = [string]$inet.ProxyServer
            if ($server -notmatch '^https?://') { $server = "http://$server" }
            $proxy = $server
        }
    } catch {}
}
if ($proxy) {
    Write-Host "Proxy: $proxy" -ForegroundColor DarkGray
    $env:HTTP_PROXY  = $proxy; $env:HTTPS_PROXY = $proxy
    $env:http_proxy  = $proxy; $env:https_proxy = $proxy
    $env:NO_PROXY    = 'localhost,127.0.0.1,::1'
    $env:no_proxy    = $env:NO_PROXY
}

Write-Host ('node: ' + (& node --version))
Write-Host ('npm:  ' + (& npm --version))

# 3. Export Windows Root CA store to a PEM file
Write-Host ''
Write-Host '=== Exporting Windows Root CA store to PEM ===' -ForegroundColor Cyan
$caFile = Join-Path $env:USERPROFILE '.corp-root-ca.pem'
$caExported = $false
try {
    $lines = New-Object System.Collections.Generic.List[string]
    $stores = @('Cert:\LocalMachine\Root', 'Cert:\CurrentUser\Root',
                'Cert:\LocalMachine\CA',   'Cert:\CurrentUser\CA')
    foreach ($storePath in $stores) {
        if (Test-Path $storePath) {
            Get-ChildItem $storePath -ErrorAction SilentlyContinue | ForEach-Object {
                $b64 = [Convert]::ToBase64String($_.RawData)
                $lines.Add('-----BEGIN CERTIFICATE-----')
                for ($i = 0; $i -lt $b64.Length; $i += 64) {
                    $lines.Add($b64.Substring($i, [Math]::Min(64, $b64.Length - $i)))
                }
                $lines.Add('-----END CERTIFICATE-----')
            }
        }
    }
    Set-Content -Path $caFile -Value $lines -Encoding ASCII
    $count = ($lines | Where-Object { $_ -eq '-----BEGIN CERTIFICATE-----' }).Count
    Write-Host "CA bundle: $caFile ($count certs)"
    $env:NODE_EXTRA_CA_CERTS = $caFile
    $caExported = $true
} catch {
    Write-Host ('CA export failed: ' + $_.Exception.Message) -ForegroundColor Yellow
}

# 4. Configure npm
if ($proxy) {
    npm config set proxy       $proxy --location=user 2>$null
    npm config set https-proxy $proxy --location=user 2>$null
}
npm config set fetch-retries          5      --location=user 2>$null
npm config set fetch-retry-mintimeout 10000  --location=user 2>$null
npm config set fetch-retry-maxtimeout 60000  --location=user 2>$null
npm config set fetch-timeout          600000 --location=user 2>$null
if ($caExported) {
    npm config set cafile $caFile --location=user 2>$null
}
# Fallback for MITM environments where the CA export missed the right root.
npm config set strict-ssl false --location=user 2>$null

# 5. Install pnpm via npm and resolve pnpm.cmd via `npm prefix -g`.
$npmPrefix = ((& npm prefix -g 2>$null) | Out-String).Trim()
if (-not $npmPrefix) { Fail 'npm prefix -g returned nothing' }
$pnpmCmd = Join-Path $npmPrefix 'pnpm.cmd'

if (-not (Test-Path $pnpmCmd)) {
    Write-Host ''
    Write-Host '=== Installing pnpm via npm ===' -ForegroundColor Cyan
    # --force: overwrite any leftover Corepack shims from previous attempts.
    npm install -g --force --loglevel=http --no-audit --no-fund pnpm@9.15.0 2>&1
    if ($LASTEXITCODE -ne 0) { Fail "npm install -g pnpm failed with exit $LASTEXITCODE" }
}
if (-not (Test-Path $pnpmCmd)) { Fail "pnpm.cmd not found at $pnpmCmd" }

$env:Path = "$npmPrefix;${env:Path}"
Write-Host ('pnpm: ' + (& $pnpmCmd --version))

# 6. Configure pnpm
if ($proxy) {
    & $pnpmCmd config set proxy       $proxy --location=user 2>$null | Out-Null
    & $pnpmCmd config set https-proxy $proxy --location=user 2>$null | Out-Null
}
& $pnpmCmd config set fetch-retries          5      --location=user 2>$null | Out-Null
& $pnpmCmd config set fetch-retry-mintimeout 10000  --location=user 2>$null | Out-Null
& $pnpmCmd config set fetch-retry-maxtimeout 60000  --location=user 2>$null | Out-Null
& $pnpmCmd config set network-timeout        600000 --location=user 2>$null | Out-Null
if ($caExported) {
    & $pnpmCmd config set cafile $caFile --location=user 2>$null | Out-Null
}
& $pnpmCmd config set strict-ssl false --location=user 2>$null | Out-Null

# 7. Install workspace
Set-Location (Split-Path -Parent $PSScriptRoot)
Write-Host ''
Write-Host '=== pnpm install (workspace) ===' -ForegroundColor Cyan
& $pnpmCmd install --no-frozen-lockfile
if ($LASTEXITCODE -ne 0) { Fail "pnpm install failed with exit $LASTEXITCODE" }

Write-Host ''
Write-Host 'SUCCESS. Run:' -ForegroundColor Green
Write-Host '  pnpm dev:public   -> http://localhost:5173'
Write-Host '  pnpm dev:admin    -> http://localhost:5174'
