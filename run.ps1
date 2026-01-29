# Prompt for REPO_TYPE
$repoType = Read-Host "Enter REPO_TYPE (e.g., mongo, memory)"
$env:REPO_TYPE = $repoType

# Prompt for SEED_DATA
$seedData = Read-Host "Enter SEED_DATA (true/false)"
$env:SEED_DATA = $seedData

# Run the application
Write-Host "Starting application with REPO_TYPE=$repoType and SEED_DATA=$seedData..." -ForegroundColor Green
go run ./cmd/main.go
