package app

import (
	"e-shop-api/internal/middleware"
	"e-shop-api/internal/model"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type MiddlewareRegistry struct {
    Auth      gin.HandlerFunc
	Admin     gin.HandlerFunc
	Seller    gin.HandlerFunc
	Buyer     gin.HandlerFunc
	RequestID gin.HandlerFunc
}

func NewMiddlewareRegistry(app *gin.Engine) *MiddlewareRegistry {
	// Global Setup Trusted Proxies
	proxies := os.Getenv("TRUSTED_PROXIES")
    if proxies != "" {
        app.SetTrustedProxies(strings.Split(proxies, ","))
    } else {
        app.SetTrustedProxies(nil)
    }

	// Global Middleware
	// Note: please don't change the order
	app.Use(middleware.RequestID()) // Request ID - must be first for logging
	app.Use(middleware.LoggerMiddleware()) // Logging HTTP request
	app.Use(gin.Recovery()) // Recover from panic
	app.Use(middleware.InitCORS()) // CORS middleware
	app.Use(middleware.ResponseMiddleware()) // Response middleware

	// Selective Middlewares
	return &MiddlewareRegistry{
        Auth:      middleware.AuthMiddleware(),
		Admin:     middleware.RoleMiddleware(model.Admin),
		Seller:    middleware.RoleMiddleware(model.Seller, model.Admin),
		Buyer:     middleware.RoleMiddleware(model.Buyer),
		RequestID: middleware.RequestID(),

		// Register new middlewares here
    }
}