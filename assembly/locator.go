package assembly

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/graduate-work-mirea/api-gateway/config"
	"github.com/graduate-work-mirea/api-gateway/controller"
	"github.com/graduate-work-mirea/api-gateway/repository"
	"github.com/graduate-work-mirea/api-gateway/server"
	"github.com/graduate-work-mirea/api-gateway/service"
)

// ServiceLocator is a container for all services
type ServiceLocator struct {
	config     *config.Config
	dbRepo     repository.DBRepository
	cacheRepo  repository.CacheRepository
	service    service.Service
	controller *controller.Controller
	server     *server.Server
	router     *gin.Engine
}

// NewServiceLocator creates a new service locator
func NewServiceLocator(cfg *config.Config) *ServiceLocator {
	// Create a new service locator
	locator := &ServiceLocator{
		config: cfg,
		router: gin.Default(),
	}

	log.Println("Creating new service locator")

	// Initialize services
	locator.initServices()

	return locator
}

// initServices initializes all services
func (l *ServiceLocator) initServices() {
	log.Println("Initializing services...")

	// Create database repository
	log.Println("Creating database repository...")
	dbRepo, err := repository.NewPostgreRepository(l.config)
	if err != nil {
		log.Fatalf("Failed to create database repository: %v", err)
	}
	l.dbRepo = dbRepo
	log.Println("Database repository created successfully")

	// Create cache repository
	log.Println("Creating cache repository...")
	cacheRepo, err := repository.NewCacheRepository(l.config)
	if err != nil {
		log.Fatalf("Failed to create cache repository: %v", err)
	}
	l.cacheRepo = cacheRepo
	log.Println("Cache repository created successfully")

	// Create service
	log.Println("Creating service layer...")
	l.service = service.NewService(l.config, l.dbRepo, l.cacheRepo)
	log.Println("Service layer created successfully")

	// Create controller and pass the router
	log.Println("Creating controller...")
	l.controller = controller.NewController(l.service, l.router)
	log.Println("Controller created successfully")

	// Create server and pass the router
	log.Println("Creating server...")
	l.server = server.NewServer(l.config, l.controller, l.router)
	log.Println("Server created successfully")

	log.Println("All services initialized successfully")
}

// GetService returns the service
func (l *ServiceLocator) GetService() service.Service {
	return l.service
}

// GetDBRepository returns the database repository
func (l *ServiceLocator) GetDBRepository() repository.DBRepository {
	return l.dbRepo
}

// GetCacheRepository returns the cache repository
func (l *ServiceLocator) GetCacheRepository() repository.CacheRepository {
	return l.cacheRepo
}

// GetController returns the controller
func (l *ServiceLocator) GetController() *controller.Controller {
	return l.controller
}

// GetServer returns the server
func (l *ServiceLocator) GetServer() *server.Server {
	return l.server
}

// GetRouter returns the router
func (l *ServiceLocator) GetRouter() *gin.Engine {
	return l.router
}
