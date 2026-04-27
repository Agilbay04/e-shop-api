package apps

import (
	"e-shop-api/internal/handlers"

	"gorm.io/gorm"
)

type HandlerRegistry struct {
	AuthHandler    *handlers.AuthHandler
	UserHandler    *handlers.UserHandler
	StoreHandler   *handlers.StoreHandler
	ProductHandler *handlers.ProductHandler
	OrderHandler   *handlers.OrderHandler
	HealthHandler  *handlers.HealthHandler
}

func NewHandlerRegistry(svc *ServiceRegistry, db *gorm.DB, client *ClientRegistry) *HandlerRegistry {
	return &HandlerRegistry{
		AuthHandler:    handlers.NewAuthHandler(svc.AuthService),
		UserHandler:    handlers.NewUserHandler(svc.UserService),
		StoreHandler:   handlers.NewStoreHandler(svc.StoreService),
		ProductHandler: handlers.NewProductHandler(svc.ProductService),
		OrderHandler:   handlers.NewOrderHandler(svc.OrderService),
		HealthHandler:  handlers.NewHealthHandler(db, client.Redis),

		// Register new handlers here
	}
}