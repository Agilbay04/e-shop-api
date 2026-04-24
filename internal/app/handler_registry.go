package app

import (
	"e-shop-api/internal/handler"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HandlerRegistry struct {
	AuthHandler   	*handler.AuthHandler
	StoreHandler 	*handler.StoreHandler
	ProductHandler 	*handler.ProductHandler
	OrderHandler  	*handler.OrderHandler
	HealthHandler 	*handler.HealthHandler
}

func NewHandlerRegistry(svc *ServiceRegistry, db *gorm.DB, rdb *redis.Client) *HandlerRegistry {
	return &HandlerRegistry{
		AuthHandler:  	handler.NewAuthHandler(svc.AuthService),
		StoreHandler: 	handler.NewStoreHandler(svc.StoreService),
		ProductHandler: handler.NewProductHandler(svc.ProductService),
		OrderHandler: 	handler.NewOrderHandler(svc.OrderService),
		HealthHandler: 	handler.NewHealthHandler(db, rdb),
	}
}