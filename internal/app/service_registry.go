package app

import (
	"e-shop-api/internal/service"
)

type ServiceRegistry struct {
	AuthService    service.AuthService
	StoreService   service.StoreService
	ProductService service.ProductService
}

func NewServiceRegistry(repo *RepositoryRegistry) *ServiceRegistry {
	return &ServiceRegistry{
		AuthService:    service.NewAuthService(repo.UserRepo, repo.UserQuery),
		StoreService:   service.NewStoreService(repo.StoreRepo, repo.StoreQuery, repo.UserQuery),
		ProductService: service.NewProductService(repo.ProductRepo, repo.ProductQuery, repo.StoreQuery),
	}
}