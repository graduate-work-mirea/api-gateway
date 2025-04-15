FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway .

# Create a smaller image for running the application
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/api-gateway .

# Expose the application port
EXPOSE 8000

# Set the JWT secret as an environment variable
ENV JWT_SECRET=your_secret_key_here
ENV SERVER_PORT=8000
ENV AUTH_SERVICE_HOST=auth-service
ENV AUTH_SERVICE_PORT=8080
ENV ML_SERVICE_HOST=ml-service
ENV ML_SERVICE_PORT=8080
ENV POSTGRES_HOST=postgres
ENV POSTGRES_PORT=5432
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres
ENV POSTGRES_DB=marketplace_data
ENV POSTGRES_SSLMODE=disable
ENV CACHE_SIZE=1000
ENV CORS_ORIGIN=http://localhost

# Run the application
CMD ["./api-gateway"] 