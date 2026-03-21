# User API Documentation

A production-ready REST API for managing employees, departments, attendance tracking, payroll processing, and attendance reports.

**Base URL:** `http://localhost:8080/api/v1`

**Swagger UI:** `http://localhost:8080/swagger/index.html`
**Scalar UI:** `http://localhost:8080/scalar`

---

## Authentication

Most endpoints require JWT authentication. Include the token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register a new user |
| POST | `/auth/login` | Authenticate and get JWT token |

---

## Common Features

### Pagination

All list endpoints support pagination with query parameters:

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number |
| `page_size` | int | 10 | Items per page |

**Response Format:**
```json
{
  "page": 1,
  "page_size": 10,
  "total": 50,
  "total_pages": 5,
  "data": [...]
}
```

### Error Responses

**Validation Error (400):**
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

**Standard Error:**
```json
{
  "code": "404",
  "message": "Resource not found"
}
```

---

## Auth

### Register User

Register a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secret123"
}
```

**Validation:** `name`, `email` (valid format), `password` (min 6 chars) are required.

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### Login

Authenticate and receive JWT token.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "secret123"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## Workers

### Worker Object

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "department_id": 1,
  "department": {
    "id": 1,
    "name": "Engineering",
    "description": "Software development team"
  },
  "salary": 50000.00,
  "hire_date": "2024-01-15"
}
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/workers` | List all workers (paginated) |
| POST | `/workers` | Create a new worker |
| GET | `/workers/:id` | Get worker by ID |
| PUT | `/workers` | Update a worker |
| DELETE | `/workers/:id` | Delete a worker |

### Create Worker

**Request Body:**
```json
{
  "name": "Jane Smith",
  "email": "jane@example.com",
  "department_id": 2,
  "hire_date": "2024-02-01"
}
```

**Validation:**
- `name`: required, 3-100 characters
- `email`: required, valid email format
- `department_id`: optional, must be > 0
- `hire_date`: optional, format YYYY-MM-DD (defaults to today)

**Response (201):** Returns created worker object.

---

### Update Worker

**Request Body:**
```json
{
  "id": 1,
  "name": "John Doe Updated",
  "email": "john.updated@example.com",
  "department_id": 3
}
```

**Validation:** `id`, `name`, `email` are required.

**Response (200):** Returns updated worker object.

---

## Departments

### Department Object

```json
{
  "id": 1,
  "name": "Engineering",
  "description": "Software development team"
}
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/departments` | List all departments (paginated) |
| POST | `/departments` | Create a new department |
| GET | `/departments/:id` | Get department by ID |
| PUT | `/departments` | Update a department |
| DELETE | `/departments/:id` | Delete a department |

### Create Department

**Request Body:**
```json
{
  "name": "Marketing",
  "description": "Marketing and sales team"
}
```

**Validation:** `name` is required.

**Response (201):** Returns created department object.

---

### Update Department

**Request Body:**
```json
{
  "id": 1,
  "name": "Engineering & Tech",
  "description": "Updated description"
}
```

**Validation:** `id` and `name` are required.

**Response (200):** Returns updated department object.

---

## Attendance

### Attendance Object

```json
{
  "id": 1,
  "worker_id": 1,
  "worker": { /* Worker object */ },
  "date": "2024-01-15",
  "check_in": "09:00:00",
  "check_out": "17:30:00",
  "status": 1
}
```

### Status Values

| Value | Status |
|-------|--------|
| 0 | CheckedIn |
| 1 | CheckedOut |
| 2 | Expired |

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/attendances` | List all attendances (paginated) |
| POST | `/attendances` | Create attendance record |
| GET | `/attendances/:id` | Get attendance by ID |
| PUT | `/attendances/:id` | Update attendance (check-out) |
| GET | `/attendances/worker/:worker_id` | Get attendance for a worker |

### Create Attendance

**Request Body:**
```json
{
  "worker_id": 1,
  "date": "2024-01-15",
  "check_in": "09:00:00",
  "check_out": "17:00:00"
}
```

**Validation:** `worker_id`, `date`, `check_in`, `check_out` are all required.

**Response (201):** Returns created attendance object.

---

### Update Attendance (Check-Out)

**Request Body:**
```json
{
  "id": 1,
  "check_out": "18:00:00"
}
```

**Validation:** `id` and `check_out` (format: HH:MM:SS) are required.

**Response (200):** Returns updated attendance object.

---

## Payroll

### Payroll Object

```json
{
  "id": 1,
  "worker_id": 1,
  "worker": { /* Worker object */ },
  "month": "2024-01",
  "base_salary": 5000.00,
  "bonus": 500.00,
  "deductions": 200.00,
  "net_salary": 5300.00,
  "status": 1
}
```

### Status Values

| Value | Status |
|-------|--------|
| 0 | Pending |
| 1 | Processed |
| 2 | Failed |

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/payrolls` | List all payroll records (paginated) |
| GET | `/payrolls/:workerId` | Get payroll for a worker |
| POST | `/payroll/calculate` | Calculate payroll (base + bonus - deductions) |

### Calculate Payroll

Calculates net salary from base salary, bonus, and deductions.

**Endpoint:** `POST /payroll/calculate`

**Request Body:**
```json
{
  "worker_id": 1,
  "month": "2024-01",
  "base_salary": 5000.00,
  "bonus": 500.00,
  "deductions": 200.00
}
```

**Validation:** `worker_id` and `month` are required.

**Response (200):**
```json
{
  "worker_id": 1,
  "month": "2024-01",
  "base_salary": 5000.00,
  "bonus": 500.00,
  "deductions": 200.00,
  "net_salary": 5300.00
}
```

**Formula:** `net_salary = base_salary + bonus - deductions`

---

### Get Payroll by Worker

**Endpoint:** `GET /payrolls/:workerId`

**Response (200):** Returns array of payroll records for the worker.

---

## Reports

### Worker Attendance Report Object

```json
{
  "workerId": 1,
  "workerName": "John Doe",
  "department": "Engineering",
  "daysPresent": 22,
  "daysAbsent": 3,
  "hoursWorked": 176.5,
  "attendanceRate": 88.0
}
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/reports/workers/attendance` | Get attendance report for workers |

### Get Worker Attendance Report

Aggregated attendance data for workers within a date range.

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `start_date` | string | Yes | Start date (YYYY-MM-DD) |
| `end_date` | string | Yes | End date (YYYY-MM-DD) |
| `department_id` | integer | No | Filter by department |
| `worker_id` | integer | No | Filter by specific worker |

**Example Request:**
```
GET /reports/workers/attendance?start_date=2024-01-01&end_date=2024-01-31&department_id=1
```

**Response (200):**
```json
[
  {
    "workerId": 1,
    "workerName": "John Doe",
    "department": "Engineering",
    "daysPresent": 22,
    "daysAbsent": 3,
    "hoursWorked": 176.5,
    "attendanceRate": 88.0
  }
]
```

---

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (missing/invalid token) |
| 404 | Not Found |
| 429 | Too Many Requests (rate limited) |
| 500 | Internal Server Error |

---

## Rate Limiting

The API limits requests to **5 requests per second** with a burst of 10.

If exceeded, returns `429 Too Many Requests`.

---

## Project Structure

```
user_api/
├── cmd/
│   ├── main.go              # Entry point + graceful shutdown
│   └── docs/                # Generated Swagger docs
├── internal/
│   ├── app/
│   │   ├── application.go   # DI setup, CORS, rate limiting
│   │   └── router.go        # Route definitions
│   ├── attendance/           # Attendance module
│   ├── auth/                # JWT authentication
│   ├── common/              # Shared utilities (errors, validation)
│   ├── database/            # DB connection + migrations
│   ├── departments/         # Department module
│   ├── dto/                 # Data Transfer Objects
│   │   ├── pagination/      # Pagination DTOs
│   │   └── ...
│   ├── logger/              # zerolog wrapper
│   ├── middleware/          # Auth & logging middleware
│   ├── payroll/            # Payroll module
│   ├── reports/            # Reports module
│   └── workers/            # Worker module
├── .env.example
└── readme.md
```

---

## Setup

1. Install dependencies:
   ```sh
   go mod download
   ```

2. Create `.env` file:
   ```env
   DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable
   JWT_SECRET=your_secret_here
   APP_ENV=development
   ```

3. Generate Swagger docs (after any endpoint changes):
   ```sh
   swag init -g cmd/main.go --output cmd/docs
   ```

4. Run the server:
   ```sh
   go run ./cmd/main.go
   ```

## Testing

```sh
go test ./...
```

---

## Generate Swagger Docs

```sh
swag init -g cmd/main.go --output cmd/docs
```
