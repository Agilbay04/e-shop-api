package app

import (
	"e-shop-api/internal/middleware"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type MiddlewareRegistry struct {
    Auth gin.HandlerFunc
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
	app.Use(middleware.InitCORS())
	app.Use(middleware.ResponseMiddleware())

	// Selective Middlewares
	return &MiddlewareRegistry{
        Auth: middleware.AuthMiddleware(),
    }
}