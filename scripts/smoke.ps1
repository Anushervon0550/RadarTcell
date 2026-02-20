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
        $Headers = $null,
        $BodyObj = $null
    )

    $uri = "$BaseUrl$Path"
    if ($null -ne $BodyObj) {
        $json = $BodyObj | ConvertTo-Json -Depth 10
        return Invoke-RestMethod -Method $Method -Uri $uri -Headers $Headers -ContentType "application/json" -Body $json
    }
    return Invoke-RestMethod -Method $Method -Uri $uri -Headers $Headers
}

Write-Host "1) Health checks..."
Invoke-Api GET "/healthz" | Out-Null
Invoke-Api GET "/readyz"  | Out-Null

Write-Host "2) Login..."
$login = Invoke-Api POST "/api/admin/login" $null @{ username = $AdminUser; password = $AdminPassword }
$token = $login.token
if (-not $token) { throw "No token returned from login" }
$headers = @{ Authorization = "Bearer $token" }

Invoke-Api GET "/api/admin/me" $headers | Out-Null

$suffix = ([guid]::NewGuid().ToString("N")).Substring(0,8)
$trendSlug = "smoke-trend-$suffix"
$tagSlug   = "smoke-tag-$suffix"
$orgSlug   = "smoke-org-$suffix"
$techSlug  = "smoke-tech-$suffix"

Write-Host "3) Create Trend + Tag + Org + Metric..."
Invoke-Api POST "/api/admin/trends" $headers @{ slug=$trendSlug; name="Smoke Trend"; order_index=99 } | Out-Null
Invoke-Api POST "/api/admin/tags" $headers @{ slug=$tagSlug; title="Smoke Tag"; category="Domain"; description="smoke" } | Out-Null
Invoke-Api POST "/api/admin/organizations" $headers @{ slug=$orgSlug; name="Smoke Org"; logo_url="https://example.com/smoke.png" } | Out-Null

$metric = Invoke-Api POST "/api/admin/metrics" $headers @{ name="Smoke Metric $suffix"; type="bar"; description="smoke"; orderable=$true }
$metricId = $metric.id
if (-not $metricId) { throw "No metric id returned" }

Write-Host "4) Create Technology referencing new trend/tag/org..."
Invoke-Api POST "/api/admin/technologies" $headers @{
    slug=$techSlug
    index=10
    name="Smoke Tech"
    trl=5
    trend_slug=$trendSlug
    tag_slugs=@($tagSlug)
    sdg_codes=@("SDG 09")
    organization_slugs=@($orgSlug)
    custom_metric_1=1
    custom_metric_2=2
    custom_metric_3=3
    custom_metric_4=4
} | Out-Null

Write-Host "5) Verify public endpoints..."
Invoke-Api GET "/api/trends" | Out-Null
Invoke-Api GET "/api/tags" | Out-Null
Invoke-Api GET "/api/organizations" | Out-Null
Invoke-Api GET "/api/metrics" | Out-Null
Invoke-Api GET "/api/technologies/$techSlug" | Out-Null

Write-Host "6) Cleanup (delete tech -> tag/trend/org -> metric)..."
Invoke-Api DELETE "/api/admin/technologies/$techSlug" $headers | Out-Null
Invoke-Api DELETE "/api/admin/tags/$tagSlug" $headers | Out-Null
Invoke-Api DELETE "/api/admin/trends/$trendSlug" $headers | Out-Null
Invoke-Api DELETE "/api/admin/organizations/$orgSlug" $headers | Out-Null
Invoke-Api DELETE "/api/admin/metrics/$metricId" $headers | Out-Null

Write-Host "SMOKE OK ✅"