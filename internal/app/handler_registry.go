package app

import (
	"e-shop-api/internal/handler"
)

type HandlerRegistry struct {
	AuthHandler  	*handler.AuthHandler
	StoreHandler 	*handler.StoreHandler
	ProductHandler 	*handler.ProductHandler
}

func NewHandlerRegistry(svc *ServiceRegistry) *HandlerRegistry {
	return &HandlerRegistry{
		AuthHandler:  	handler.NewAuthHandler(svc.AuthService),
		StoreHandler: 	handler.NewStoreHandler(svc.StoreService),
		ProductHandler: handler.NewProductHandler(svc.ProductService),
	}
}