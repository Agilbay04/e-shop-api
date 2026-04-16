package app

import (
	"e-shop-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewMiddlewareRegistry(app *gin.Engine) {
	app.Use(middleware.ResponseMiddleware())
}