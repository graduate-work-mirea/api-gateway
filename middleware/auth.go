package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/graduate-work-mirea/api-gateway/config"
)

// JWTClaims represents the claims in a JWT
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware creates a middleware for authentication
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	log.Println("Middleware: Creating authentication middleware")
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		log.Printf("Middleware: Processing authentication for path: %s", path)

		// Get the JWT token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("Middleware: Authorization header is missing for path: %s", path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// Check if the header starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("Middleware: Invalid authorization header format for path: %s", path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		// Parse the JWT token
		tokenString := parts[1]
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Middleware: Unexpected signing method: %v for path: %s", token.Header["alg"], path)
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		// Check if the token is valid
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				log.Printf("Middleware: Token has expired for path: %s", path)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
				return
			}
			log.Printf("Middleware: Invalid token: %v for path: %s", err, path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Check if the token is valid
		if !token.Valid {
			log.Printf("Middleware: Token is invalid for path: %s", path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Check if the token has expired
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			log.Printf("Middleware: Token has expired for path: %s", path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
			return
		}

		// Parse UUID
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			log.Printf("Middleware: Invalid user ID: %v for path: %s", err, path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
			return
		}

		// Store the user ID in the context
		c.Set("userID", userID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		log.Printf("Middleware: Authentication successful for user: %s, role: %s, path: %s", claims.UserID, claims.Role, path)
		c.Next()
	}
}

// GetUserID gets the user ID from the context
func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("Middleware: User ID not found in context")
		return uuid.Nil, errors.New("user ID not found in context")
	}

	return userID.(uuid.UUID), nil
}
