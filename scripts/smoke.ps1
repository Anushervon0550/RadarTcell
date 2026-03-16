param(
    [string]$BaseUrl = "http://localhost:8080",
    [string]$AdminUser = "admin",
    [string]$AdminPassword = "admin123",
    [string]$TrustedOrigin = "http://localhost:3000",
    [int]$RequestTimeoutSec = 15,
    [int]$HealthRetries = 6,
    [int]$HealthRetryDelaySec = 2
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

Assert-ValidBaseUrl -Url $BaseUrl
$BaseUri = [System.Uri]$BaseUrl

function Test-TcpPortReachable {
    param(
        [string]$HostName,
        [int]$Port,
        [int]$TimeoutMs = 1500
    )

    $client = New-Object System.Net.Sockets.TcpClient
    try {
        $iar = $client.BeginConnect($HostName, $Port, $null, $null)
        if (-not $iar.AsyncWaitHandle.WaitOne($TimeoutMs, $false)) {
            return $false
        }
        $client.EndConnect($iar) | Out-Null
        return $true
    }
    catch {
        return $false
    }
    finally {
        $client.Close()
    }
}

function Add-CsrfHeaders {
    param(
        [string]$Method,
        [hashtable]$Headers
    )

    $effective = @{}
    if ($Headers) {
        foreach ($key in $Headers.Keys) {
            $effective[$key] = $Headers[$key]
        }
    }

    $m = $Method.ToUpperInvariant()
    if ($m -in @("POST", "PUT", "PATCH", "DELETE")) {
        if (-not $effective.ContainsKey("Origin")) {
            $effective["Origin"] = $TrustedOrigin
        }
        if (-not $effective.ContainsKey("Referer")) {
            $effective["Referer"] = "$TrustedOrigin/"
        }
    }

    return $effective
}

function Invoke-Api {
    param(
        [string]$Method,
        [string]$Path,
        [hashtable]$Headers = $null,
        $BodyObj = $null
    )

    $uri = "$BaseUrl$Path"
    $Headers = Add-CsrfHeaders -Method $Method -Headers $Headers

    if ($null -ne $BodyObj) {
        $json = $BodyObj | ConvertTo-Json -Depth 10
        return Invoke-RestMethod -Method $Method -Uri $uri -Headers $Headers -ContentType "application/json" -Body $json -TimeoutSec $RequestTimeoutSec
    }

    return Invoke-RestMethod -Method $Method -Uri $uri -Headers $Headers -TimeoutSec $RequestTimeoutSec
}

function Invoke-ApiExpectError {
    param(
        [string]$Method,
        [string]$Path,
        [hashtable]$Headers = $null,
        $BodyObj = $null,
        [int[]]$ExpectedStatus = @()
    )

    $uri = "$BaseUrl$Path"
    $Headers = Add-CsrfHeaders -Method $Method -Headers $Headers

    try {
        if ($null -ne $BodyObj) {
            $json = $BodyObj | ConvertTo-Json -Depth 10
            Invoke-WebRequest -Method $Method -Uri $uri -Headers $Headers -ContentType "application/json" -Body $json -TimeoutSec $RequestTimeoutSec -ErrorAction Stop | Out-Null
        }
        else {
            Invoke-WebRequest -Method $Method -Uri $uri -Headers $Headers -TimeoutSec $RequestTimeoutSec -ErrorAction Stop | Out-Null
        }

        throw "Expected error for $Method $Path, but request succeeded"
    }
    catch {
        $status = $null

        try {
            if ($_.Exception.Response -and $_.Exception.Response.StatusCode) {
                $status = [int]$_.Exception.Response.StatusCode
            }
        } catch {}

        if ($null -eq $status) {
            try {
                $status = [int]$_.Exception.Response.StatusCode.value__
            } catch {}
        }

        $body = $null
        if ($_.ErrorDetails -and $_.ErrorDetails.Message) {
            $body = $_.ErrorDetails.Message
        }
        else {
            $body = $_.Exception.Message
        }

        if ($null -eq $status) {
            throw "Failed to read HTTP status for $Method $Path. Raw error: $body"
        }

        if ($ExpectedStatus.Count -gt 0 -and ($ExpectedStatus -notcontains $status)) {
            throw "Unexpected status for $Method $Path. Expected: $($ExpectedStatus -join ', '), got: $status. Body: $body"
        }

        return @{
            status = $status
            body   = $body
        }
    }
}

function Invoke-WithRetry {
    param(
        [scriptblock]$Action,
        [string]$OperationName,
        [int]$MaxAttempts,
        [int]$DelaySeconds
    )

    for ($attempt = 1; $attempt -le $MaxAttempts; $attempt++) {
        try {
            return & $Action
        }
        catch {
            $lastError = $_.Exception.Message
            if ($attempt -eq $MaxAttempts) {
                throw "$OperationName failed after $MaxAttempts attempts. Last error: $lastError"
            }

            Write-Host "$OperationName attempt $attempt/$MaxAttempts failed, retry in $DelaySeconds sec..."
            Start-Sleep -Seconds $DelaySeconds
        }
    }
}

Write-Host "== Smoke test started ==" -ForegroundColor Cyan

# cleanup vars (pre-init for safety)
$headers = $null
$TmpMetricId = $null
$TmpMetricCreated = $false
$metricId = $null

# ✅ SDG smoke vars
$sdgCode = $null
$sdgCreated = $false

# unique suffix
$suffix = ([guid]::NewGuid().ToString("N")).Substring(0, 8)

$trendSlug = "smoke-trend-$suffix"
$tagSlug   = "smoke-tag-$suffix"
$orgSlug   = "smoke-org-$suffix"
$techSlug  = "smoke-tech-$suffix"

# ✅ SDG code without spaces (URL-friendly)
$sdgCode = "SDG-99-$suffix"

try {
    # 1) health
    Write-Host "1) Health checks..."
    $port = if ($BaseUri.IsDefaultPort) { if ($BaseUri.Scheme -eq "https") { 443 } else { 80 } } else { $BaseUri.Port }
    if (-not (Test-TcpPortReachable -HostName $BaseUri.Host -Port $port)) {
        throw "Cannot connect to $($BaseUri.Host):$port. Start API first (example: go run ./cmd/api) or pass reachable -BaseUrl."
    }

    $h = Invoke-WithRetry -OperationName "GET /healthz" -MaxAttempts $HealthRetries -DelaySeconds $HealthRetryDelaySec -Action {
        Invoke-Api GET "/healthz"
    }
    $r = Invoke-WithRetry -OperationName "GET /readyz" -MaxAttempts $HealthRetries -DelaySeconds $HealthRetryDelaySec -Action {
        Invoke-Api GET "/readyz"
    }
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

    # 3) create catalog entities
    Write-Host "3) Create trend/tag/org/metric/SDG..."
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

    # ✅ create SDG
    $sdg = Invoke-Api POST "/api/admin/sdgs" $headers @{
        code = $sdgCode
        title = "Smoke SDG"
        description = "smoke"
        icon = "x"
    }
    if (-not $sdg.code -or $sdg.code -ne $sdgCode) {
        throw "SDG create failed"
    }
    $sdgCreated = $true

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
    $all = Invoke-Api GET "/api/technologies?limit=200"
    if ($all.total -lt 16) { throw "Expected at least 16 technologies after seed_0002" }

    $sdg09 = Invoke-Api GET "/api/sdgs/SDG%2009/technologies?limit=200"
    if ($sdg09.total -lt 10) { throw "Expected SDG 09 total >= 10" }

    # 5.1) Metric values checks (field_key)
    Write-Host "5.1) Metric values checks..."

    # Берем одну технологию
    $techs = Invoke-Api GET "/api/technologies?limit=1"
    if (-not $techs.items -or $techs.items.Count -lt 1) {
        throw "No technologies returned from /api/technologies"
    }
    $techId = $techs.items[0].id

    # Берем seed-метрику Custom Metric 01
    $metrics = Invoke-Api GET "/api/metrics"
    $customMetric = $metrics | Where-Object { $_.name -eq "Custom Metric 01" } | Select-Object -First 1
    if (-not $customMetric) {
        throw "Custom Metric 01 not found in /api/metrics"
    }

    # Проверяем endpoint values для custom_metric_1
    $mv1 = Invoke-Api GET "/api/metrics/$($customMetric.id)/values?technology_id=$techId"
    if (-not $mv1.metric_id -or $mv1.metric_id -ne $customMetric.id) {
        throw "Metric values check failed for Custom Metric 01 (metric_id mismatch)"
    }
    if ($mv1.field_key -ne "custom_metric_1") {
        throw "Metric values check failed for Custom Metric 01 (field_key expected custom_metric_1, got '$($mv1.field_key)')"
    }
    if ($null -eq $mv1.value) {
        throw "Metric values check failed for Custom Metric 01 (value is null)"
    }

    # Пытаемся использовать уже существующую list_index метрику
    $listIndexMetric = $metrics | Where-Object { $_.name -eq "List Index Metric API" } | Select-Object -First 1

    if (-not $listIndexMetric) {
        # Если нет - создаем временную
        $TmpMetricName = "Smoke List Index Metric $([guid]::NewGuid().ToString('N').Substring(0,8))"
        $tmpMetric = Invoke-Api POST "/api/admin/metrics" $headers @{
            name = $TmpMetricName
            type = "distance"
            description = "smoke list_index"
            orderable = $true
            field_key = "list_index"
        }

        if (-not $tmpMetric.id) {
            throw "Failed to create temp metric for list_index smoke check"
        }

        $TmpMetricId = $tmpMetric.id
        $TmpMetricCreated = $true
        $listIndexMetricId = $TmpMetricId
    }
    else {
        $listIndexMetricId = $listIndexMetric.id
    }

    # Проверяем endpoint values для list_index
    $mv2 = Invoke-Api GET "/api/metrics/$listIndexMetricId/values?technology_id=$techId"
    if ($mv2.field_key -ne "list_index") {
        throw "Metric values check failed for list_index metric (field_key mismatch)"
    }
    if ($null -eq $mv2.value) {
        throw "Metric values check failed for list_index metric (value is null)"
    }

    Write-Host "Metric values checks OK"

    # 5.2) Relation endpoints checks...
    Write-Host "5.2) Relation endpoints checks..."

    # trends/{slug}/technologies for created smoke trend
    $trendTechs = Invoke-Api GET "/api/trends/$trendSlug/technologies"
    if ($null -eq $trendTechs.total -or $trendTechs.total -lt 1) {
        throw "Trend technologies endpoint failed for $trendSlug"
    }

    # tags/{slug}/technologies for created smoke tag
    $tagTechs = Invoke-Api GET "/api/tags/$tagSlug/technologies"
    if ($null -eq $tagTechs.total -or $tagTechs.total -lt 1) {
        throw "Tag technologies endpoint failed for $tagSlug"
    }

    # organizations/{slug} for created smoke org
    $org = Invoke-Api GET "/api/organizations/$orgSlug"
    if ($org.slug -ne $orgSlug) {
        throw "Organization by slug endpoint failed for $orgSlug"
    }

    # SDG endpoints (ожидаем, что seed уже привязан)
    $sdg09 = Invoke-Api GET "/api/sdgs/SDG%2009/technologies"
    if ($null -eq $sdg09.total -or $sdg09.total -lt 1) {
        throw "SDG 09 technologies endpoint failed (expected total >= 1)"
    }

    $sdg03 = Invoke-Api GET "/api/sdgs/SDG%2003/technologies"
    if ($null -eq $sdg03.total -or $sdg03.total -lt 1) {
        throw "SDG 03 technologies endpoint failed (expected total >= 1)"
    }

    # highlight filters
    $hlTag = Invoke-Api GET "/api/technologies?highlight=tag:ml"
    if ($null -eq $hlTag.total) {
        throw "Highlight tag endpoint failed"
    }

    $hlCombo = Invoke-Api GET "/api/technologies?highlight=trend:ai&highlight=organization:openai"
    if ($null -eq $hlCombo.total) {
        throw "Highlight combo endpoint failed"
    }

    Write-Host "Relation endpoints checks OK"

    # ✅ 5.2.1) Admin SDG checks (create/update/delete)
    Write-Host "5.2.1) Admin SDG checks..."
    $sdgUpd = Invoke-Api PUT "/api/admin/sdgs/$sdgCode" $headers @{
        title = "Smoke SDG Updated"
        description = "smoke2"
        icon = "y"
    }
    if (-not $sdgUpd.code -or $sdgUpd.code -ne $sdgCode) {
        throw "SDG update failed"
    }

    # 5.x) Technologies validation checks (should return 400)
    Write-Host "5.x) Technologies validation checks..."

    $bad1 = Invoke-ApiExpectError GET "/api/technologies?page=0" $null $null @(400)
    if ($bad1.body -notmatch "page must be") { throw "Expected page validation error" }

    $bad2 = Invoke-ApiExpectError GET "/api/technologies?limit=0" $null $null @(400)
    if ($bad2.body -notmatch "limit must be between") { throw "Expected limit validation error" }

    $bad3 = Invoke-ApiExpectError GET "/api/technologies?sort_by=hack" $null $null @(400)
    if ($bad3.body -notmatch "sort_by") { throw "Expected sort_by validation error" }

    $bad4 = Invoke-ApiExpectError GET "/api/technologies?trl_min=8&trl_max=2" $null $null @(400)
    if ($bad4.body -notmatch "trl_min must be") { throw "Expected TRL range validation error" }

    Write-Host "Technologies validation checks OK"

    # 5.3) Negative checks...
    Write-Host "5.3) Negative checks..."

    # Public: not found (technology by slug)
    $e1 = Invoke-ApiExpectError GET "/api/technologies/not-exists-12345" $null $null @(404)
    if ($e1.body -notmatch "not found") {
        throw "Expected 'not found' message for unknown technology slug"
    }

    # Public: metric values without required query param
    $e2 = Invoke-ApiExpectError GET "/api/metrics/$metricId/values" $null $null @(400)
    if ($e2.body -notmatch "technology_id is required") {
        throw "Expected 'technology_id is required' for metric values without technology_id"
    }

    # Admin: unauthorized without token
    $e3 = Invoke-ApiExpectError POST "/api/admin/trends" $null @{
        slug = "unauth-test"
        name = "Unauthorized"
        order_index = 1
    } @(401)
    if ($e3.body -notmatch "missing bearer token") {
        throw "Expected 'missing bearer token' for admin route without token"
    }

    # Admin: validation error (invalid TRL)
    $e4 = Invoke-ApiExpectError POST "/api/admin/technologies" $headers @{
        slug = "bad-tech-$suffix"
        index = 10
        name = "Bad Tech"
        trl = 12
        trend_slug = $trendSlug
    } @(400)
    if ($e4.body -notmatch "trl must be 1..9") {
        throw "Expected TRL validation error"
    }

    # Admin: validation error (invalid index)
    $e5 = Invoke-ApiExpectError POST "/api/admin/technologies" $headers @{
        slug = "bad-tech2-$suffix"
        index = 0
        name = "Bad Tech 2"
        trl = 5
        trend_slug = $trendSlug
    } @(400)
    if ($e5.body -notmatch "index must be between 1 and 99") {
        throw "Expected index validation error"
    }

    # Admin: delete non-existing trend
    $e6 = Invoke-ApiExpectError DELETE "/api/admin/trends/not-exists-$suffix" $headers $null @(404)
    if ($e6.body -notmatch "not found") {
        throw "Expected not found on deleting non-existing trend"
    }

    Write-Host "Negative checks OK"

    # Остальные public checks
    $tech = Invoke-Api GET "/api/technologies/$techSlug"
    if ($tech.slug -ne $techSlug) { throw "GET technology by slug failed" }

    $list = Invoke-Api GET "/api/technologies?search=Smoke"
    if ($list.total -lt 1) { throw "Technology search failed" }

    Invoke-Api GET "/api/trends" | Out-Null
    Invoke-Api GET "/api/tags" | Out-Null
    Invoke-Api GET "/api/organizations" | Out-Null
    Invoke-Api GET "/api/metrics" | Out-Null

    Write-Host "SMOKE OK ✅" -ForegroundColor Green
}
finally {
    # 6) cleanup (important order)
    Write-Host "6) Cleanup..."

    # delete smoke technology first (depends on trend/tag/org)
    try {
        if ($headers) { Invoke-Api DELETE "/api/admin/technologies/$techSlug" $headers | Out-Null }
    } catch {
        Write-Host "Cleanup warning: failed to delete tech $techSlug"
    }

    # ✅ delete SDG created by smoke (ignore if already deleted)
    if ($sdgCreated -and $sdgCode) {
        try {
            if ($headers) { Invoke-Api DELETE "/api/admin/sdgs/$sdgCode" $headers | Out-Null }
        } catch {
            Write-Host "Cleanup warning: failed to delete sdg $sdgCode"
        }
    }

    # delete temp list_index metric only if created in this smoke
    if ($TmpMetricCreated -and $TmpMetricId) {
        try {
            if ($headers) { Invoke-Api DELETE "/api/admin/metrics/$TmpMetricId" $headers | Out-Null }
        } catch {
            Write-Host "Cleanup warning: failed to delete temp metric $TmpMetricId"
        }
    }

    # delete smoke metric (bar)
    if ($metricId) {
        try {
            if ($headers) { Invoke-Api DELETE "/api/admin/metrics/$metricId" $headers | Out-Null }
        } catch {
            Write-Host "Cleanup warning: failed to delete metric $metricId"
        }
    }

    # delete catalog entities
    try {
        if ($headers) { Invoke-Api DELETE "/api/admin/tags/$tagSlug" $headers | Out-Null }
    } catch {
        Write-Host "Cleanup warning: failed to delete tag $tagSlug"
    }

    try {
        if ($headers) { Invoke-Api DELETE "/api/admin/trends/$trendSlug" $headers | Out-Null }
    } catch {
        Write-Host "Cleanup warning: failed to delete trend $trendSlug"
    }

    try {
        if ($headers) { Invoke-Api DELETE "/api/admin/organizations/$orgSlug" $headers | Out-Null }
    } catch {
        Write-Host "Cleanup warning: failed to delete org $orgSlug"
    }
}