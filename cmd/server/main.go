package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
    "github.com/joho/godotenv"

	"gitlab.com/pedrokoblitz/opus-go/actions"
	"gitlab.com/pedrokoblitz/opus-go/modules/auth"
	"gitlab.com/pedrokoblitz/opus-go/modules/story"
	"gitlab.com/pedrokoblitz/opus-go/services"
)

func main() {
	resourcesDir := "./resources"
	config, err := actions.LoadServiceConfig(resourcesDir + "/service.yml")
	if err != nil {
		log.Fatal("Failed to load service config:", err)
	}

	// Initialize container with proper error handling
	container, err := actions.NewContainer(config)
	if err != nil {
		log.Fatal("Failed to initialize container:", err)
	}
	// defer func() {
	// 	if err := container.Close(); err != nil {
	// 		log.Printf("Error closing container: %v", err)
	// 	}
	// }()

	// register providers
	auth.Register(container.Registry())
	story.Register(container.Registry())

	// Initialize services using the container
	hub := services.NewPubSubHub(container)
	httpSvc := services.NewHTTPService(container, hub)

	// Create context for graceful shutdown
	shutdownCtx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start HTTP server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting HTTP server on :%d", container.Config().HTTP.Port)
		if err := httpSvc.Start(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
			stop() // Trigger shutdown if HTTP server fails to start
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case <-shutdownCtx.Done():
		log.Println("Shutdown signal received")
	case err := <-serverErr:
		log.Printf("HTTP server failed: %v", err)
	}

	// Perform graceful shutdown with timeout
	gracefulShutdown(container, httpSvc)
}

func gracefulShutdown(container actions.Container, httpSvc *services.HTTPService) {
	log.Println("Starting graceful shutdown...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server first (stop accepting new requests)
	log.Println("Shutting down HTTP server...")
	if err := httpSvc.Shutdown(shutdownCtx); err != nil {
		if err == context.DeadlineExceeded {
			log.Println("HTTP server shutdown timeout - forcing close")
		} else {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	} else {
		log.Println("HTTP server stopped gracefully")
	}

	// Container close will handle database and other services
	log.Println("Closing container services...")
	// if err := container.Close(); err != nil {
	// 	log.Printf("Container close error: %v", err)
	// } else {
	// 	log.Println("All services closed successfully")
	// }

	log.Println("Service shutdown complete")
}

// Alternative simpler version if you prefer a more concise approach:
func mainSimple() {
	resourcesDir := "./resources"
	config, err := actions.LoadServiceConfig(resourcesDir + "/service.yml")
	if err != nil {
		log.Fatal("Failed to load service config:", err)
	}

	// Initialize container
	container, err := actions.NewContainer(config)
	if err != nil {
		log.Fatal("Failed to initialize container:", err)
	}
	// defer container.Close()

	// Setup services
	hub := services.NewPubSubHub(container)
	httpSvc := services.NewHTTPService(container, hub)

	// Register modules
	registerModules(container)

	// Run server with graceful shutdown
	runServerWithShutdown(container, httpSvc)
}

func registerModules(container actions.Container) {
	for _, module := range container.Config().Modules {
		switch module {
		case "auth":
			auth.Register(container.Registry())
		case "story":
			story.Register(container.Registry())
		}
	}
}

func runServerWithShutdown(container actions.Container, httpSvc *services.HTTPService) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start server
	go func() {
		log.Printf("Server starting on port :%d", container.Config().HTTP.Port)
		if err := httpSvc.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpSvc.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}
}
