package app

import (
	"e-shop-api/internal/service"

	"gorm.io/gorm"
)

type ServiceRegistry struct {
	AuthService    service.AuthService
	UserService    service.UserService
	StoreService   service.StoreService
	ProductService service.ProductService
	OrderService   service.OrderService
	NotifService   service.NotificationService
}

func NewServiceRegistry(db *gorm.DB, repo *RepositoryRegistry, client *ClientRegistry) *ServiceRegistry {
	notifService := service.NewNotificationService()

	return &ServiceRegistry{
		AuthService:    service.NewAuthService(db, repo.UserRepo, repo.UserQuery, notifService, client.Redis),
		UserService:    service.NewUserService(db, repo.UserRepo, repo.UserQuery, client.Redis),
		StoreService:   service.NewStoreService(db, repo.StoreRepo, repo.StoreQuery, repo.OrderQuery, repo.UserQuery),
		ProductService: service.NewProductService(db, repo.ProductRepo, repo.ProductQuery, repo.StoreQuery),
		OrderService: service.NewOrderService(
			db,
			repo.OrderRepo,
			repo.OrderQuery,
			repo.ProductRepo,
			repo.ProductQuery,
			repo.StoreQuery,
			notifService,
		),
		NotifService: notifService,

		// Register new services here
	}
}
