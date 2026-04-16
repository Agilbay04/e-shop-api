package app

import (
	"e-shop-api/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.ResponseMiddleware())

	// Register repository
	repoRegistry := NewRepositoryRegistry(db)

	// Register service
	svcRegistry := NewServiceRegistry(repoRegistry)
	
	// Register handler
	handlerRegistry := NewHandlerRegistry(svcRegistry)

	// Register routes
	RegisterRoutes(r, handlerRegistry)

	return r
}