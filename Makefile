.PHONY: help run run-memory run-mongo test test-cache test-repo build clean

# Default target
help:
	@echo "Recipes Web API - Make targets:"
	@echo ""
	@echo "  make run                 - Run the server (interactive mode)"
	@echo "  make run-memory          - Run with in-memory repository"
	@echo "  make run-mongo           - Run with MongoDB repository"
	@echo "  make test                - Run all tests"
	@echo "  make test-cache          - Run Redis cache tests"
	@echo "  make test-repo           - Run repository tests"
	@echo "  make build               - Build the binary"
	@echo "  make clean               - Clean build artifacts"
	@echo ""
	@echo "Environment variables:"
	@echo "  REPO_TYPE=memory|mongo   - Repository type (default: memory)"
	@echo "  SEED_DATA=true|false     - Seed database with initial data (default: false)"
	@echo "  HTTP_ADDR=:8080          - HTTP server address (default: :8080)"

# Run with interactive mode
run:
	@echo "Starting Recipes Web API..."
	@read -p "Enter REPO_TYPE (memory/mongo) [memory]: " REPO_TYPE; \
	REPO_TYPE=$${REPO_TYPE:-memory}; \
	read -p "Enter SEED_DATA (true/false) [false]: " SEED_DATA; \
	SEED_DATA=$${SEED_DATA:-false}; \
	REPO_TYPE=$$REPO_TYPE SEED_DATA=$$SEED_DATA go run ./cmd/main.go

# Run with memory repository
run-memory:
	@echo "Starting with memory repository..."
	REPO_TYPE=memory SEED_DATA=false go run ./cmd/main.go

# Run with MongoDB
run-mongo:
	@echo "Starting with MongoDB repository..."
	REPO_TYPE=mongo SEED_DATA=true go run ./cmd/main.go

# Run all tests
test:
	go test ./...

# Run cache tests only
test-cache:
	go test ./internal/cache/redisrecipe -v

# Run repository tests
test-repo:
	go test ./internal/repository -v

# Build the binary
build:
	go build -o recipes-web ./cmd/main.go

# Clean build artifacts
clean:
	rm -f recipes-web
	go clean
