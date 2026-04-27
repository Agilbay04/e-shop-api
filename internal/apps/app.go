package apps

import (
	"e-shop-api/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.New()

	// Register JSON tag name
	utils.RegisterJSONTagName()

	// Init middlewares
	middlewareRegistry := NewMiddlewareRegistry(r)

	// Init clients
	clientRegistry := NewClientRegistry(rdb)
	
	// Init repositories
	repoRegistry := NewRepositoryRegistry(db)
	
	// Init services
	svcRegistry := NewServiceRegistry(db, repoRegistry, clientRegistry)
	
	// Init handlers
	handlerRegistry := NewHandlerRegistry(svcRegistry, db, clientRegistry)

	// Register routes
	RegisterRoutes(r, handlerRegistry, middlewareRegistry, clientRegistry.Redis)

	return r
}