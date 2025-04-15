# API Gateway Service

This API Gateway serves as a central entry point for the authentication and ML prediction services in the marketplace system.

## Features

- Proxies authentication requests to the Auth Service
- Proxies ML prediction requests to the ML Service
- Stores prediction history in PostgreSQL database
- Maintains a local cache for faster access to prediction data
- Provides additional statistics endpoint for user prediction history

## API Documentation

API documentation is available in OpenAPI format at `docs/swagger.yaml`. You can visualize it using tools like Swagger UI or Redoc.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Docker and Docker Compose

### Environment Variables

The service can be configured using the following environment variables:

- `SERVER_PORT`: Port for the API Gateway (default: 8000)
- `AUTH_SERVICE_HOST`: Host for the Auth Service (default: localhost)
- `AUTH_SERVICE_PORT`: Port for the Auth Service (default: 8080)
- `ML_SERVICE_HOST`: Host for the ML Service (default: localhost)
- `ML_SERVICE_PORT`: Port for the ML Service (default: 6785)
- `POSTGRES_HOST`: Host for the PostgreSQL database (default: localhost)
- `POSTGRES_PORT`: Port for the PostgreSQL database (default: 5432)
- `POSTGRES_USER`: Username for the PostgreSQL database (default: postgres)
- `POSTGRES_PASSWORD`: Password for the PostgreSQL database (default: postgres)
- `POSTGRES_DB`: Database name for PostgreSQL (default: marketplace_data)
- `POSTGRES_SSLMODE`: SSL mode for PostgreSQL connection (default: disable)
- `CACHE_SIZE`: Size of the LRU cache (default: 1000)
- `JWT_SECRET`: Secret key for JWT token validation (default: your_secret_key_here)
- `CORS_ORIGIN`: Allowed CORS origin (default: http://localhost)

### Running with Docker Compose

The easiest way to run the entire system is to use Docker Compose:

```bash
docker-compose up -d
```

This will start the API Gateway, Auth Service, ML Service, PostgreSQL database, and the UI Service.

### Running Locally

To run the API Gateway locally:

1. Clone the repository
2. Install dependencies:
```bash
go mod download
```
3. Build the application:
```bash
go build -o api-gateway .
```
4. Run the application:
```bash
./api-gateway
```

### Testing the API

You can use the provided REST file at `docs/api.rest` to test the API endpoints. This file can be used with tools like VS Code's REST Client extension or Postman.

## Project Structure

- `assembly`: Service locator for dependency injection
- `config`: Configuration loading and management
- `controller`: HTTP request handlers
- `docs`: API documentation and testing files
- `middleware`: Authentication middleware
- `model`: Data models and structures
- `repository`: Database and cache repositories
- `server`: HTTP server implementation
- `service`: Business logic implementation

## License

This project is licensed under the MIT License - see the LICENSE file for details. 