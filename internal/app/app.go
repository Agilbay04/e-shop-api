package app

import (
	"e-shop-api/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Register json tag name
	util.RegisterJSONTagName()

	// Middleware
	middlewareRegistry := NewMiddlewareRegistry(r)

	// Register repository
	repoRegistry := NewRepositoryRegistry(db)

	// Register service
	svcRegistry := NewServiceRegistry(repoRegistry, db)
	
	// Register handler
	handlerRegistry := NewHandlerRegistry(svcRegistry)

	// Register routes
	RegisterRoutes(r, handlerRegistry, middlewareRegistry)

	return r
}