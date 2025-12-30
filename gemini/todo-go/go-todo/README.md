# Go TODO App

This is a production-ready scaffold for a Go web TODO application.

## Overview

The application is a RESTful JSON API for managing TODO items. It follows a clean architecture pattern with a clear separation of concerns between the domain, application, and infrastructure layers.

- **`cmd/server`**: The main application entry point.
- **`internal/todo`**: The core domain logic for TODOs.
- **`internal/http`**: The HTTP handlers, routing, and middleware.
- **`internal/storage`**: The storage implementations (in-memory and SQLite).
- **`internal/config`**: Configuration loading.
- **`pkg/logger`**: A simple structured logger.
- **`migrations`**: Database migrations.

## Requirements

- Go 1.21+
- SQLite

## Getting Started

### Running the application

To run the application, use the following command:

```bash
make run
```

This will start the server on `:8080` by default.

You can also run it with `go run`:

```bash
go mod tidy
go run ./cmd/server
```

### Running tests

To run the tests, use the following command:

```bash
make test
```

### Configuration

The application can be configured using environment variables:

- `HTTP_ADDR`: The address for the HTTP server to listen on. Default: `:8080`.
- `SQLITE_DSN`: The Data Source Name for the SQLite database. Default: `./data/todos.db`.
- `LOG_LEVEL`: The log level (`debug`, `info`, `warn`, `error`). Default: `info`.
- `CORS_ALLOWED_ORIGINS`: Comma-separated list of allowed CORS origins. Default: `http://localhost:3000`.

## API Usage

Here are some example `curl` commands to interact with the API:

### Create a new TODO

```bash
curl -X POST http://localhost:8080/api/todos \
-H "Content-Type: application/json" \
-d '{"title": "My first TODO", "description": "This is a description."}'
```

### List all TODOs

```bash
curl http://localhost:8080/api/todos
```

### List completed TODOs

```bash
curl http://localhost:8080/api/todos?completed=true
```

### Get a single TODO

```bash
# Replace {id} with the ID of the TODO
curl http://localhost:8080/api/todos/{id}
```

### Update a TODO

```bash
# Replace {id} with the ID of the TODO
curl -X PUT http://localhost:8080/api/todos/{id} \
-H "Content-Type: application/json" \
-d '{"title": "Updated Title", "description": "Updated description", "completed": true}'
```

### Delete a TODO

```bash
# Replace {id} with the ID of the TODO
curl -X DELETE http://localhost:8080/api/todos/{id}
```
