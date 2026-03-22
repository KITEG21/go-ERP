# User API

A production-ready Go REST API for managing employees, departments, attendance tracking, payroll processing, and attendance reports.

**Base URL:** `http://localhost:8080`

---

## Features

| Feature | Description |
|---------|-------------|
| **CRUD Operations** | Workers, Departments, Attendance, Payroll |
| **Authentication** | JWT-based (register/login) |
| **Validation** | Centralized struct validation with custom error messages |
| **Pagination** | All list endpoints support `page` and `page_size` params |
| **Structured Logging** | zerolog with dev/prod formatting |
| **Rate Limiting** | 5 requests/sec with burst of 10 |
| **CORS** | Configured for local development |
| **API Documentation** | Swagger + Scalar UI |
| **Graceful Shutdown** | 10-second timeout on SIGINT/SIGTERM |
| **Database** | PostgreSQL via GORM with connection pooling |
| **Monitoring** | Prometheus metrics + Grafana dashboards |

---

## Prerequisites

- Go 1.21+
- PostgreSQL
- `.env` file (see `.env.example`)

---

## Quick Start

```sh
# 1. Install dependencies
go mod download

# 2. Configure environment
cp .env.example .env
# Edit .env with your DATABASE_URL and JWT_SECRET

# 3. Run the server
go run ./cmd/main.go
```

---

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes | - | PostgreSQL connection string |
| `JWT_SECRET` | No | `default_secret` | JWT signing key (warns if missing) |
| `APP_ENV` | No | `development` | `development` or `production` (affects log format) |

### Example `.env`

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/user_api?sslmode=disable
JWT_SECRET=your-super-secret-key
APP_ENV=development
```

---

## API Documentation

Generate Swagger docs:
```sh
swag init -g cmd/main.go --output cmd/docs
```

Access documentation:
- **Swagger UI:** `http://localhost:8080/swagger/index.html`
- **Swagger JSON:** `http://localhost:8080/swagger/doc.json`
- **Scalar UI:** `http://localhost:8080/scalar`

---

## API Endpoints

All endpoints are prefixed with `/api/v1`. Authentication required unless marked **(public)**.

### Authentication (public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login, returns JWT |

### Health Endpoints (public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/healthz` | Liveness probe - returns OK if service is running |
| GET | `/readyz` | Readiness probe - checks DB connectivity |
| GET | `/health/liveness` | Same as `/healthz` |
| GET | `/health/readiness` | Same as `/readyz` |

**Liveness Response (200):**
```json
{"status": "ok"}
```

**Readiness Response (200):**
```json
{"status": "ready"}
```

**Readiness Response (503 - DB unavailable):**
```json
{"status": "unavailable", "error": "connection refused"}
```

### Workers

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/workers` | List workers (paginated) |
| POST | `/api/v1/workers` | Create worker |
| GET | `/api/v1/workers/:id` | Get worker by ID |
| PUT | `/api/v1/workers` | Update worker |
| DELETE | `/api/v1/workers/:id` | Delete worker |

### Departments

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/departments` | List departments (paginated) |
| POST | `/api/v1/departments` | Create department |
| GET | `/api/v1/departments/:id` | Get department by ID |
| PUT | `/api/v1/departments` | Update department |
| DELETE | `/api/v1/departments/:id` | Delete department |

### Attendance

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/attendances` | List attendances (paginated) |
| POST | `/api/v1/attendances` | Create attendance record |
| GET | `/api/v1/attendances/:id` | Get attendance by ID |
| PUT | `/api/v1/attendances/:id` | Update attendance (check-out) |
| GET | `/api/v1/attendances/worker/:worker_id` | Get by worker |

### Payroll

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/payrolls` | List payrolls (paginated) |
| GET | `/api/v1/payrolls/:workerId` | Get payroll by worker |
| POST | `/api/v1/payroll/calculate` | Calculate net salary |

### Reports

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/reports/workers/attendance` | Worker attendance report |

Query params: `start_date`, `end_date` (required), `department_id`, `worker_id` (optional)

---

## Authentication

Include JWT token in Authorization header:

```
Authorization: Bearer <your_token>
```

### Example

```sh
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret"}'

# Use token
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/workers
```

---

## Pagination

All list endpoints support pagination:

```
GET /api/v1/workers?page=1&page_size=10
```

Response:
```json
{
  "page": 1,
  "page_size": 10,
  "total": 50,
  "total_pages": 5,
  "data": [...]
}
```

---

## Error Responses

### Validation Error (400)
```json
[
  {
    "field": "email",
    "tag": "email",
    "value": "",
    "message": "Must be a valid email address"
  }
]
```

### Standard Error
```json
{
  "code": "404",
  "message": "Resource not found"
}
```

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 404 | Not Found |
| 429 | Too Many Requests |
| 500 | Internal Server Error |

---

## Project Structure

```
user_api/
├── cmd/
│   ├── main.go              # Entry point, graceful shutdown
│   └── docs/                # Generated Swagger docs
├── internal/
│   ├── app/
│   │   ├── application.go   # DI setup, CORS, rate limiting
│   │   └── router.go       # Route definitions
│   ├── attendance/
│   │   ├── attendance_handler.go
│   │   ├── attendance_model.go
│   │   ├── attendance_repository.go
│   │   └── attendance_service.go
│   ├── auth/
│   │   ├── auth_handler.go
│   │   ├── auth_service.go
│   │   ├── jwt_service.go
│   │   └── user_model.go
│   ├── common/
│   │   ├── api_error.go
│   │   ├── validator.go          # Centralized validation
│   │   └── validation_error_response.go
│   ├── database/
│   │   ├── db.go           # Connection + pooling
│   │   └── migrations.go
│   ├── departments/
│   │   ├── department_handler.go
│   │   ├── department_model.go
│   │   ├── department_repository.go
│   │   ├── department_repository_interface.go
│   │   └── department_service.go
│   ├── dto/
│   │   ├── attendance/
│   │   ├── auth/
│   │   ├── department/
│   │   ├── pagination/    # Pagination DTOs
│   │   ├── payroll/
│   │   ├── report/
│   │   └── worker/
│   ├── logger/
│   │   └── logger.go      # zerolog wrapper
│   ├── metrics/
│   │   └── metrics.go     # Prometheus metrics definitions
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   ├── logger_middleware.go
│   │   └── metric_middleware.go  # Prometheus middleware
│   ├── payroll/
│   │   ├── payroll_handler.go
│   │   ├── payroll_model.go
│   │   ├── payroll_repository.go
│   │   ├── payroll_repository_interface.go
│   │   └── payroll_service.go
│   ├── reports/
│   │   ├── report_handler.go
│   │   └── report_service.go
│   └── workers/
│       ├── worker_handler.go
│       ├── worker_model.go
│       ├── worker_repository.go
│       ├── worker_repository_interface.go
│       └── worker_service.go
├── grafana/
│   ├── dashboards/
│   │   └── http-metrics.json     # Pre-built Grafana dashboard
│   └── provisioning/
│       ├── dashboards/
│       │   └── dashboards.yml    # Dashboard provisioning
│       └── datasources/
│           └── datasources.yml   # Prometheus datasource
├── .dockerignore
├── .env.example
├── .env.test.example
├── docker-compose.yml       # Docker Compose stack
├── Dockerfile              # Multi-stage build
├── prometheus.yml          # Prometheus config
├── go.mod
└── readme.md
```

---

## Testing

```sh
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run integration tests (uses .env.test database)
go test -tags=integration ./...
```

---

## Configuration Notes

### CORS Origins
Edit `internal/app/application.go` to add your frontend domain:
```go
AllowOrigins: []string{"http://localhost:3000", "https://your-frontend.com"},
```

### Rate Limiting
Default: 5 requests/sec, burst 10. Edit in `application.go`:
```go
limiter := tollbooth.NewLimiter(5, nil)
limiter.SetBurst(10)
```

### Database Connection Pool
Configure in `internal/database/db.go`:
```go
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| `DATABASE_URL not set` | Add to `.env` file |
| `JWT_SECRET not set` warning | Set in `.env` (recommended for production) |
| CORS errors | Add your origin to `AllowOrigins` in `application.go` |
| Port 8080 in use | Change `Addr` in `main.go` or use reverse proxy |

---

## Generating Documentation

After modifying Swagger annotations in handlers:

```sh
swag init -g cmd/main.go --output cmd/docs
```

---

## Contributing

1. Follow existing project structure
2. Add unit tests for services, integration tests for handlers
3. Update Swagger annotations for new endpoints
4. Run `swag init` after adding/modifying endpoints
5. Ensure `go fmt` and `go vet` pass before committing

---

## Monitoring (Prometheus + Grafana)

A complete monitoring stack is provided via Docker Compose.

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      DOCKER COMPOSE                          │
│                                                             │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐   │
│  │   Your API   │    │  Prometheus │    │   Grafana   │   │
│  │   :8080      │───▶│   :9090     │───▶│   :3000     │   │
│  │   /metrics   │scrape│             │datasource│             │   │
│  └─────────────┘    └─────────────┘    └─────────────┘   │
│                              │                               │
│                              ▼                               │
│                       ┌─────────────┐                       │
│                       │  prometheus │                       │
│                       │  data/      │  (15d retention)     │
│                       └─────────────┘                       │
└─────────────────────────────────────────────────────────────┘
```

### Quick Start

```sh
# 1. Build and start all services
docker-compose up -d

# 2. Check services are running
docker-compose ps

# 3. Access services
# - API:        http://localhost:8080
# - Prometheus: http://localhost:9090
# - Grafana:    http://localhost:3000 (admin/admin)

# 4. Verify metrics are being scraped
curl http://localhost:8080/metrics | grep user_api

# 5. Generate test traffic
for i in {1..100}; do curl -s http://localhost:8080/api/v1/workers > /dev/null; done
```

### Services

| Service | Port | Description |
|---------|------|-------------|
| API | 8080 | Your Go application |
| Prometheus | 9090 | Metrics storage & querying |
| Grafana | 3000 | Dashboards (admin/admin) |

### Metrics Exposed

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `http_requests_total` | Counter | method, path, status | Total HTTP requests |
| `user_api_http_request_duration_seconds` | Histogram | method, path | Request latency |
| `db_connections_open` | Gauge | - | Open DB connections |

### Grafana Dashboard

Pre-configured dashboard: **User API - HTTP Metrics**

Panels:
- **Overview**: Requests/sec, Error Rate %, p95 Latency, Total Requests (24h)
- **Latency**: Percentile graphs (p50/p95/p99), Latency by endpoint
- **Status Codes**: Pie chart breakdown, Status code trends
- **Endpoints**: Stacked bar chart by endpoint

### Prometheus UI

Access at `http://localhost:9090` to:
- Query metrics using PromQL
- View scraped targets
- Check configuration

Example PromQL queries:
```promql
# Requests per second
sum(rate(http_requests_total[5m]))

# Error rate
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m])) * 100

# p95 latency
histogram_quantile(0.95, sum(rate(user_api_http_request_duration_seconds_bucket[5m])) by (le))
```

### Stopping

```sh
# Stop all services
docker-compose down

# Stop and remove volumes (clears all data)
docker-compose down -v
```

### Files

| File | Purpose |
|------|---------|
| `docker-compose.yml` | Container orchestration |
| `prometheus.yml` | Prometheus scrape configuration |
| `grafana/provisioning/datasources/datasources.yml` | Auto-configure Prometheus datasource |
| `grafana/provisioning/dashboards/dashboards.yml` | Dashboard provisioning config |
| `grafana/dashboards/http-metrics.json` | Pre-built dashboard |
| `Dockerfile` | Multi-stage Go build |
| `.dockerignore` | Exclude unnecessary files |
