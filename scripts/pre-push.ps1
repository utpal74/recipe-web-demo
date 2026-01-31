$branch = git rev-parse --abbrev-ref HEAD

if ($branch -eq "main") {
    Write-Host "Running tests before push to main..." -ForegroundColor Yellow
    Set-Location recipes-web
    go test ./...
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Tests failed! Push aborted." -ForegroundColor Red
        exit 1
    }
    Write-Host "✅ Tests passed!" -ForegroundColor Green
}
exit 0
