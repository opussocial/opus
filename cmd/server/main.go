package main

import (
    // "fmt"
    "log"
    "time"
    "context"
    "net/http"
    "syscall"
    "os/signal"

    "gitlab.com/pedrokoblitz/opus-go/modules/auth"
    "gitlab.com/pedrokoblitz/opus-go/modules/story"
    "gitlab.com/pedrokoblitz/opus-go/services"
)

func main() {

    resourcesDir := "./resources"
    config, err := services.LoadServiceConfig(resourcesDir + "/service.yml")
    if err != nil {
        log.Fatal(err)
    }

    // Initialize services
    container := services.NewContainer(config)
    hub := services.NewPubSubHub(container)
    httpSvc := services.NewHTTPService(container, hub)

    // Register modules
    for _, module := range config.Service.Modules {
        switch module {
        case "auth":
            auth.Register(container.Registry)
        case "story":
            story.Register(container.Registry)
        }
    }

    // Create context for shutdown
    shutdownCtx, stop := signal.NotifyContext(context.Background(), 
        syscall.SIGINT, syscall.SIGTERM)
    defer stop()

    go func() {
        log.Printf("Starting HTTP server on :%d", config.Service.HTTP.Port)
        if err := httpSvc.Start(); err != nil && err != http.ErrServerClosed {
            log.Printf("HTTP service failed: %v", err)
            stop() // Trigger shutdown if HTTP fails
        }
    }()

    // Wait for shutdown signal
    <-shutdownCtx.Done()
    log.Println("Shutdown signal received")

    // Start graceful shutdown with timeout
    gracefulCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Shutdown HTTP server
    log.Println("Shutting down HTTP server...")
    if err := httpSvc.Shutdown(gracefulCtx); err != nil {
        log.Printf("HTTP server shutdown error: %v", err)
    } else {
        log.Println("HTTP server stopped gracefully")
    }

    // Close database connection
    log.Println("Closing database connection...")
    db := container.DB
    if err := db.Close(); err != nil {
        log.Printf("Database close error: %v", err)
    } else {
        log.Println("Database connection closed")
    }

    log.Println("Service shutdown complete")
}
