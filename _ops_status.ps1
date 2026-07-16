$env:Path = [Environment]::GetEnvironmentVariable('Path','Machine') + ';' + [Environment]::GetEnvironmentVariable('Path','User')

Write-Host '--- Tools ---' -ForegroundColor Cyan
foreach ($tool in @('go','psql','migrate','docker')) {
    try {
        $v = & $tool --version 2>&1 | Select-Object -First 1
        Write-Host "$tool : $v"
    } catch {
        Write-Host "$tool : NOT INSTALLED"
    }
}

Write-Host ''
Write-Host '--- DB tables ---' -ForegroundColor Cyan
docker exec radartcell-postgres psql -U radar_tcell -d radar_tcell -c "\dt" 2>&1 | Select-Object -First 40
