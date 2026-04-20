package app

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Validator tag json
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        v.RegisterTagNameFunc(func(fld reflect.StructField) string {
            name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
            if name == "-" {
                return ""
            }
            return name
        })
    }

	// Middleware
	NewMiddlewareRegistry(r)

	// Register repository
	repoRegistry := NewRepositoryRegistry(db)

	// Register service
	svcRegistry := NewServiceRegistry(repoRegistry, db)
	
	// Register handler
	handlerRegistry := NewHandlerRegistry(svcRegistry)

	// Register routes
	RegisterRoutes(r, handlerRegistry)

	return r
}