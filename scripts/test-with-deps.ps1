Param()
Set-StrictMode -Version Latest

$Root = Split-Path -Path $PSScriptRoot -Parent
Set-Location $Root

Write-Host "Starting test dependencies via docker compose..."
docker compose up -d

Write-Host "Waiting for services to be healthy..."
$tries = 0
while ($tries -lt 30) {
    try {
        docker compose exec -T mongo mongosh --eval "db.adminCommand('ping')" | Out-Null
        docker compose exec -T redis redis-cli ping | Out-Null
        break
    } catch {
        Start-Sleep -Seconds 1
        $tries++
    }
}

if ($tries -ge 30) {
    Write-Error "Timed out waiting for services"
    docker compose logs --no-color
    docker compose down --volumes --remove-orphans
    exit 1
}

Write-Host "Services are ready — running tests..."
& go test ./...
$exitCode = $LASTEXITCODE

Write-Host "Tearing down test dependencies..."
docker compose down --volumes --remove-orphans

exit $exitCode
