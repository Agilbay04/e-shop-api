package app

import (
	"e-shop-api/internal/repository"
	"gorm.io/gorm"
)

type RepositoryRegistry struct {
	UserRepo         repository.UserRepository
	UserQuery        repository.UserQueryRepository
	StoreRepo        repository.StoreRepository
	StoreQuery       repository.StoreQueryRepository
	ProductRepo      repository.ProductRepository
	ProductQuery     repository.ProductQueryRepository
	OrderRepo         repository.OrderRepository
	OrderQuery        repository.OrderQueryRepository
	OrderSequenceRepo repository.OrderSequenceRepository
}

func NewRepositoryRegistry(db *gorm.DB) *RepositoryRegistry {
	return &RepositoryRegistry{
		UserRepo:         repository.NewUserRepository(db),
		UserQuery:        repository.NewUserQueryRepository(db),
		StoreRepo:        repository.NewStoreRepository(db),
		StoreQuery:       repository.NewStoreQueryRepository(db),
		ProductRepo:      repository.NewProductRepository(db),
		ProductQuery:     repository.NewProductQueryRepository(db),
		OrderRepo:        repository.NewOrderRepository(db),
		OrderQuery:       repository.NewOrderQueryRepository(db),
		OrderSequenceRepo: repository.NewOrderSequenceRepository(db),

		// Register new repositories here
	}
}
