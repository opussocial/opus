package main

import (
    "context"
    "log"
    "net/http"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/opussocialcontent/opus-go/actions"
    "github.com/opussocialcontent/opus-go/services"
)

func main() {
    resourcesDir := "./resources"
    
    // Define default service configuration
    defaultServiceConfig := &actions.ServiceConfig{
        HTTP: actions.HTTPConfig{
            Host: "localhost",
            Port: 8080,
        },
        DB: actions.DBConfig{
            Driver:   "mysql",
            Host:     "localhost",
            Port:     3306,
            User:     "dbuser",
            Password: "1234",
            Database: "opus",
        },
    }
    
    // Load service config with optional default
    config, err := actions.LoadServiceConfig(resourcesDir+"/service.yml", defaultServiceConfig)
    if err != nil {
        log.Fatal("Failed to load service config:", err)
    }

    // Initialize container with DB and theme services
    container, err := actions.NewContainer(config)
    if err != nil {
        log.Fatal("Failed to initialize container:", err)
    }
    defer func() {
        if err := container.Close(); err != nil {
            log.Printf("Error closing container: %v", err)
        }
    }()

    // Initialize HTTP service without PubSub hub
    httpSvc := services.NewHTTPService(container, nil)

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

    // Container close will handle database
    log.Println("Closing container services...")
    if err := container.Close(); err != nil {
        log.Printf("Container close error: %v", err)
    } else {
        log.Println("All services closed successfully")
    }

    log.Println("Service shutdown complete")
}
