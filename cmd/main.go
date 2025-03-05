package main

import (
    "log"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "api-gateway/internal/cache"
    "api-gateway/internal/grpc"
    "api-gateway/internal/handlers"
    "api-gateway/internal/middlewares"
)

func main() {
    // Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, relying on system env vars")
    }

    // Initialize gRPC client for Auth Service
    grpcClient, err := grpc.NewAuthClient(os.Getenv("AUTH_GRPC_ADDR"))
    if err != nil {
        log.Fatalf("Failed to connect to Auth Service gRPC: %v", err)
    }
    defer grpcClient.Close()

    // Initialize token cache
    tokenCache := cache.NewTokenCache()
    
    // Запуск периодической очистки кэша токенов
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        defer ticker.Stop()
        
        for range ticker.C {
            tokenCache.Cleanup()
        }
    }()

    // Set up Gin router
    r := gin.Default()

    // Initialize handlers with dependencies
    authHandler := handlers.NewAuthHandler(os.Getenv("AUTH_HTTP_ADDR"))
    analyticsHandler := handlers.NewAnalyticsHandler()

    // Public routes
    r.POST("/login", authHandler.Login)
    r.POST("/register", authHandler.Register)

    // Protected routes with authentication middleware
    protected := r.Group("/")
    protected.Use(middlewares.AuthMiddleware(grpcClient, tokenCache))
    {
        protected.GET("/analytics/demand", analyticsHandler.GetDemand)
    }

    // Start server
    port := os.Getenv("API_GATEWAY_PORT")
    if port == "" {
        port = "8081"
    }
    log.Printf("API Gateway starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
