package apps

import (
	"e-shop-api/internal/constants"
	"e-shop-api/internal/middlewares"
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
	app.Use(middlewares.RequestID()) 		   //# Request ID - must be first for logging
	app.Use(middlewares.LoggerMiddleware())    //# Logging HTTP request
	app.Use(gin.Recovery()) 				   //# Recover from panic
	app.Use(middlewares.InitCORS()) 		   //# CORS middleware
	app.Use(middlewares.ResponseMiddleware())  //# Response middleware

	// Selective Middlewares
	return &MiddlewareRegistry{
        Auth:      middlewares.AuthMiddleware(),
		Admin:     middlewares.RoleMiddleware(constants.Admin),
		Seller:    middlewares.RoleMiddleware(constants.Seller, constants.Admin),
		Buyer:     middlewares.RoleMiddleware(constants.Buyer),
		RequestID: middlewares.RequestID(),

		// Register new middlewares here
    }
}