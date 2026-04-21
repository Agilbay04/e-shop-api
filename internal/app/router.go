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
			protected.GET("/stores", h.StoreHandler.GetStores)
			protected.POST("/stores", h.StoreHandler.CreateStore)
			protected.PUT("/stores/:id", h.StoreHandler.UpdateStore)
			protected.PATCH("/stores/:id", h.StoreHandler.DeleteStore)
			protected.PATCH("/stores/activate", h.StoreHandler.ActivateStore)

			// Product Routes
			protected.GET("/products", h.ProductHandler.Index)
			protected.POST("/products", h.ProductHandler.CreateProduct)
			protected.PUT("/products/:id", h.ProductHandler.UpdateProduct)
			protected.PATCH("/products/:id", h.ProductHandler.DeleteProduct)
			protected.PATCH("/products/activate", h.ProductHandler.ActivateProduct)

			// Order Routes
			protected.POST("/orders", h.OrderHandler.CreateOrder)
		}
	}
}
