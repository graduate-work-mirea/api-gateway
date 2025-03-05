package middlewares

import (
    "log"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "api-gateway/internal/cache"
    "api-gateway/internal/grpc"
)

func AuthMiddleware(client *grpc.AuthClient, tokenCache *cache.TokenCache) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token from Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
            return
        }
        token := parts[1]

        // Check cache first
        if valid, exists := tokenCache.Get(token); exists {
            if !valid {
                c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
                return
            }
            c.Next()
            return
        }

        // Validate token via gRPC
        valid, err := client.ValidateToken(token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
            log.Printf("Token validation error: %v", err)
            return
        }

        // Cache the result
        tokenCache.Set(token, valid)

        if !valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        c.Next()
    }
}
