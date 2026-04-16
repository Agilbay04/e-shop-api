# E-Shop API

A RESTful API for e-commerce built with Go using the Gin framework and PostgreSQL.

## Tech Stack

- **Language**: Go 1.25
- **Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Utilities**: godotenv, uuid, bcrypt (via golang.org/x/crypto)

## Project Structure

```bash
e-shop-api/
├── cmd/
│   ├── api/          # Main API server
│   ├── migrate/     # Database migration
│   └── seed/        # Data seeding
├── internal/
│   ├── app/         # Application setup (router, DI registries)
│   ├── config/      # Configuration (database, migration, seeder)
│   ├── dto/         # Data Transfer Objects
│   ├── handler/     # HTTP handlers
│   ├── middleware/  # Middleware (auth, response)
│   ├── model/       # Database models
│   ├── repository/  # Database repositories
│   ├── service/     # Business logic services
│   └── pkg/util/    # Utility packages (JWT, exceptions)
├── docker-compose.yml
├── .env.example
└── go.mod
```

## Features

- **User Authentication**: Register and login with JWT-based authentication
- **Store Management**: Create and manage stores
- **Product Management**: Manage products (model layer ready)

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

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
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=e_shop_db
DB_PORT=5432
SERVER_PORT=8001
APP_ENV=development
JWT_SECRET_KEY=your_jwt_secret
```

### 3. Start Database

Using Docker Compose:
```bash
docker-compose up -d
```

Or use an existing PostgreSQL instance.

### 4. Run Migrations

```bash
go run cmd/migrate/main.go
```

### 5. Seed Data (Optional)

```bash
go run cmd/seed/main.go
```

### 6. Run the Server

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8001` (or the port specified in `.env`).

## API Endpoints

### Public Routes

| Method | Endpoint         | Description       |
|--------|------------------|-------------------|
| POST   | /api/v1/auth/register | Register new user |
| POST   | /api/v1/auth/login    | Login user        |

### Protected Routes (Requires JWT)

| Method | Endpoint   | Description   |
|--------|------------|---------------|
| POST   | /api/v1/store | Create store  |

*Note: Additional endpoints for products and stores are available in the codebase.*

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

### Create Store (Protected)
```bash
POST /api/v1/store
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "My Store",
  "description": "Store description"
}
```

## License

MIT
