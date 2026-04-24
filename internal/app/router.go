package app

import (
	"e-shop-api/internal/middleware"
	"e-shop-api/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(r *gin.Engine, h *HandlerRegistry, m *MiddlewareRegistry, rdb *redis.Client) {
	// Health routes (no auth required)
	r.GET("/health", h.HealthHandler.Health)
	r.GET("/ready", h.HealthHandler.Readiness)

	// Request ID middleware
	r.Use(m.RequestID)

	api := r.Group("/api/v1")
	{
		// Auth Routes (Public)
		registerAuthRoutes(api, h, m, rdb)

		// Protected Routes (Common Auth)
		protected := api.Group("/")
		protected.Use(m.Auth)
		{
			registerStoreRoutes(protected, h, m)
			registerProductRoutes(protected, h, m)
			registerOrderRoutes(protected, h, m)
		}
	}
}

// Sub Routes
func registerAuthRoutes(api *gin.RouterGroup, h *HandlerRegistry, m *MiddlewareRegistry, rdb *redis.Client) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.AuthHandler.Register)
		
		auth.POST("/login", 
			middleware.RateLimiter(rdb, "login", util.TimeParse("5s")),
			h.AuthHandler.Login,
		)
		
		auth.POST("/forgot-password", 
			middleware.RateLimiter(rdb, "forgot-password", util.TimeParse("1m")), 
			h.AuthHandler.ForgotPassword,
		)

		auth.PUT("/reset-password", h.AuthHandler.ResetPassword)

		// Protected Routes (Common Auth)
		protected := auth.Group("/")
		protected.Use(m.Auth)
		{
			protected.GET("/profile", h.AuthHandler.Profile)
			protected.POST("/upload-picture", h.AuthHandler.UploadPicture)
		}
	}
}

func registerStoreRoutes(rg *gin.RouterGroup, h *HandlerRegistry, m *MiddlewareRegistry) {
	stores := rg.Group("/stores")
	{
		// Publicly available for all authenticated users
		stores.GET("/", h.StoreHandler.GetStores)

		// Actions requiring Seller or Admin privileges
		privileged := stores.Group("/")

		// Admin included in m.Seller registry
		privileged.Use(m.Seller)
		{
			privileged.POST("/", h.StoreHandler.CreateStore)
			privileged.PUT("/:id", h.StoreHandler.UpdateStore)
			privileged.PATCH("/:id", h.StoreHandler.DeleteStore)
			privileged.PATCH("/activate", h.StoreHandler.ActivateStore)
		}
	}
}

func registerProductRoutes(rg *gin.RouterGroup, h *HandlerRegistry, m *MiddlewareRegistry) {
	products := rg.Group("/products")
	{
		products.GET("/", h.ProductHandler.Index)

		privileged := products.Group("/")
		privileged.Use(m.Seller)
		{
			privileged.POST("/", h.ProductHandler.CreateProduct)
			privileged.PUT("/:id", h.ProductHandler.UpdateProduct)
			privileged.PATCH("/:id", h.ProductHandler.DeleteProduct)
			privileged.PATCH("/activate", h.ProductHandler.ActivateProduct)
		}
	}
}

func registerOrderRoutes(rg *gin.RouterGroup, h *HandlerRegistry, m *MiddlewareRegistry) {
	orders := rg.Group("/orders")
	{
		orders.GET("/", h.OrderHandler.GetOrders)

		buyerOnly := orders.Group("/")
		buyerOnly.Use(m.Buyer)
		{
			buyerOnly.POST("/", h.OrderHandler.CreateOrder)
			buyerOnly.PUT("/:id", h.OrderHandler.UpdateOrder)
			buyerOnly.PATCH("/:id/cancel", h.OrderHandler.CancelOrder)
			buyerOnly.PATCH("/:id/confirm", h.OrderHandler.ConfirmOrder)
		}
	}
}
