package server

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graduate-work-mirea/api-gateway/config"
	"github.com/graduate-work-mirea/api-gateway/controller"
	"github.com/graduate-work-mirea/api-gateway/middleware"
)

// Server represents the HTTP server
type Server struct {
	router     *gin.Engine
	config     *config.Config
	controller *controller.Controller
}

// NewServer creates a new server
func NewServer(cfg *config.Config, controller *controller.Controller, router *gin.Engine) *Server {
	log.Println("Server: Configuring router with middleware...")

	// Add recovery middleware
	router.Use(gin.Recovery())

	// Configure CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.CorsOrigin}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(corsConfig))

	log.Println("Server: CORS configured with origins:", cfg.CorsOrigin)

	return &Server{
		router:     router,
		config:     cfg,
		controller: controller,
	}
}

// SetupRoutes sets up all routes
func (s *Server) SetupRoutes() {
	log.Println("Server: Setting up routes...")

	// Create auth middleware
	authMiddleware := middleware.AuthMiddleware(s.config)
	log.Println("Server: Auth middleware created")

	// Register routes
	s.controller.RegisterRoutes(authMiddleware)
	log.Println("Server: Controller routes registered")

	// Add a health check endpoint
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	log.Println("Server: Health check endpoint added")

	// Log all registered routes
	routes := s.router.Routes()
	log.Println("Server: Registered routes:")
	for _, route := range routes {
		log.Printf("  %s %s", route.Method, route.Path)
	}
}

// Start starts the server
func (s *Server) Start() error {
	// Set up routes
	s.SetupRoutes()

	// Start the server
	addr := fmt.Sprintf(":%s", s.config.Server.Port)
	log.Printf("Server: Starting on %s", addr)
	return s.router.Run(addr)
}

// GetRouter returns the router
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
