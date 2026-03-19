# User API

A small Go-based REST API for managing users, departments, attendance, payroll, and reports.

## ✅ Features
- CRUD endpoints for **workers**, **departments**, **attendance**, and **payroll**
- Validation using `go-playground/validator`
- PostgreSQL persistence via GORM
- Swagger/OpenAPI documentation via `swag` (Swagger annotations in source)

## 🚀 Prerequisites
- Go 1.20+ installed
- PostgreSQL database
- A `.env` file in the repo root with at least:
  - `DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable`
  - `JWT_SECRET=your_secret_here`

## ⚙️ Setup
1. Install dependencies:
   ```sh
   go mod download
   ```
2. Create or update `.env` with your database connection string.
3. Run migrations (auto-run by the server on startup):
   ```sh
   go run ./cmd/main.go
   ```

## ▶️ Run the server
```sh
go run ./cmd/main.go
```
The server will run on `:8080` by default.

## 🧪 Tests
### Run all tests
```sh
go test ./...
```

### Integration tests
Integration tests use the same database configured in `.env`. Make sure it points to a safe test database.

### Test environment
Use `.env.test.example` to set environment variables during CI or local test runs:
- `JWT_SECRET=super-secret-key`
- `DATABASE_URL=postgres://postgres:postgres@localhost:5432/databaseTest?sslmode=disable`

Run with `go test` via a script or by setting `DATABASE_URL`/`JWT_SECRET` in your environment.

## 🧱 Project structure details

The repo is intentionally cleanly separated:
- `cmd/` : service bootstrap, route wiring, DI for handlers/services
- `internal/database/` : connection helper (`database.Connect`) and shared `database.DB`
- `internal/*/model.go` : GORM model definitions
- `internal/*/repository.go` : data access (DB queries)
- `internal/*/service.go` : business logic layer
- `internal/*/handler.go` : Gin HTTP handlers + validation
- `internal/*/*_integration_test.go` : endpoint integration tests (in DB)
- `internal/*/*_service_test.go` : pure unit tests using fake repos (where interface exists)

This is effective for separation of concerns and for swapping persistence (e.g. in-memory/mocks for tests) while keeping endpoint contract stable.

## 🛠️ Using golang-migrate

In this repository we currently do migrations with `gorm.DB.AutoMigrate` at startup. To switch to `golang-migrate` for explicit migration control, do:

1. Install CLI:
```sh
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

2. Add migration files:
```shn
migrations/
  000001_create_workers_table.up.sql
  000001_create_workers_table.down.sql
  000002_create_departments_table.up.sql
  000002_create_departments_table.down.sql
  000003_create_attendances_table.up.sql
  000003_create_attendances_table.down.sql
  000004_create_payrolls_table.up.sql
  000004_create_payrolls_table.down.sql
```

3. Run migrations:
```sh
migrate -path migrations -database "$DATABASE_URL" up
```

4. Rollback one step:
```sh
migrate -path migrations -database "$DATABASE_URL" down 1
```

5. In `cmd/main.go`, replace `AutoMigrate(...)` with a call to migrate first (safe startup):
```go
import (
  "github.com/golang-migrate/migrate/v4"
  "github.com/golang-migrate/migrate/v4/database/postgres"
  "github.com/golang-migrate/migrate/v4/source/file"
)

// in main
m, err := migrate.New("file://migrations", os.Getenv("DATABASE_URL"))
if err != nil { panic(err) }
if err := m.Up(); err != nil && err != migrate.ErrNoChange { panic(err) }
```


## 📄 API Documentation (Swagger)
Generate/update Swagger docs:
```sh
swag init -g cmd/main.go --output cmd/docs
```

Serve Swagger UI by running the server and visiting:
```
http://localhost:8080/swagger/index.html
```

## 📁 Project structure

- `cmd/main.go` - application entry point and router setup
- `internal/` - application modules
  - `workers/` - worker CRUD + handler/tests
  - `departments/` - department CRUD + handler/tests
  - `attendance/` - attendance endpoints + handler/tests
  - `payroll/` - payroll calculation + handler/tests
  - `auth/` - JWT authentication
  - `database/` - DB connection helper

---
