package app

import (
	"e-shop-api/internal/config"
	"e-shop-api/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	rdb := config.ConnectRedis()

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
	RegisterRoutes(r, handlerRegistry, middlewareRegistry)

	return r
}