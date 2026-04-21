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
│   ├── api/            # Main API server
│   ├── migrate/        # Database migration
│   └── seed/           # Data seeding
├── internal/
│   ├── app/            # Application setup (router, DI registries)
│   ├── config/         # Configuration (database, migration, seeder)
│   ├── dto/            # Data Transfer Objects
│   ├── handler/        # HTTP handlers
│   ├── middleware/     # Middleware (auth, response)
│   ├── model/          # Database models
│   ├── repository/     # Database repositories
│   ├── service/        # Business logic services
│   └── pkg/util/       # Utility packages (JWT, exceptions)
├── docker-compose.yml
├── .env.example
└── go.mod
```

## Features

- **User Authentication**: Register and login with JWT-based authentication
- **Store Management**: Create, update, delete, and activate/deactivate stores
- **Product Management**: CRUD operations for products with categories
- **Order Management**: Create, list, update, cancel, and confirm orders with order items
- **Pagination**: Built-in pagination support for list endpoints
- **Custom Validation**: Request validation with custom validators
- **Transaction Support**: Database transactions for data integrity

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

| Method | Endpoint                  | Description           |
|--------|---------------------------|-----------------------|
| POST   | /api/v1/auth/register    | Register new user    |
| POST   | /api/v1/auth/login       | Login user           |

### Protected Routes (Requires JWT)

| Method | Endpoint                      | Description              |
|--------|-------------------------------|--------------------------|
| POST   | /api/v1/stores                | Create store            |
| GET    | /api/v1/stores                | List stores (paginated)|
| PUT    | /api/v1/stores/:id            | Update store           |
| PATCH  | /api/v1/stores/:id            | Delete store           |
| PATCH  | /api/v1/stores/activate       | Activate/deactivate store|
| POST   | /api/v1/products              | Create product          |
| GET    | /api/v1/products              | List products (paginated)|
| PUT    | /api/v1/products/:id          | Update product          |
| PATCH  | /api/v1/products/:id          | Delete product          |
| PATCH  | /api/v1/products/activate     | Activate/deactivate product|
| POST   | /api/v1/orders                | Create order            |
| GET    | /api/v1/orders                | List orders (paginated)|
| PUT    | /api/v1/orders/:id            | Update order            |
| PATCH  | /api/v1/orders/:id/cancel     | Cancel order           |
| PATCH  | /api/v1/orders/:id/confirm    | Confirm order          |

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
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_items": 50,
    "total_pages": 5
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

## License

MIT
