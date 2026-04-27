package apps

import (
	"e-shop-api/internal/services"

	"gorm.io/gorm"
)

type ServiceRegistry struct {
	AuthService    services.AuthService
	UserService    services.UserService
	StoreService   services.StoreService
	ProductService services.ProductService
	OrderService   services.OrderService
	NotifService   services.NotificationService
}

func NewServiceRegistry(db *gorm.DB, repo *RepositoryRegistry, client *ClientRegistry) *ServiceRegistry {
	notifService := services.NewNotificationService()

	return &ServiceRegistry{
		AuthService:    services.NewAuthService(db, repo.UserRepo, repo.UserQuery, notifService, client.Redis),
		UserService:    services.NewUserService(db, repo.UserRepo, repo.UserQuery, client.Redis),
		StoreService:   services.NewStoreService(db, repo.StoreRepo, repo.StoreQuery, repo.OrderQuery, repo.UserQuery),
		ProductService: services.NewProductService(db, repo.ProductRepo, repo.ProductQuery, repo.StoreQuery),
		OrderService: 	services.NewOrderService(
			db,
			repo.OrderRepo,
			repo.OrderQuery,
			repo.ProductRepo,
			repo.ProductQuery,
			repo.StoreQuery,
			repo.OrderSequenceRepo,
			notifService,
		),
		NotifService: 	notifService,

		// Register new services here
	}
}
