package main

import (
	"e-shop-api/internal/apps"
	"e-shop-api/internal/configs"
	"e-shop-api/internal/pkg/logger"
	"e-shop-api/internal/pkg/utils"
	"errors"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load file .env first (ensures APP_ENV is available for logger config)
	err := godotenv.Load()

	// Init logger before any logging calls
	logger.InitLogger()
	defer logger.L.Sync()

	// Log .env warning now that logger is initialized
	if err != nil {
		logger.L.Info("Warning: .env file not found, using system environment variables")
	}

	logger.L.Info("Starting server...")

	// Connect database
	db := configs.ConnectDatabase()

	// Connect redis
	rdb := configs.ConnectRedis()

	// Init setup for all dependencies
	r := apps.Setup(db, rdb)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8001"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Run server in goroutine
	go func() {
		useHTTPS := utils.GetEnvBool("USE_HTTPS", "false")
		cert := os.Getenv("SSL_CERT_PATH")
		key := os.Getenv("SSL_KEY_PATH")

		if useHTTPS && cert != "" && key != "" {
			logger.L.Info("Server e-shop-api starting on https://localhost:"+port, zap.String("port", port))
			if err := srv.ListenAndServeTLS(cert, key); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.L.Fatal("HTTPS Listen error", zap.Error(err))
			}
		} else {
			logger.L.Info("Server e-shop-api starting on http://localhost:"+port, zap.String("port", port))
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.L.Fatal("HTTP Listen error", zap.Error(err))
			}
		}
	}()

	// Graceful shutdown
	utils.GracefulShutdown(srv, db, rdb, utils.TimeParse("5s"))
}