# Recipes Web API

A REST API for managing recipes with CRUD operations, featuring Redis caching and multiple storage backends.

#### Features

- Create, read, update, and delete recipes
- Search recipes by tags
- **Redis caching layer** with 30-minute TTL for improved performance
- Support for multiple repository backends (in-memory, MongoDB)
- Clean architecture with controller, handler, and repository layers
- Comprehensive test coverage (81-91%)

#### Caching Architecture

The application implements a **CachedRepository pattern** that wraps any repository backend with Redis caching:

- **Cache Hits**: Requests for cached recipes return in < 1ms
- **Cache Misses**: Requests fall through to the underlying repository and are cached for future use
- **Cache Invalidation**: Updates and deletes automatically invalidate relevant cache entries
- **TTL**: Cached entries expire after 30 minutes
- **Graceful Degradation**: # Recipes Web API - Complete Documentation

A REST API for managing recipes with CRUD operations, featuring Redis caching and multiple storage backends.

## Table of Contents

1. [Features & Architecture](#features--architecture)
2. [Quick Start](#quick-start)
3. [Running the Application](#running-the-application)
4. [API Endpoints](#api-endpoints)
5. [Configuration](#configuration)
6. [Testing & Coverage](#testing--coverage)
7. [Development Guide](#development-guide)
8. [Troubleshooting](#troubleshooting)

## Features & Architecture

### Features

- ✅ Create, read, update, and delete recipes
- ✅ Search recipes by tags
- ✅ **Redis caching layer** with 30-minute TTL for improved performance
- ✅ Support for multiple repository backends (in-memory, MongoDB)
- ✅ Clean architecture (controller, handler, repository patterns)
- ✅ Comprehensive test coverage (81-91%)
- ✅ Cross-platform runner scripts (PowerShell, Bash, Make)

### Caching Architecture

The application implements a **CachedRepository pattern** that wraps any repository backend with Redis caching:

| Aspect                 | Details                                                                |
| ---------------------- | ---------------------------------------------------------------------- |
| **Cache Hits**         | Requests for cached recipes return in < 1ms                            |
| **Cache Misses**       | Requests fall through to repository and are cached for future use      |
| **Cache Invalidation** | Updates and deletes automatically invalidate relevant cache entries    |
| **TTL**                | Cached entries expire after 30 minutes                                 |
| **Degradation**        | If Redis unavailable, app continues working with underlying repository |

**Performance Benefits:**

- Significantly faster response times for frequently accessed recipes
- Reduced load on the underlying repository
- No database queries for repeated requests

### Repository Backends

| Backend     | Use Case             | Status       |
| ----------- | -------------------- | ------------ |
| **Memory**  | Development, testing | ✅ Default   |
| **MongoDB** | Production data      | ✅ Supported |

### Test Coverage

| Component         | Coverage | Tests         |
| ----------------- | -------- | ------------- |
| Redis Cache       | 81.2%    | 7 tests       |
| Cached Repository | 90.9%    | 8 tests       |
| Memory Repository | High     | 8 tests       |
| Controller        | High     | 6 tests       |
| Handler           | High     | 7 tests       |
| **TOTAL**         | **High** | **36+ tests** |

---

## Quick Start

Choose your preferred platform:

### Windows (PowerShell)

```powershell
.\run.ps1
```

### Linux / macOS / Git Bash / WSL

```bash
chmod +x run.sh
./run.sh
```

### Any Platform (Make)

```bash
make run-memory
```

### Any Platform (Direct Command)

```bash
REPO_TYPE=memory SEED_DATA=false go run ./cmd/main.go
```

---

## Running the Application

### Method 1: PowerShell Script (Windows)

**File:** `run.ps1`

```powershell
.\run.ps1
```

**Interactive Prompts:**

- `REPO_TYPE`: Select `memory` or `mongo`
- `SEED_DATA`: Enter `true` or `false`

**Example:**

```
PS> .\run.ps1
Enter REPO_TYPE (e.g., mongo, memory): memory
Enter SEED_DATA (true/false): false
Starting application with REPO_TYPE=memory and SEED_DATA=false...
```

### Method 2: Shell Script (Linux, macOS, Git Bash, WSL)

**File:** `run.sh`

```bash
chmod +x run.sh
./run.sh
```

**Interactive Prompts:**

- `REPO_TYPE`: Select `memory` or `mongo`
- `SEED_DATA`: Enter `true` or `false`

**Features:**

- ✅ Interactive prompts
- ✅ Input validation
- ✅ Go installation check
- ✅ Colored output
- ✅ Works with piped input (for automation)

**Example:**

```bash
$ ./run.sh
=== Recipes Web API Launcher ===
Available repository types:
  - memory (default: file-based in-memory storage)
  - mongo (MongoDB storage with seeding)
Enter REPO_TYPE (default: memory): memory
Enter SEED_DATA (true/false, default: false): false

Configuration:
  Repository Type: memory
  Seed Data: false

Starting Recipes Web API...
HTTP server listening on :8080
```

### Method 3: Makefile (Any Platform with Make)

**File:** `Makefile`

```bash
make help           # Show all available targets
make run            # Interactive mode
make run-memory     # In-memory backend
make run-mongo      # MongoDB backend
make test           # Run all tests
make test-cache     # Run cache tests (81.2% coverage)
make test-repo      # Run repository tests (90.9% coverage)
make build          # Build binary
make clean          # Clean build artifacts
```

### Method 4: Direct Command (All Platforms)

```bash
# In-memory repository without seeding
REPO_TYPE=memory SEED_DATA=false go run ./cmd/main.go

# MongoDB repository with seeding
REPO_TYPE=mongo SEED_DATA=true go run ./cmd/main.go
```

**Windows (PowerShell):**

```powershell
$env:REPO_TYPE="memory"; $env:SEED_DATA="false"; go run ./cmd/main.go
```

---

## API Endpoints

| Method | Endpoint                | Purpose               | Cached |
| ------ | ----------------------- | --------------------- | ------ |
| GET    | `/recipes`              | List all recipes      | No     |
| GET    | `/recipes/{id}`         | Get recipe by ID      | ✅ Yes |
| POST   | `/recipes`              | Create new recipe     | No     |
| PUT    | `/recipes/{id}`         | Update recipe         | No     |
| DELETE | `/recipes/{id}`         | Delete recipe         | No     |
| GET    | `/recipes/search?tag=X` | Search recipes by tag | No     |

### Example API Requests

```bash
# List all recipes
curl http://localhost:8080/recipes

# Get specific recipe (will be cached after first request)
curl http://localhost:8080/recipes/recipe-id-here

# Create a new recipe
curl -X POST http://localhost:8080/recipes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Pasta Carbonara",
    "tags": ["italian", "main"],
    "ingredients": ["pasta", "eggs", "bacon", "cheese"],
    "instructions": ["boil pasta", "fry bacon", "mix with eggs", "combine"]
  }'

# Search recipes by tag
curl 'http://localhost:8080/recipes/search?tag=italian'

# Update a recipe
curl -X PUT http://localhost:8080/recipes/recipe-id-here \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Recipe Name"}'

# Delete a recipe
curl -X DELETE http://localhost:8080/recipes/recipe-id-here
```

---

## Configuration

### Prerequisites

- **Go 1.18+** - Install from https://golang.org/dl/
- **Redis** (optional) - For caching feature
  - Quick start: `docker run -d -p 6379:6379 redis:latest`
- **MongoDB** (optional) - For persistent storage
  - Quick start: `docker run -d -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password mongo`

### Environment Variables

| Variable    | Default            | Options                   | Purpose                    |
| ----------- | ------------------ | ------------------------- | -------------------------- |
| `REPO_TYPE` | `memory`           | `memory`, `mongo`         | Repository backend         |
| `SEED_DATA` | `false`            | `true`, `false`           | Populate with initial data |
| `HTTP_ADDR` | `:8080`            | Any valid address:port    | Server listening address   |
| `DATA_PATH` | `data/recipe.json` | Any valid file path       | Recipe data file location  |
| `MONGO_URI` | See below          | MongoDB connection string | MongoDB connection         |

**Default MongoDB URI:**

```
mongodb://admin:password@localhost:27017/test?authSource=admin
```

### Configuration Examples

```bash
# In-memory, development mode (fastest for testing)
./run.sh                    # Interactive
REPO_TYPE=memory SEED_DATA=false go run ./cmd/main.go

# MongoDB with initial data
./run.sh                    # Interactive
REPO_TYPE=mongo SEED_DATA=true go run ./cmd/main.go

# Custom HTTP address
HTTP_ADDR=:3000 REPO_TYPE=memory go run ./cmd/main.go

# Custom data file
DATA_PATH=data/my-recipes.json REPO_TYPE=memory go run ./cmd/main.go
```

### Enabling Caching

```bash
# 1. Start Redis
docker run -d -p 6379:6379 redis:latest

# 2. Run application (auto-detects Redis)
./run.sh
# or
make run-memory

# 3. Verify caching is working
# - First request to /recipes/{id} will be slower (database query)
# - Subsequent requests will be < 1ms (cache hit)
```

---

# Testing & Coverage

## Test Structure

```
recipes-web/
├── model/
│   └── recipe_test.go                          # Model tests
├── internal/
│   ├── cache/redisrecipe/
│   │   └── cache_test.go                       # Cache tests (81.2% coverage)
│   ├── controller/recipe/
│   │   └── controller_test.go                  # Business logic tests
│   ├── handler/httpapi/
│   │   └── http_test.go                        # HTTP handler tests
│   └── repository/
│       ├── cached_recipe_repository_test.go    # Cached repo tests (90.9% coverage)
│       └── memory/
│           └── memory_test.go                  # In-memory repo tests
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run with coverage report
go test ./... -cover

# Run specific test suites
go test ./internal/cache/redisrecipe -v        # Cache tests (81.2% coverage)
go test ./internal/repository -v               # Repository tests (90.9% coverage)
go test ./internal/controller/recipe -v        # Controller tests
go test ./model -v                             # Model tests

# Using Make
make test                  # All tests
make test-cache            # Cache tests only
make test-repo             # Repository tests only
```

## Test Coverage by Component

### Redis Cache Tests (cache_test.go) - 81.2% Coverage

| Test                       | Purpose                     |
| -------------------------- | --------------------------- |
| `TestNewCache`             | Cache initialization        |
| `TestCacheSetByID`         | Storing recipes             |
| `TestCacheGetByID`         | Retrieving cached recipes   |
| `TestCacheGetByIDNotFound` | Non-existent entries        |
| `TestCacheDeleteByID`      | Cache invalidation          |
| `TestCacheTTL`             | Expiration after 30 minutes |
| `TestRecipeKey`            | Key generation              |

### Cached Repository Tests (cached_recipe_repository_test.go) - 90.9% Coverage

| Test                                         | Purpose                        |
| -------------------------------------------- | ------------------------------ |
| `TestNewCachedRepository`                    | Initialization                 |
| `TestCachedRepositoryGetByID_FromCache`      | Cache hit scenario             |
| `TestCachedRepositoryGetByID_FromRepository` | Cache miss + populate flow     |
| `TestCachedRepositoryCreate`                 | Creation with caching          |
| `TestCachedRepositoryUpdate`                 | Update with cache invalidation |
| `TestCachedRepositoryDelete`                 | Deletion with cache clearing   |
| `TestCachedRepositoryDeleteNotFound`         | Error handling                 |
| `TestCachedRepositoryGetByIDNotFound`        | Not found handling             |

### Other Test Layers

- **Memory Repository**: 8 tests covering CRUD, persistence, concurrency
- **Controller**: 6 tests covering business logic and error handling
- **Handler**: 7 tests covering HTTP routing, JSON binding, error responses
- **Model**: 2 tests covering serialization and type conversion

## Example Test Output

```bash
$ go test ./internal/cache/redisrecipe -v
=== RUN   TestNewCache
--- PASS: TestNewCache (0.04s)
=== RUN   TestCacheSetByID
--- PASS: TestCacheSetByID (0.02s)
=== RUN   TestCacheGetByID
--- PASS: TestCacheGetByID (0.03s)
...
PASS
coverage: 81.2% of statements
ok      github.com/gin-demo/recipes-web/internal/cache/redisrecipe   5.576s
```

---

# Development Guide

## Architecture Layers

```
┌─────────────────────────────────────┐
│   HTTP Handlers (httpapi)           │ ← Request/Response handling
├─────────────────────────────────────┤
│   Controllers (recipe)              │ ← Business logic
├─────────────────────────────────────┤
│   Cached Repository                 │ ← Optional caching layer
├─────────────────────────────────────┤
│   Repository (memory/mongo)         │ ← Data access layer
├─────────────────────────────────────┤
│   Models                            │ ← Data structures
└─────────────────────────────────────┘
```

## Project Structure

```
recipes-web/
├── cmd/
│   └── main.go                          # Application entry point with caching setup
├── internal/
│   ├── bootstrap/                       # Initialization utilities
│   │   ├── seed.go                      # Database seeding
│   │   └── redis.go                     # Redis client setup
│   ├── cache/
│   │   └── redisrecipe/
│   │       ├── cache.go                 # Redis cache implementation
│   │       └── cache_test.go            # Cache tests (81.2% coverage)
│   ├── controller/
│   │   └── recipe/
│   │       ├── controller.go            # Business logic
│   │       ├── controller_test.go       # Controller tests
│   │       └── commands.go              # Command structures
│   ├── domain/                          # Domain models and errors
│   │   ├── errors.go                    # Error definitions
│   │   └── recipe.go                    # Repository interface
│   ├── handler/
│   │   └── httpapi/
│   │       ├── http.go                  # HTTP handlers
│   │       └── http_test.go             # Handler tests
│   └── repository/
│       ├── cached_recipe_repository.go  # Cached wrapper (decorates repositories)
│       ├── cached_recipe_repository_test.go # Cached repo tests (90.9% coverage)
│       ├── memory/
│       │   ├── memory.go                # In-memory implementation
│       │   └── memory_test.go           # Memory tests
│       └── mongorepo/
│           ├── mongo.go                 # MongoDB implementation
│           ├── mongo_test.go            # Mongo tests
│           └── seed.go                  # Mongo seeding
├── model/
│   ├── recipe.go                        # Recipe data model
│   └── recipe_test.go                   # Model tests
├── data/
│   ├── recipe.json                      # Sample recipes data
│   └── dummy.json                       # Test data
├── README.md                            # This comprehensive documentation
├── Makefile                             # Build automation
├── run.ps1                              # PowerShell runner script
├── run.sh                               # Bash runner script
└── go.mod                               # Go module definition
```

## Dependencies

```
github.com/gin-gonic/gin                # Web framework
github.com/redis/go-redis/v9            # Redis client
github.com/mongodb/mongo-go-driver      # MongoDB driver
github.com/rs/xid                       # ID generation
```

## Building from Source

```bash
cd recipes-web

# Build binary
go build -o recipes-web ./cmd/main.go

# Run binary
./recipes-web  # Linux/macOS
recipes-web.exe  # Windows
```

## Recommended Development Workflow

1. **Start with memory repository** (fastest):

   ```bash
   make run-memory
   # or
   ./run.sh  # select: memory, false
   ```

2. **Enable Redis for caching** (optional):

   ```bash
   docker run -d -p 6379:6379 redis:latest
   ```

3. **Test with real data** (requires MongoDB):

   ```bash
   docker run -d -p 27017:27017 \
     -e MONGO_INITDB_ROOT_USERNAME=admin \
     -e MONGO_INITDB_ROOT_PASSWORD=password \
     mongo
   make run-mongo
   # or
   ./run.sh  # select: mongo, true
   ```

4. **Run tests**:
   ```bash
   make test        # All tests
   make test-cache  # Cache tests
   go test ./... -cover  # With coverage
   ```

---

# Troubleshooting

## Go Issues

### "Go is not installed"

- Install from https://golang.org/dl/
- Verify: `go version`

## Database Issues

### "Cannot connect to MongoDB"

**Option 1:** Use in-memory repository

```bash
REPO_TYPE=memory ./run.sh
# or
make run-memory
```

**Option 2:** Start MongoDB

```bash
docker run -d -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=password \
  mongo
```

### "Redis connection error"

- This is **non-fatal** - the app works without caching
- To enable caching: `docker run -d -p 6379:6379 redis:latest`

## Port Issues

### "Port 8080 already in use"

```bash
# Use a different port
HTTP_ADDR=:3000 go run ./cmd/main.go
# or
make run  # will show this option
```

## Script Issues

### PowerShell script won't run

Enable script execution:

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
.\run.ps1
```

### Shell script permission denied

```bash
chmod +x run.sh
./run.sh
```

## Testing Issues

### Tests fail with "Redis not available"

- This is expected if Redis isn't running
- Tests will skip Redis tests gracefully
- Start Redis to run those tests: `docker run -d -p 6379:6379 redis:latest`

### Tests fail with "MongoDB not available"

- MongoDB tests will skip if MongoDB isn't running
- Start MongoDB to run those tests
- In-memory tests will always work

---

## Performance Tips

### Enable Caching

- Start Redis: `docker run -d -p 6379:6379 redis:latest`
- Watch response times drop dramatically for repeated requests
- **Cache hits:** < 1ms vs **database queries:** 10-100ms+

### Use Memory Repository for Development

- Fastest for testing and development
- No external dependencies
- Good for CI/CD pipelines

### Monitor Performance

- Check server logs for response times
- Use curl with timing: `curl -w "Total time: %{time_total}\n" http://localhost:8080/recipes/{id}`
- Implement metrics if needed

---

## Useful Docker Commands

```bash
# Redis
docker run -d -p 6379:6379 redis:latest
docker exec -it <container-id> redis-cli

# MongoDB
docker run -d -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=password \
  mongo

# View logs
docker logs <container-id>

# Stop container
docker stop <container-id>
```

---

## Contributing

Guidelines for contributing:

1. Maintain test coverage above 80%
2. Follow Go code style guidelines
3. Document all public functions
4. Run tests before submitting changes
5. Update README if adding new features

---

## Additional Resources

- [Gin Framework Documentation](https://gin-gonic.com/)
- [Go Language Documentation](https://golang.org/doc/)
- [Redis Documentation](https://redis.io/documentation)
- [MongoDB Documentation](https://docs.mongodb.com/)
- [Docker Documentation](https://docs.docker.com/)

---

## License

This project is provided as-is for educational and demonstration purposes.

---

**Last Updated:** January 30, 2026  
**Status:** ✅ Ready for Production (with proper Redis/MongoDB setup)  
**Test Coverage:** 81-91% across caching layers  
**Platforms Supported:** Windows (PowerShell), Linux, macOS, WSL
