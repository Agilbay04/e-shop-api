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
- **Utilities**: godotenv, uuid, air, gorm, golangcli-lint, redis, bcrypt (via golang.org/x/crypto), zap

## Project Structure

```bash
e-shop-api/
├── cmd/
│   ├── api/            # Main API server
│   ├── gen/            # Migration generator
│   ├── migrate/        # Database migration
│   └── seed/           # Data seeding
├── doc/
│   ├── api/            # API documentation (HTML, YAML, PNG)
│   └── erd/            # Database documentation (DBML, SQL, ERD PNG)
├── internal/
│   ├── app/            # Application setup (router, DI registries)
│   ├── config/         # Configuration (database, seeder, redis)
│   ├── dto/            # Data Transfer Objects
│   ├── handler/        # HTTP handlers
│   ├── middleware/     # Middleware
│   ├── migrations/     # Database migrations (gormigrate)
│   ├── model/          # Database models
│   ├── repository/     # Database repositories
│   ├── service/        # Business logic services
│   └── pkg/util/       # Utility packages (logger, auth, pagination, etc.)
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
|                      Handler Layer                          |
|                     (internal/handler)                      |
|   - HTTP request/response handling                          |
|   - Input validation & binding                              |
|   - Calls services, returns formatted responses             |
+-------------------------------------------------------------+
                          | ^
                          v |
+-------------------------------------------------------------+
|                      Service Layer                          |
|                     (internal/service)                      |
|   - Business logic & orchestration                          |
|   - Transaction management (Begin/Commit/Rollback)          |
|   - Depends on repository interfaces                        |
+-------------------------------------------------------------+
                          | ^
                          v |
+-------------------------------------------------------------+
|                    Repository Layer                         |
|                    (internal/repository)                    |
|   - *Repository: Write operations (Create, Update)          |
|   - *QueryRepository: Read operations (Find, List)          |
|   - Database operations via GORM                            |
+-------------------------------------------------------------+
```

### Dependency Injection (DI)

The project uses a **Registry Pattern** to manage dependencies:

```go
// internal/app/
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

| Pattern                 | Where Used                                                         |
|-------------------------|--------------------------------------------------------------------|
| **DTO Pattern**         | `internal/dto/` - Request/Response objects separate from DB models |
| **Registry Pattern**    | `internal/app/` - Centralized DI container                         |
| **Transaction Script**  | Services manage explicit DB transactions                           |
| **Factory Pattern**     | `NewXxxService()`, `NewXxxHandler()` constructors                  |

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
- **Rate Limiting**: Request rate limiting using Redis (5 req/5s for login, 1 req/1min for forgot-password)
- **Graceful Shutdown**: Graceful shutdown handling for server, DB, and Redis connections
- **Logging**: Structured logging with Zap

## Logging

This project uses **Zap** for structured logging. The logger is configured in `internal/pkg/logger/logger.go`.

### Initialization

```go
func main() {
    logger.InitLogger()
    defer logger.L.Sync()
    // ... rest of your code
}
```

### Usage

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
- **Production mode**: JSON output for log aggregation

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL 15+
- Docker & Docker Compose (optional)
- `air` for hot reloading (`go install github.com/air-verse/air@latest`)
- `golangci-lint` for linting (`curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.11.4`)
- `redis` for caching (`go get github.com/redis/go-redis/v9`)

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
SERVER_PORT=<server_port>
APP_ENV=development

# DB
DB_HOST=localhost
DB_USER=<db_username>
DB_PASSWORD=<db_password>
DB_NAME=<db_name>
DB_PORT=<db_port>
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME=60 #minutes
DB_CONN_MAX_IDLETIME=15 #minutes

# JWT
JWT_SECRET_KEY=<jwt_secret_key>
JWT_TTL=3600 #seconds

# SMTP EMAIL
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_SENDER_NAME="E-Shop Admin"
SMTP_AUTH_EMAIL=<auth_email>
SMTP_AUTH_PASSWORD=<auth_password>

# PROXIES
TRUSTED_PROXIES=127.0.0.1

# CORS
CORS_ALLOWED_ORIGINS=http://127.0.0.1:8001,http://127.0.0.1:5500

# REDIS
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_CACHE_TTL=5 #minutes
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

The server will start on `http://localhost:8001` (or the port specified in `.env`).
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

| Method | Endpoint                     | Description               |
|--------|------------------------------|---------------------------|
| POST   | /api/v1/auth/register        | Register new user         |
| POST   | /api/v1/auth/login           | Login user                |
| POST   | /api/v1/auth/forgot-password | Request password reset    |
| PUT    | /api/v1/auth/reset-password  | Reset password with token |

### Protected Routes (Requires JWT)

| Method | Endpoint                      | Description                |
|--------|-------------------------------|----------------------------|
| GET    | /api/v1/auth/profile          | Get profile user           |
| POST   | /api/v1/auth/upload-picture   | Upload profile picture     |
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
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {...}
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
