package app

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *HandlerRegistry, m *MiddlewareRegistry) {
	api := r.Group("/api/v1")
	{
		/** PUBLIC ROUTES */
		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.AuthHandler.Register)
			auth.POST("/login", h.AuthHandler.Login)
		}

		/** PROTECTED ROUTES */
		protected := api.Group("/")
		protected.Use(m.Auth)
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
			protected.GET("/orders", h.OrderHandler.GetOrders)
			protected.POST("/orders", h.OrderHandler.CreateOrder)
			protected.PUT("/orders/:id", h.OrderHandler.UpdateOrder)
			protected.PATCH("/orders/:id/cancel", h.OrderHandler.CancelOrder)
			protected.PATCH("/orders/:id/confirm", h.OrderHandler.ConfirmOrder)
		}
	}
}
