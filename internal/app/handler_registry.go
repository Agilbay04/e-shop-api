package app

import (
	"e-shop-api/internal/handler"

	"gorm.io/gorm"
)

type HandlerRegistry struct {
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	StoreHandler   *handler.StoreHandler
	ProductHandler *handler.ProductHandler
	OrderHandler   *handler.OrderHandler
	HealthHandler  *handler.HealthHandler
}

func NewHandlerRegistry(svc *ServiceRegistry, db *gorm.DB, client *ClientRegistry) *HandlerRegistry {
	return &HandlerRegistry{
		AuthHandler:    handler.NewAuthHandler(svc.AuthService),
		UserHandler:    handler.NewUserHandler(svc.UserService),
		StoreHandler:   handler.NewStoreHandler(svc.StoreService),
		ProductHandler: handler.NewProductHandler(svc.ProductService),
		OrderHandler:   handler.NewOrderHandler(svc.OrderService),
		HealthHandler:  handler.NewHealthHandler(db, client.Redis),

		// Register new handlers here
	}
}