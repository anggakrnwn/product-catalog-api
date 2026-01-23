# product-catalog-api

A RESTful API for managing product categories built with Go and Gin framework. This project implements clean architecture with in-memory storage and is developed as part of CodeWithUmam.

## Features

- **Full Operations** - Create, Read, Update, Delete categories
- **Clean Architecture** - Separation of concerns with domain, usecase, repository layers
- **In-Memory Storage** - Fast development and testing without database setup
- **RESTful Design** - Standard HTTP methods with proper status codes
- **API Versioning** - Organized under `/api/v1` endpoint
- **CORS Support** - Cross-origin requests enabled
- **Health Check** - Endpoint for monitoring and deployment verification
- **Bulk Operations** - Create multiple categories in single request
- **Input Validation** - Request validation with proper error messages
- **Unit Tests** - Test coverage for business logic

## Project Structure

```
product-catalog-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── delivery/
│   │   └── http/
│   │       └── category_handler.go  # HTTP handlers
│   ├── domain/
│   │   └── category.go          # Domain model
│   ├── dto/
│   │   └── category_dto.go      # Data Transfer Objects
│   ├── middleware/
│   │   └── cors.go              # CORS middleware
│   ├── repository/
│   │   ├── category_repository.go      # Repository interface
│   │   └── inmemory/
│   │       └── category_repository_inmemory.go  # In-memory implementation
│   ├── seed/
│   │   └── seed.go              # Data seeder
│   └── usecase/
│       ├── interfaces.go        # Usecase interfaces
│       ├── category_usecase.go  # Business logic
│       └── category_usecase_test.go  # Unit tests
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/anggakrnwn/product-catalog-api.git
cd product-catalog-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### Health Check
```
GET /health
```
Response:
```json
{
  "status": "healthy",
  "storage": "in-memory"
}
```

#### Categories

##### Get All Categories
```
GET /api/v1/categories
```
Response:
```json
{
  "success": true,
  "data": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Elektronik",
      "description": "Perangkat elektronik dan gadget",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "meta": {
    "count": 5,
    "total": 5
  }
}
```

##### Get Single Category
```
GET /api/v1/categories/{id}
```

##### Create Category
```
POST /api/v1/categories
```
Request:
```json
{
  "name": "Fashion",
  "description": "Pakaian dan aksesori"
}
```
Response (201 Created):
```json
{
  "success": true,
  "message": "Category created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Fashion",
    "description": "Pakaian dan aksesori",
    "created_at": "2024-01-15T11:45:00Z",
    "updated_at": "2024-01-15T11:45:00Z"
  }
}
```

##### Update Category
```
PUT /api/v1/categories/{id}
```
Request:
```json
{
  "name": "Fashion Updated",
  "description": "Updated description"
}
```

##### Delete Category
```
DELETE /api/v1/categories/{id}
```

##### Bulk Create Categories
```
POST /api/v1/categories/bulk
```
Request:
```json
{
  "categories": [
    {"name": "Elektronik", "description": "Gadget"},
    {"name": "Buku", "description": "Bacaan"}
  ]
}
```

## Testing

Run unit tests:
```bash
go test ./internal/usecase/...
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `ENVIRONMENT` | `development` | Runtime environment |
| `SEED_DATA` | `true` | Whether to seed initial data |

Set environment variables:
```bash
export PORT=3000
export ENVIRONMENT=production
export SEED_DATA=false
```

## Contributing

This project is developed as part of CodeWithUmam. While contributions are welcome, please note this is primarily a learning project.

## License

This project is created for educational purposes as part of CodeWithUmam.

---
**Built with ❤️ using Go & Gin**