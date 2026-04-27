# E-Shop API

A RESTful API for e-commerce built with Go using the Gin framework and PostgreSQL.

## Tech Stack

- **Language**: Go 1.25
- **Framework**: Gin
- **Database**: PostgreSQL
- **Pool**: DB connection pooling
- **ORM**: GORM
- **Migrations**: gormigrate
- **Authentication**: JWT (golang-jwt/jwt)
- **Utilities**: godotenv, uuid, air, gorm, golangcli-lint, redis, bcrypt (via golang.org/x/crypto), zap, gobreaker
- **HTTPS/TLS**: SSL certificate support for secure HTTPS connections

## Project Structure

```bash
e-shop-api/
├── cmd/
│   ├── api/            # Command main API server
│   ├── gen/            # Command migration generator
│   ├── migrate/        # Command database migration
│   └── seed/           # Command data seeder
├── docs/
│   ├── api/            # API documentation (HTML, YAML, PNG)
│   └── erd/            # Database documentation (DBML, SQL, ERD PNG)
├── internal/
│   ├── apps/           # Application setup (router, DI registries, etc.)
│   ├── configs/        # Configuration (database, redis, etc.)
│   ├── constants/      # Constants for common values
│   ├── dtos/           # Data Transfer Objects
│   ├── handlers/       # HTTP handlers
│   ├── middlewares/    # Middleware
│   ├── migrations/     # Database migrations (gormigrate)
│   ├── models/         # Database models
│   ├── repositories/   # Database repositories
│   ├── seeders/        # Data seeders
│   ├── services/       # Business logic services
│   └── pkg/utils/      # Utility packages (logger, auth, pagination, etc.)
├── uploads/            # Static file storage
├── .env.example        # Example environment file
├── docker-compose.yml  # Docker Compose file
├── go.mod              # Go module file
├── go.sum              # Go module checksum file
└── Makefile            # Project Makefile
```

## Architecture

This project implements a **Layered Architecture (3-Tier)** with principles inspired by **Clean Architecture** and **Dependency Injection (DI)**.

### Clean Architecture Principles

| Principle                   | Implementation                                                                     |
|-----------------------------|------------------------------------------------------------------------------------|
| **Separation of Concerns**  | Handler (I/O), Service (Business), Repository (Data) layers are strictly separated |
| **Dependency Inversion**    | Services depend on repository interfaces, not concrete implementations             |
| **Interface Segregation**   | Each repository has separate Write (`*Repo`) and Read (`*QueryRepo`) interfaces    |
| **Single Responsibility**   | Each layer has one job - handlers receive, services process, repositories persist  |

### Architecture Layers

```diagram
+-------------------------------------------------------------+
|                       Handler Layer                         |
|                    (internal/handlers)                      |
|   - HTTP request/response handling                          |
|   - Input validation & binding                              |
|   - Calls services, returns formatted responses             |
+-------------------------------------------------------------+
                            | ^
                            v |
+-------------------------------------------------------------+
|                       Service Layer                         |
|                    (internal/services)                      |
|   - Business logic & orchestration                          |
|   - Transaction management (Begin/Commit/Rollback)          |
|   - Depends on repository interfaces                        |
|   - Caches data in Redis                                    |
+-------------------------------------------------------------+
                            | ^
                            v |
+-------------------------------------------------------------+
|                      Repository Layer                       |
|                   (internal/repositories)                    |
|   - *Repository: Write operations (Create, Update)          |
|   - *QueryRepository: Read operations (Find, List)          |
|   - Database operations via GORM                            |
+-------------------------------------------------------------+
```

### Dependency Injection (DI)

The project uses a **Registry Pattern** to manage dependencies:

```go
// internal/apps/
NewRepositoryRegistry(db)    // Creates all repository instances
NewServiceRegistry(...)      // Injects repositories into services
NewHandlerRegistry(...)    // Injects services into handlers
RegisterRoutes(...)       // Wires up HTTP handlers
```

### Data Flow

```diagram
                HTTP Request
                      |
    [Handler] - Validates input, binds JSON
                      |
[Service] - Business logic, transaction management
                      |
    [Repository] - Database operations (GORM)
                      |
          [Database (PostgreSQL)]
                      ^
    [Response] - Formatted JSON via middleware
```

### Additional Patterns

| Pattern                 | Where Used                                                          |
|-------------------------|---------------------------------------------------------------------|
| **DTO Pattern**         | `internal/dtos/` - Request/Response objects separate from DB models |
| **Registry Pattern**    | `internal/apps/` - Centralized DI container                         |
| **Transaction Script**  | Services manage explicit DB transactions                            |
| **Factory Pattern**     | `NewXxxService()`, `NewXxxHandler()` constructors                   |

This architecture provides:

- **Testability** - Services can be mocked via interfaces
- **Maintainability** - Changes isolated to specific layers
- **Flexibility** - Easy to swap database or transport layers

## Features

- **User Authentication**: Register and login with JWT-based authentication
- **Forgot/Reset Password**: Password reset via email with Redis token storage (5-min TTL)
- **Upload File Handling**: Upload file handling with custom options and validation
- **RBAC**: Role-based access control
- **Store Management**: Create, update, delete, and activate/deactivate stores
- **Product Management**: CRUD operations for products with categories
- **Order Management**: Create, list, update, cancel, and confirm orders with order items
- **Email Notifications**: Automatic email notifications for order events (create, update, cancel, confirm), implement goroutine to send emails
- **Pagination**: Built-in pagination support for list endpoints
- **Custom Validation**: Request validation with custom validators
- **Custom Exceptions**: Custom exception handling
- **Error Handling**: Custom error responses with status codes
- **Response Formatting**: Custom response formatting
- **Transaction Support**: Database transactions for data integrity
- **Database Pooling**: Database connection pooling for improved performance
- **Redis Caching**: Redis caching for improved performance
- **Rate Limiting**: Request rate limiting using Redis with fixed window counter
  - `/auth/login`: 5 requests per 5 seconds
  - `/auth/forgot-password`: 3 requests per 1 minute
  - Headers: X-RateLimit-Limit, X-RateLimit-Remaining, Retry-After
  - Atomic INCR + EXPIRE for race-condition-free counting
- **Circuit Breaker**: Circuit breaker pattern using gobreaker to handle external service failures gracefully
- **Slow Query Tracker**: Track and log all database queries with execution timing (configurable threshold)
- **Auto Retry**: Automatic retry mechanism with configurable attempts and delay
- **HTTPS Support**: TLS/SSL support for secure HTTPS connections
- **Graceful Shutdown**: Graceful shutdown handling for server, DB, and Redis connections
- **Logging**: Structured logging with Zap

## Logging

This project uses **Zap** for structured logging. The logger is configured in `internal/pkg/logger/logger.go`.

### Initialization Logging

```go
func main() {
    logger.InitLogger()
    defer logger.L.Sync()
    // ... rest of your code
}
```

### Usage Logging

The logger exposes two aliases: `logger.Log` and `logger.L` (shorter syntax).

```go
import (
    "e-shop-api/internal/pkg/logger"
    "go.uber.org/zap"
)

// Info logging
logger.L.Info("Server started", zap.String("port", "8001"))

// Error logging
logger.L.Fatal("Failed to connect to database", zap.Error(err))

// Warning logging
logger.L.Warn("Rate limit exceeded", zap.String("ip", clientIP))
```

### Available Fields

| Field Type                | Usage                 |
|---------------------------|-----------------------|
| `zap.String(key, value)`  | Log a string field    |
| `zap.Int(key, value)`     | Log an integer field  |
| `zap.Error(err)`          | Log an error          |
| `zap.Any(key, value)`     | Log any type          |
| `zap.Bool(key, value)`    | Log a boolean         |
| `zap.Float64(key, value)` | Log a float           |

### Output Format

- **Development mode** (`APP_ENV=development`): Console output with colors
- **Production mode** (`APP_ENV=production`): JSON output for log aggregation

## Circuit Breaker

This project uses **gobreaker** (from sony/gobreaker) for circuit breaker pattern and auto-retry mechanism. The utility is configured in `internal/pkg/utils/circuit_breaker.go`.

### Initialization Circuit Breaker

```go
import (
    "e-shop-api/internal/pkg/utils"
)

// Create a circuit breaker
cb := util.NewCircuitBreaker("order-service")
```

### Usage Circuit Breaker

```go
// Circuit breaker usage
result, err := cb.Execute(func() (interface{}, error) {
    return orderService.CreateOrder(req, user)
})
if err != nil {
    if cb.IsStateOpen() {
        logger.L.Error("Circuit breaker is open, service unavailable")
    }
    return dto.OrderResponse{}, err
}
```

### Auto Retry

```go
// Auto retry with automatic function name detection
err := util.AutoRetry(func() error {
    return externalAPI.Call()
})
```

### Manual Retry Helper

```go
// Manual retry with custom function name
err := util.RetryHelper("create-order", 3, 2*time.Second, func() error {
    return orderService.CreateOrder(req, user)
})
```

### Configuration

Circuit breaker and retry can be configured via environment variables:

| Variable            | Default      | Description                                                           |
|---------------------|--------------|-----------------------------------------------------------------------|
| `CB_MAX_REQUESTS`   | 3            | Maximum requests allowed when circuit is half-open                    |
| `CB_INTERVAL`       | 5s (second)  | Period the circuit will stay open before testing if fault is resolved |
| `CB_TIMEOUT`        | 30s (second) | Time the circuit stays open before attempting to close                |
| `CB_THRESHOLD`      | 3            | Number of consecutive failures to trigger circuit open                |
| `RETRY_ATTEMPTS`    | 3            | Number of retry attempts                                              |
| `RETRY_DELAY`       | 2s (second)  | Delay between retry attempts                                          |

### State Change Logging

The circuit breaker logs state changes:

```code
WARN - Circuit Breaker State Change name=order-service from=closed to=open
```

**States:**

- **Closed**: Normal operation, requests pass through
- **Open**: Service unavailable, requests fail fast
- **Half-Open**: Testing if service recovered

## Slow Query Tracker

This project uses a GORM plugin to track and log all database queries with execution timing.

### Features of the Slow Query Tracker

- Tracks all database queries (SELECT, INSERT, UPDATE, DELETE)
- Logs query execution time in milliseconds
- Configurable slow query threshold via environment variable
- Slow queries (> threshold) are logged at WARN level with `[SLOW QUERY]` prefix
- Normal queries are logged at INFO level

### Configuration of the Slow Query Tracker

| Environment Variable    | Default | Description                                        |
|-------------------------|---------|----------------------------------------------------|
| `SLOW_QUERY_THRESHOLD`  | `200ms` | Threshold in milliseconds for slow query detection |

### Log Output Example

```code
INFO Query executed: 15ms - products
INFO Query executed: 50ms - orders
WARN [SLOW QUERY] Query executed: 245ms - products
```

### Implementation

The slow query tracker is implemented as a GORM plugin in `internal/pkg/querytracker/querytracker.go` and is automatically registered when the database connection is established.

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL 15+
- Docker & Docker Compose (optional)
- `air` for hot reloading (`go install github.com/air-verse/air@latest`)
- `golangci-lint` for linting (`curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.11.4`)
- `redis` for caching (`go get github.com/redis/go-redis/v9`)
- SSL certificates (for HTTPS, optional)

### 1. Clone and Setup

```bash
# Clone the repository
cd e-shop-api

# Copy environment file
cp .env.example .env
```

### 2. Configure Environment

Edit `.env` with your database and application settings:

```env
# APP
SERVER_PORT="<http_port>"
HTTPS_PORT="<https_port>"
SSL_CERT_PATH="<cert>.pem"
SSL_KEY_PATH="<cert>-key.pem"
APP_ENV="development"
USE_HTTPS=true

# DB
DB_HOST=127.0.0.1
DB_USER=<db_username>
DB_PASSWORD=<db_password>
DB_NAME=<db_name>
DB_PORT=<db_port>
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME="60m" #minutes
DB_CONN_MAX_IDLETIME="15m" #minutes

# SLOW QUERY TRACKER
SLOW_QUERY_THRESHOLD=200

# JWT
JWT_SECRET_KEY=<jwt_secret_key>
JWT_TTL="3600s" #seconds
JWT_ACCESS_TTL="900s" #seconds (access token - 15 min)
JWT_REFRESH_TTL="604800s" #seconds (refresh token - 7 days)

# SMTP EMAIL
SMTP_HOST=127.0.0.1
SMTP_PORT=1025
SMTP_SENDER_NAME="E-Shop Admin"
SMTP_AUTH_EMAIL=<auth_email>
SMTP_AUTH_PASSWORD=<auth_password>

# PROXIES
TRUSTED_PROXIES=127.0.0.1

# CORS
CORS_ALLOWED_ORIGINS=http://127.0.0.1:8001,https://127.0.0.1:8001,http://127.0.0.1:5500

# REDIS
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_CACHE_TTL="5m" #minutes

# CIRCUIT BREAKER
CB_MAX_REQUESTS=3
CB_INTERVAL="5s" #seconds
CB_TIMEOUT="30s" #seconds
CB_THRESHOLD=3
RETRY_ATTEMPTS=3
RETRY_DELAY="2s" #seconds
```

### 3. Start Database

Using Docker Compose:

```bash
docker-compose up -d
```

Or use an existing PostgreSQL instance.

### 4. Run Migrations

Run pending migrations:

```bash
make migrate
```

Generate a new migration for a model:

```bash
make add-migrate name=<ModelName>
```

### 5. Seed Data (Optional)

```bash
make seed
```

### 6. Run the Server

```bash
# Run the server
make run

# Run the server with hot reload
make dev
```

The server will start on `http://localhost:8001` if `USE_HTTPS` is set to `true` in `.env` server will start on `https://localhost:8001` (or the port specified in `.env`).
> **Note:** Before using command `make dev` you need to install `air` with `go install github.com/air-verse/air@latest`, then run command `air init`, and update `.air.toml` file with your configuration.

### 7. Linting (Run golangci-lint)

```bash
make lint
```

### 8. Tidy (Cleaning and verifying go.mod and go.sum)

```bash
make tidy
```

### 9. Clean (Delete binary and temp files)

```bash
make clean
```

## API Endpoints

### Public Routes

| Method | Endpoint                     | Description                                 |
|--------|------------------------------|---------------------------------------------|
| GET    | /health                      | Health check                                |
| GET    | /ready                       | Readiness check                             |
| POST   | /api/v1/auth/register        | Register new user                           |
| POST   | /api/v1/auth/login           | Login user (returns access + refresh token) |
| POST   | /api/v1/auth/refresh-token   | Refresh access token                        |
| POST   | /api/v1/auth/forgot-password | Request password reset                      |
| PUT    | /api/v1/auth/reset-password  | Reset password with token                   |

### Protected Routes (Requires JWT)

| Method | Endpoint                      | Description                |
|--------|-------------------------------|----------------------------|
| GET    | /api/v1/users/profile         | Get profile user (cached)  |
| POST   | /api/v1/users/upload-picture  | Upload profile picture     |
| POST   | /api/v1/stores                | Create store               |
| GET    | /api/v1/stores                | List stores (paginated)    |
| PUT    | /api/v1/stores/:id            | Update store               |
| PATCH  | /api/v1/stores/:id            | Delete store               |
| PATCH  | /api/v1/stores/activate       | Activate/deactivate store  |
| POST   | /api/v1/products              | Create product             |
| GET    | /api/v1/products              | List products (paginated)  |
| PUT    | /api/v1/products/:id          | Update product             |
| PATCH  | /api/v1/products/:id          | Delete product             |
| PATCH  | /api/v1/products/activate     | Activate/deactivate product|
| POST   | /api/v1/orders                | Create order               |
| GET    | /api/v1/orders                | List orders (paginated)    |
| PUT    | /api/v1/orders/:id            | Update order               |
| PATCH  | /api/v1/orders/:id/cancel     | Cancel order               |
| PATCH  | /api/v1/orders/:id/confirm    | Confirm order              |

## Request/Response Examples

### Register

```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

### Login

```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Forgot Password

```bash
POST /api/v1/auth/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

### Reset Password

```bash
PUT /api/v1/auth/reset-password
Content-Type: application/json

{
  "token": "uuid-token-from-email",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```

### Create Store (Protected)

```bash
POST /api/v1/stores
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "My Store",
  "description": "Store description"
}
```

### Create Product (Protected)

```bash
POST /api/v1/products
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Product Name",
  "description": "Product description",
  "price": 99.99,
  "stock": 100,
  "category_id": "uuid-of-category"
}
```

### List Products (Paginated)

```bash
GET /api/v1/products?page=1&limit=10
Authorization: Bearer <token>
```

Response:

```json
{
  "items": [...],
  "meta": {
    "current_page": 1,
    "total_page": 5,
    "total_data": 50,
    "limit": 10
  }
}
```

### Create Order (Protected)

```bash
POST /api/v1/orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "store_id": "uuid-of-store",
  "items": [
    {
      "product_id": "uuid-of-product",
      "quantity": 2
    }
  ]
}
```

## Documentation

### API Documentation

The API documentation is available as an interactive HTML collection in `doc/api/e-shop-api-documentation.html`.

**How to view:**

1. Install the **Live Server** extension in VS Code (by Ritwick Dey)
2. Right-click `doc/api/e-shop-api-documentation.html`
3. Select "Open with Live Server"

The documentation will open in your default browser, allowing you to test all API endpoints with the configured environment (base URL: `http://localhost:8001`).

### Database Schema (DBML)

The database schema is defined in `doc/erd/e_shop_db.dbml` using DBML (Database Markup Language).

**How to view:**

1. Install a DBML viewer extension in VS Code (e.g., **DBML** by mohsen1)
2. Open `doc/erd/e_shop_db.dbml` to see an interactive ERD diagram
3. Alternatively, view the pre-generated ERD image at `doc/erd/e_shop_db_erd.png`

The SQL schema export is also available at `doc/erd/e_shop_db.sql`.

## License

MIT
