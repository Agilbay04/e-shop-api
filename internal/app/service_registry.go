package app

import (
	"e-shop-api/internal/service"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceRegistry struct {
	AuthService    service.AuthService
	StoreService   service.StoreService
	ProductService service.ProductService
	OrderService   service.OrderService
	notifService   service.NotificationService
}

func NewServiceRegistry(repo *RepositoryRegistry, db *gorm.DB, rdb *redis.Client) *ServiceRegistry {
	return &ServiceRegistry{
		AuthService:    service.NewAuthService(db, repo.UserRepo, repo.UserQuery, rdb),
		StoreService:   service.NewStoreService(db, repo.StoreRepo, repo.StoreQuery, repo.OrderQuery, repo.UserQuery),
		ProductService: service.NewProductService(db, repo.ProductRepo, repo.ProductQuery, repo.StoreQuery),
		OrderService:   service.NewOrderService(
			db, 
			repo.OrderRepo, 
			repo.OrderQuery, 
			repo.ProductRepo, 
			repo.ProductQuery, 
			repo.StoreQuery, 
			service.NewNotificationService(),
		),
	}
}
