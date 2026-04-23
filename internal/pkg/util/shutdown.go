package util

import (
	"context"
	"e-shop-api/internal/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GracefulShutdown handles graceful shutdown of the server
func GracefulShutdown(srv *http.Server, db *gorm.DB, rdb *redis.Client, timeout time.Duration) {
	// Channel for receiving signal os
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	logger.L.Info("Shutting down server...")

	// Create new context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		logger.L.Fatal("Server forced to shutdown: %v", zap.Error(err))
	}

	// Close database connection
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
			logger.L.Info("Database connection closed.")
		}
	}

	// Close Redis connection
	if rdb != nil {
		rdb.Close()
		logger.L.Info("Redis connection closed.")
	}

	logger.L.Info("Server exited gracefully")
}