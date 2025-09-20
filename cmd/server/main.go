package main

import (
    // "os"
    // "fmt"
    "log"
    "time"
    "context"
    "net/http"
    "syscall"
    "os/signal"

    "github.com/joho/godotenv"

    "gitlab.com/pedrokoblitz/opus-go/providers/auth"
    "gitlab.com/pedrokoblitz/opus-go/services"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    debug := os.Getenv("DEBUG")
    environment := os.Getenv("ENVIRONMENT")
    port := os.Getenv("APP_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")

    resourcesDir := "./resources"
    config, err := services.LoadServiceConfig(resourcesDir + "/service.yml")
    if err != nil {
        log.Fatal(err)
    }
    config.Service.Debug = debug
    config.Service.HTTP.Port = port
    config.Service.DSN = dbURL

    // Initialize services
    container := services.NewContainer(config)
    hub := services.NewPubSubHub(container)
    httpSvc := services.NewHTTPService(container, hub)

    // Register providers
    for _, module := range config.Service.Modules {
        switch module {
        case "auth":
            auth.Register(container.Registry)
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
}









package main

import (
    "fmt"
    "log"
    "strconv"
    
    "github.com/joho/godotenv"
)

// type Config struct {
//     Port        int
//     DatabaseURL string
//     Debug       bool
//     APIKey      string
//     Environment string
// }

// func LoadConfig() (*Config, error) {
//     err := godotenv.Load()
//     if err != nil {
//         log.Println("No .env file found, using system environment variables")
//     }

//     port, _ := strconv.Atoi(getEnv("APP_PORT", "8080"))
//     debug, _ := strconv.ParseBool(getEnv("DEBUG", "false"))

//     return &Config{
//         Port:        port,
//         DatabaseURL: getEnv("DATABASE_URL", ""),
//         Debug:       debug,
//         APIKey:      getEnv("API_KEY", ""),
//         Environment: getEnv("ENVIRONMENT", "development"),
//     }, nil
// }

// func getEnv(key, defaultValue string) string {
//     value := os.Getenv(key)
//     if value == "" {
//         return defaultValue
//     }
//     return value
// }

// func main() {
//     config, err := LoadConfig()
//     if err != nil {
//         log.Fatal(err)
//     }

//     fmt.Printf("Config: %+v\n", config)
// }

// func loadEnv() error {
//     env := os.Getenv("GO_ENV")
//     if env == "" {
//         env = "development"
//     }
    
//     // Try to load environment-specific file first
//     godotenv.Load(fmt.Sprintf(".env.%s", env))
    
//     // Fall back to default .env file
//     return godotenv.Load()
// }
