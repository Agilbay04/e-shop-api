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
├── doc/
│   ├── api/            # API documentation (HTML, YAML, PNG)
│   └── erd/            # Database documentation (DBML, SQL, ERD PNG)
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
# APP
SERVER_PORT=<server_port>
APP_ENV=development

# DB
DB_HOST=localhost
DB_USER=<db_username>
DB_PASSWORD=<db_password>
DB_NAME=<db_name>
DB_PORT=<db_port>

# JWT
JWT_SECRET_KEY=<jwt_secret_key>

# SMTP EMAIL
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_SENDER_NAME="E-Shop Admin"
SMTP_AUTH_EMAIL=<auth_email>
SMTP_AUTH_PASSWORD=<auth_password>
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
