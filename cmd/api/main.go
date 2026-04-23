package main

import (
	"e-shop-api/internal/app"
	"e-shop-api/internal/config"
	"e-shop-api/internal/pkg/util"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Connect database
	db := config.ConnectDatabase()

	// Connect redis
	rdb := config.ConnectRedis()

	// Setup router
	r := app.SetupRouter(db, rdb)
	
	port := os.Getenv("SERVER_PORT")
    if port == "" { port = "8001" }

    srv := &http.Server{
        Addr:    ":" + port,
        Handler: r,
    }

    // Run server in goroutine
    go func() {
        log.Printf("Server starting on port %s", port)
        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("Listen error: %v", err)
        }
    }()

    // Graceful shutdown
    util.GracefulShutdown(srv, db, rdb, 5*time.Second)
}