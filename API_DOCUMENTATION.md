# User API Documentation

A REST API for managing employees, departments, attendance tracking, payroll processing, and attendance reports.

**Base URL:** `http://localhost:8080/api/v1`

**Swagger UI:** `http://localhost:8080/swagger/index.html`

---

## Authentication

Most endpoints require JWT authentication. Include the token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

### Endpoints (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register a new user |
| POST | `/auth/login` | Authenticate and get JWT token |

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
  "department": { "id": 1, "name": "Engineering", "description": "..." },
  "salary": 50000.00,
  "hire_date": "2024-01-15"
}
```

### Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/workers` | Required | List all workers |
| GET | `/workers/:id` | Required | Get worker by ID |
| POST | `/workers` | Required | Create a new worker |
| PUT | `/workers` | Required | Update a worker |
| DELETE | `/workers/:id` | Required | Delete a worker |

#### Create Worker

**Request Body:**
```json
{
  "name": "Jane Smith",
  "email": "jane@example.com",
  "department_id": 2
}
```

**Validation:** `name` (required, 3-100 chars), `email` (required, valid email format)

#### Update Worker

**Request Body:**
```json
{
  "id": 1,
  "name": "John Doe Updated",
  "email": "john.updated@example.com"
}
```

**Validation:** `id`, `name`, `email` are required.

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

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/departments` | Required | List all departments |
| GET | `/departments/:id` | Required | Get department by ID |
| POST | `/departments` | Required | Create a new department |
| PUT | `/departments` | Required | Update a department |
| DELETE | `/departments/:id` | Required | Delete a department |

#### Create Department

**Request Body:**
```json
{
  "name": "Marketing",
  "description": "Marketing and sales team"
}
```

**Validation:** `name` (required)

#### Update Department

**Request Body:**
```json
{
  "id": 1,
  "name": "Engineering & Tech",
  "description": "Updated description"
}
```

**Validation:** `id` and `name` (required)

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

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/attendances` | Required | List all attendance records |
| GET | `/attendances/:id` | Required | Get attendance by ID |
| GET | `/attendances/worker/:worker_id` | Required | Get attendance for a worker |
| POST | `/attendances` | Required | Create attendance record (check-in) |
| PUT | `/attendances/:id` | Required | Update attendance (check-out) |

#### Create Attendance (Check-In)

**Request Body:**
```json
{
  "worker_id": 1,
  "date": "2024-01-15",
  "check_in": "09:00:00",
  "check_out": "17:00:00"
}
```

**Validation:** `worker_id`, `date`, `check_in`, `check_out` (all required)

#### Update Attendance (Check-Out)

**Request Body:**
```json
{
  "id": 1,
  "check_out": "18:00:00"
}
```

**Validation:** `id` and `check_out` (required)

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

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/payrolls` | Required | List all payroll records |
| GET | `/payrolls/:workerId` | Required | Get payroll for a worker |
| POST | `/payroll/calculate` | Required | Calculate payroll (base + bonus - deductions) |

#### Calculate Payroll

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

**Response:**
```json
{
  "net_salary": 5300.00,
  "message": "Payroll calculated successfully"
}
```

**Validation:** `worker_id` and `month` (required)

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

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/reports/workers/attendance` | Required | Get attendance report for workers |

#### Query Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `start_date` | string | Yes | Start date (YYYY-MM-DD) |
| `end_date` | string | Yes | End date (YYYY-MM-DD) |
| `department_id` | integer | No | Filter by department |
| `worker_id` | integer | No | Filter by worker |

**Example:**
```
GET /reports/workers/attendance?start_date=2024-01-01&end_date=2024-01-31&department_id=1
```

---

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Email is required"
}
```

### Common HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request (validation error) |
| 401 | Unauthorized (missing/invalid token) |
| 404 | Not Found |
| 500 | Internal Server Error |

---

## Project Structure

```
user_api/
├── cmd/
│   ├── main.go              # Application entry point
│   └── docs/                 # Swagger generated docs
├── internal/
│   ├── attendance/           # Attendance module
│   │   ├── attendance_handler.go
│   │   ├── attendance_model.go
│   │   ├── attendance_repository.go
│   │   └── attendance_service.go
│   ├── auth/                 # Authentication module
│   │   ├── auth_handler.go
│   │   ├── auth_service.go
│   │   ├── jwt_service.go
│   │   └── user_model.go
│   ├── departments/          # Departments module
│   │   ├── department_handler.go
│   │   ├── department_model.go
│   │   └── department_repository.go
│   ├── dto/                  # Data Transfer Objects
│   │   ├── attendance/
│   │   ├── auth/
│   │   ├── department/
│   │   ├── payroll/
│   │   ├── report/
│   │   └── worker/
│   ├── middleware/          # Auth & logging middleware
│   ├── payroll/              # Payroll module
│   │   ├── payroll_handler.go
│   │   ├── payroll_model.go
│   │   ├── payroll_repository.go
│   │   └── payroll_service.go
│   ├── reports/              # Reports module
│   │   ├── report_handler.go
│   │   └── report_service.go
│   ├── workers/              # Workers module
│   │   ├── worker_handler.go
│   │   ├── worker_model.go
│   │   ├── worker_repository.go
│   │   └── worker_service.go
│   ├── common/               # Shared utilities
│   └── database/             # Database connection
└── readme.md
```

## Setup

1. Install dependencies:
   ```sh
   go mod download
   ```

2. Create `.env` file:
   ```env
   DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable
   JWT_SECRET=your_secret_here
   ```

3. Run the server:
   ```sh
   go run ./cmd/main.go
   ```

## Testing

```sh
go test ./...
```

## Generate Swagger Docs

```sh
swag init -g cmd/main.go --output cmd/docs
```
