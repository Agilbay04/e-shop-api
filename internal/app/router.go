package app

import (
	"e-shop-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *HandlerRegistry) {
	api := r.Group("/api/v1")
	{
		/**
		PUBLIC ROUTES
		*/
		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.AuthHandler.Register)
			auth.POST("/login", h.AuthHandler.Login)
		}

		/**
		PROTECTED ROUTES
		*/
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Store Routes
			protected.POST("/store", h.StoreHandler.CreateStore)

			// Product Routes
			protected.GET("/product", h.ProductHandler.Index)
			protected.POST("/product", h.ProductHandler.CreateProduct)

			// Order Routes
			protected.POST("/order", h.OrderHandler.CreateOrder)
		}
	}
}