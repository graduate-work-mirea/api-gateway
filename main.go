package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/graduate-work-mirea/api-gateway/assembly"
	"github.com/graduate-work-mirea/api-gateway/config"
)

func main() {
	// Configure logging
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Println("API Gateway starting...")
	log.Printf("Go version: %s, OS: %s, Arch: %s", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	// Log environment variables
	log.Println("Environment configuration:")
	log.Printf("  SERVER_PORT: %s", os.Getenv("SERVER_PORT"))
	log.Printf("  AUTH_SERVICE_HOST: %s", os.Getenv("AUTH_SERVICE_HOST"))
	log.Printf("  AUTH_SERVICE_PORT: %s", os.Getenv("AUTH_SERVICE_PORT"))
	log.Printf("  ML_SERVICE_HOST: %s", os.Getenv("ML_SERVICE_HOST"))
	log.Printf("  ML_SERVICE_PORT: %s", os.Getenv("ML_SERVICE_PORT"))

	// Load configuration from environment
	startTime := time.Now()
	log.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded in %v", time.Since(startTime))

	// Create service locator
	log.Println("Creating service locator...")
	locatorStartTime := time.Now()
	locator := assembly.NewServiceLocator(cfg)
	log.Printf("Service locator created in %v", time.Since(locatorStartTime))

	// Start the server
	log.Println("Starting the server...")
	server := locator.GetServer()
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
