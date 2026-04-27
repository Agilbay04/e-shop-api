package apps

import (
	"e-shop-api/internal/repositories"
	"gorm.io/gorm"
)

type RepositoryRegistry struct {
	UserRepo         	repositories.UserRepository
	UserQuery        	repositories.UserQueryRepository
	StoreRepo        	repositories.StoreRepository
	StoreQuery       	repositories.StoreQueryRepository
	ProductRepo      	repositories.ProductRepository
	ProductQuery     	repositories.ProductQueryRepository
	OrderRepo         	repositories.OrderRepository
	OrderQuery        	repositories.OrderQueryRepository
	OrderSequenceRepo 	repositories.OrderSequenceRepository
}

func NewRepositoryRegistry(db *gorm.DB) *RepositoryRegistry {
	return &RepositoryRegistry{
		UserRepo:         	repositories.NewUserRepository(db),
		UserQuery:        	repositories.NewUserQueryRepository(db),
		StoreRepo:        	repositories.NewStoreRepository(db),
		StoreQuery:       	repositories.NewStoreQueryRepository(db),
		ProductRepo:      	repositories.NewProductRepository(db),
		ProductQuery:     	repositories.NewProductQueryRepository(db),
		OrderRepo:        	repositories.NewOrderRepository(db),
		OrderQuery:       	repositories.NewOrderQueryRepository(db),
		OrderSequenceRepo: 	repositories.NewOrderSequenceRepository(db),

		// Register new repositories here
	}
}
