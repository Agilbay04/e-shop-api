package util

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// GracefulShutdown handles graceful shutdown of the server
func GracefulShutdown(srv *http.Server, db *gorm.DB, rdb *redis.Client, timeout time.Duration) {
	// Channel for receiving signal os
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	log.Println("Shutting down server...")

	// Create new context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Close database connection
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("Database connection closed.")
		}
	}

	// Close Redis connection
	if rdb != nil {
		rdb.Close()
		log.Println("Redis connection closed.")
	}

	log.Println("Server exited gracefully")
}