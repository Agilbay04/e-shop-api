package main

import (
	"e-shop-api/internal/app"
	"e-shop-api/internal/config"
	"e-shop-api/internal/pkg/logger"
	"e-shop-api/internal/pkg/util"
	"errors"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		logger.L.Info("Warning: .env file not found, using system environment variables")
	}

	// Init logger
	logger.InitLogger()
	defer logger.L.Sync()

	logger.L.Info("Starting server...")

	// Connect database
	db := config.ConnectDatabase()

	// Connect redis
	rdb := config.ConnectRedis()

	// Setup router
	r := app.SetupRouter(db, rdb)

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
		useHTTPS := util.GetEnvBool("USE_HTTPS", "false")
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
	util.GracefulShutdown(srv, db, rdb, util.TimeParse("5s"))
}