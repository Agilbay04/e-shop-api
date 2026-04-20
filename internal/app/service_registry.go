package app

import (
	"e-shop-api/internal/service"

	"gorm.io/gorm"
)

type ServiceRegistry struct {
	AuthService    service.AuthService
	StoreService   service.StoreService
	ProductService service.ProductService
	OrderService   service.OrderService
}

func NewServiceRegistry(repo *RepositoryRegistry, db *gorm.DB) *ServiceRegistry {
	return &ServiceRegistry{
		AuthService:    service.NewAuthService(repo.UserRepo, repo.UserQuery),
		StoreService:   service.NewStoreService(repo.StoreRepo, repo.StoreQuery, repo.UserQuery),
		ProductService: service.NewProductService(repo.ProductRepo, repo.ProductQuery, repo.StoreQuery),
		OrderService:   service.NewOrderService(db, repo.OrderRepo, repo.ProductRepo, repo.ProductQuery),
	}
}