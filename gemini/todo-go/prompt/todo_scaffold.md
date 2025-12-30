

# Generate a detailed prompt for Gemini to scaffold a Go web TODO app

Use this as a detailed prompt for Gemini (CLI or UI) to scaffold a Go web TODO app.

## Role and goal

You are a senior Go backend engineer.
Generate a complete, production-ready scaffold for a **Go web TODO application** with a clean, testable architecture and clear separation of concerns.[^1][^2]

The output must be a project that can be cloned, run, and extended by another engineer.

***

## Functional requirements

- Implement a **RESTful JSON API** for managing TODO items.
- TODO fields:
    - `id` (server-generated).
    - `title` (required, non-empty).
    - `description` (optional).
    - `completed` (boolean).
    - `created_at`, `updated_at` (server-managed timestamps).
- API endpoints:
    - `POST /api/todos` – create todo.
    - `GET /api/todos` – list todos, with optional `completed` filter.
    - `GET /api/todos/{id}` – get a single todo.
    - `PUT /api/todos/{id}` – update title, description, completed.
    - `DELETE /api/todos/{id}` – delete todo.
- API must be usable by a browser-based frontend (CORS, JSON, standard HTTP verbs).[^3]

***

## Non-functional requirements

- Use **Go modules** and idiomatic Go 1.21+ features where helpful.[^4]
- Favor **simplicity and readability** over premature optimization.
- No unnecessary third-party dependencies; prefer standard library unless a library clearly improves clarity.

***

## Architecture and project layout

Use a simple clean/hexagonal style structure, inspired by common Go clean architecture examples.[^2][^1]

Target layout:

```text
go-todo/
  cmd/
    server/
      main.go
  internal/
    todo/          # domain model + use cases
      entity.go
      service.go
    http/
      handler.go   # HTTP handlers & routing
      middleware.go
    storage/
      memory/
        repo.go    # in-memory impl for tests/dev
      sqlite/
        repo.go    # production repository
    config/
      config.go
  pkg/
    logger/
      logger.go
  api/
    openapi.yaml   # optional: OpenAPI spec (nice-to-have)
  migrations/
    001_create_todos.sql
  go.mod
  go.sum
  Makefile (or simple make-like commands in README)
  README.md
```

Constraints:

- **Domain layer (`internal/todo`)**:
    - Define the `Todo` entity and interfaces for repository and service.
    - Implement business rules such as validation, default values, and state transitions.
- **Storage layer (`internal/storage`)**:
    - Implement repository interfaces for in-memory and SQLite backends.
- **HTTP layer (`internal/http`)**:
    - Implement handlers that depend only on service interfaces, not on concrete storage.
    - Use `net/http` plus a light router (e.g., `chi`) or just `http.ServeMux`.[^3]
- **Config and wiring**:
    - Wire dependencies in `cmd/server/main.go` (config → storage → services → HTTP).

***

## Persistence (SQLite)

- Use `database/sql` and a well-supported SQLite driver (`github.com/mattn/go-sqlite3` or `modernc.org/sqlite`).[^5][^3]
- Provide:
    - `migrations/001_create_todos.sql` that creates the `todos` table.
    - Code that runs migrations at startup if needed.
- Table design:
    - Primary key integer `id` with AUTOINCREMENT.
    - Index on `completed` for filtered queries.
- Use:
    - Connection pooling via `database/sql`.
    - Context-aware queries (`ctx context.Context`).
    - Prepared statements or parameterized queries to avoid injection.
- Provide a **configurable DSN** via environment variable, defaulting to `./data/todos.db`.

***

## HTTP API design

- Implement routes under `/api/` prefix.
- Behavior:
    - All responses are JSON with `Content-Type: application/json`.
    - On validation errors, return:
        - Status `400`.
        - Body: `{"error": "validation_error", "details": {"field": "message"}}`.
    - On not found:
        - Status `404`.
        - Body: `{"error": "not_found", "message": "todo not found"}`.
    - On internal errors:
        - Status `500`.
        - Body: `{"error": "internal_error", "message": "..."}`.
- Add minimal middleware:
    - Request logging (method, path, status, duration).
    - Panic recovery → 500 with structured error.
    - CORS allowing local dev (e.g., `http://localhost:3000`) with configurable origins.

***

## Configuration and environment

- Small config package reading from env with sane defaults:
    - `HTTP_ADDR` default `:8080`.
    - `SQLITE_DSN` default `./data/todos.db`.
    - `LOG_LEVEL` default `info`.
- Provide a `Config` struct and `Load()` function that:
    - Reads env vars.
    - Validates them.
    - Returns a ready-to-use config instance.

***

## Testing requirements

Write **tests first or alongside code** so the scaffold includes realistic examples.[^6][^7][^8]

- **Service layer tests (`internal/todo/service_test.go`)**:
    - Use a **mock or in-memory repository** to test business logic in isolation.[^9][^10]
    - Table-driven tests for:
        - Creating valid/invalid todos.
        - Updating fields and completed status.
        - Handling repository errors.
- **HTTP tests (`internal/http/handler_test.go`)**:
    - Use `net/http/httptest` to test handlers.
    - Cover:
        - Happy paths for each endpoint.
        - Invalid JSON / missing required fields → 400.
        - Not found → 404.
        - Repository error → 500.
- Aim for **good coverage of core flows**, but keep tests clear and focused.
- Provide commands:
    - `go test ./...` in README and/or Makefile.[^1]

***

## Best practices to follow

- Use **dependency injection**:
    - Pass interfaces into constructors (services, handlers) rather than using globals.
- Keep functions small and focused; handle errors explicitly (`if err != nil { ... }`).[^4]
- Avoid global mutable state (no global DB handles; only store them in structs that are injected).
- Prefer **flat packages** where possible, but keep domain, HTTP, and storage concerns separated.[^11][^4]
- Add brief doc comments on exported types and functions.

***

## Developer ergonomics

- Provide a simple **Makefile** or at least a section in README with commands:
    - `make run` → `go run ./cmd/server`.
    - `make test` → `go test ./...`.
    - `make lint` (optional) with `go vet` or similar.
- Include example **curl** commands in README to:
    - Create a todo.
    - List todos.
    - Filter by `completed`.
    - Update and delete.

***

## Output format

Return the scaffold in this order:

1. Short high-level overview of the architecture and reasoning.
2. The `go.mod` file.
3. Directory tree of the scaffold.
4. Full source code for:
    - `cmd/server/main.go`.
    - `internal/todo/entity.go`, `service.go`, and tests.
    - `internal/storage/memory/repo.go` and tests (if needed).
    - `internal/storage/sqlite/repo.go`.
    - `internal/http/handler.go`, `middleware.go`, and tests.
    - `internal/config/config.go`.
    - `pkg/logger/logger.go`.
5. `migrations/001_create_todos.sql`.
6. `Makefile` (or equivalent) and `README.md` with run/test instructions and sample curl commands.

Ensure the final project compiles and is runnable with:

```bash
go mod tidy
go run ./cmd/server
```
