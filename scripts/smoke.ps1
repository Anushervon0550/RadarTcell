param(
    [string]$BaseUrl = "http://localhost:8080",
    [string]$AdminUser = "admin",
    [string]$AdminPassword = "admin123"
)

$ErrorActionPreference = "Stop"

function Invoke-Api {
    param(
        [string]$Method,
        [string]$Path,
        [hashtable]$Headers = $null,
        $BodyObj = $null
    )

    $uri = "$BaseUrl$Path"

    if ($null -ne $BodyObj) {
        $json = $BodyObj | ConvertTo-Json -Depth 10
        return Invoke-RestMethod -Method $Method -Uri $uri -Headers $Headers -ContentType "application/json" -Body $json
    }

    return Invoke-RestMethod -Method $Method -Uri $uri -Headers $Headers
}

Write-Host "== Smoke test started ==" -ForegroundColor Cyan

# 1) health
Write-Host "1) Health checks..."
$h = Invoke-Api GET "/healthz"
$r = Invoke-Api GET "/readyz"
if ($h.status -ne "ok") { throw "healthz failed" }
if ($r.status -ne "ready") { throw "readyz failed" }

# 2) login
Write-Host "2) Admin login..."
$login = Invoke-Api POST "/api/admin/login" $null @{
    username = $AdminUser
    password = $AdminPassword
}
if (-not $login.token) { throw "No token returned" }

$headers = @{ Authorization = "Bearer $($login.token)" }

$me = Invoke-Api GET "/api/admin/me" $headers
if ($me.role -ne "admin") { throw "admin/me failed" }

# unique suffix
$suffix = ([guid]::NewGuid().ToString("N")).Substring(0, 8)

$trendSlug = "smoke-trend-$suffix"
$tagSlug   = "smoke-tag-$suffix"
$orgSlug   = "smoke-org-$suffix"
$techSlug  = "smoke-tech-$suffix"

# 3) create catalog entities
Write-Host "3) Create trend/tag/org/metric..."
Invoke-Api POST "/api/admin/trends" $headers @{
    slug = $trendSlug
    name = "Smoke Trend"
    order_index = 99
} | Out-Null

Invoke-Api POST "/api/admin/tags" $headers @{
    slug = $tagSlug
    title = "Smoke Tag"
    category = "Domain"
    description = "smoke"
} | Out-Null

Invoke-Api POST "/api/admin/organizations" $headers @{
    slug = $orgSlug
    name = "Smoke Org"
    logo_url = "https://example.com/smoke.png"
} | Out-Null

$metric = Invoke-Api POST "/api/admin/metrics" $headers @{
    name = "Smoke Metric $suffix"
    type = "bar"
    description = "smoke"
    orderable = $true
}
if (-not $metric.id) { throw "Metric create failed: no id" }
$metricId = $metric.id

# 4) create technology linked to all of them
Write-Host "4) Create technology..."
Invoke-Api POST "/api/admin/technologies" $headers @{
    slug = $techSlug
    index = 99
    name = "Smoke Tech"
    trl = 5
    trend_slug = $trendSlug
    tag_slugs = @($tagSlug)
    sdg_codes = @("SDG 09")
    organization_slugs = @($orgSlug)
    custom_metric_1 = 1
    custom_metric_2 = 2
    custom_metric_3 = 3
    custom_metric_4 = 4
} | Out-Null

# 5) public checks
Write-Host "5) Public API checks..."
$tech = Invoke-Api GET "/api/technologies/$techSlug"
if ($tech.slug -ne $techSlug) { throw "GET technology by slug failed" }

$list = Invoke-Api GET "/api/technologies?search=Smoke"
if ($list.total -lt 1) { throw "Technology search failed" }

Invoke-Api GET "/api/trends" | Out-Null
Invoke-Api GET "/api/tags" | Out-Null
Invoke-Api GET "/api/organizations" | Out-Null
Invoke-Api GET "/api/metrics" | Out-Null

# 6) cleanup (important order)
Write-Host "6) Cleanup..."
Invoke-Api DELETE "/api/admin/technologies/$techSlug" $headers | Out-Null
Invoke-Api DELETE "/api/admin/tags/$tagSlug" $headers | Out-Null
Invoke-Api DELETE "/api/admin/trends/$trendSlug" $headers | Out-Null
Invoke-Api DELETE "/api/admin/organizations/$orgSlug" $headers | Out-Null
Invoke-Api DELETE "/api/admin/metrics/$metricId" $headers | Out-Null

Write-Host "SMOKE OK ✅" -ForegroundColor Green