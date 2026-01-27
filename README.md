# Gin Demo Projects

A collection of Go applications demonstrating the use of the Gin web framework for building REST APIs.

## Projects

### Recipes Web API (`recipes-web/`)

A REST API for managing recipes with CRUD operations.

#### Features

- Create, read, update, and delete recipes
- Search recipes by tags
- JSON-based data storage (currently in-memory/file-based)
- Clean architecture with controller, handler, and repository layers

#### API Endpoints

- `GET /recipes` - List all recipes
- `GET /recipes/{id}` - Get recipe by ID
- `POST /recipes` - Create a new recipe
- `PUT /recipes/{id}` - Update an existing recipe
- `DELETE /recipes/{id}` - Delete a recipe
- `GET /recipes/search?tag=X` - Search recipes by tag

#### Running the Application

```bash
cd recipes-web
go run cmd/main.go
```

The server will start on the default Gin port (usually 8080).

### Tasks API (`tasks/`)

A simple task management API with basic CRUD operations and additional features like login simulation and file uploads.

#### Features

- Task CRUD operations
- Task listing with status filtering and pagination
- Simulated login endpoint
- File upload handling

#### API Endpoints

- `GET /tasks` - List tasks with optional status and page filters
- `GET /tasks/{id}` - Get task by ID
- `POST /tasks` - Create a new task
- `POST /login` - Simulated login
- `POST /tasks/{id}/attachment` - Upload attachment for a task

#### Running the Application

```bash
cd tasks
go run main.go
```

### Demo (`demo/`)

Basic Gin examples demonstrating form binding, URI parameters, and JSON handling.

#### Running the Demo

```bash
cd demo
go run main.go
```

## Architecture

The projects follow clean architecture principles:

- **Handlers**: HTTP request/response handling
- **Controllers**: Business logic
- **Repositories**: Data access layer
- **Models**: Data structures

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [rs/xid](https://github.com/rs/xid) - Globally unique identifier generation

## Testing

The recipes-web project includes comprehensive unit tests. Run tests with:

```bash
cd recipes-web
go test ./...
```

See `TESTING_GUIDE.md` and `TEST_SUMMARY.md` for more details on testing.

## Development

This workspace contains multiple Go modules and applications. Each subdirectory may have its own dependencies and can be developed independently.

## License

This project is for educational purposes.</content>
<parameter name="filePath">d:\Golang workspace\rest-world\gin-demo\README.md
