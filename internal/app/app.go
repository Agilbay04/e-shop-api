package app

import (
	"e-shop-api/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	// Register json tag name
	util.RegisterJSONTagName()

	// Middleware
	middlewareRegistry := NewMiddlewareRegistry(r)

	// Register repository
	repoRegistry := NewRepositoryRegistry(db)

	// Register service
	svcRegistry := NewServiceRegistry(repoRegistry, db, rdb)
	
	// Register handler
	handlerRegistry := NewHandlerRegistry(svcRegistry)

	// Register routes
	RegisterRoutes(r, handlerRegistry, middlewareRegistry, rdb)

	return r
}