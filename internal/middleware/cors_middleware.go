package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitCORS() gin.HandlerFunc {
	// Get allowed origins
	originsConfig := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowedOrigins []string

	if originsConfig == "" {
		// Default allowed origins
		allowedOrigins = []string{"http://localhost:8001", "https://localhost:8001", "http://localhost:5500"}
	} else {
		allowedOrigins = strings.Split(originsConfig, ",")
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}