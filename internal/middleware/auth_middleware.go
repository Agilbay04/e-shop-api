package middleware

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/pkg/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Error(util.UnauthorizedException("Missing authorization header"))
			ctx.Abort()
			return
		}

		// Format Authorization header must be "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.Error(util.UnauthorizedException("Invalid authorization format"))
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		// Parse token and validate
		claims, err := util.ParseToken(tokenString)
		if err != nil {
			ctx.Error(util.UnauthorizedException("Token is invalid or expired"))
			ctx.Abort()
			return
		}

		// Map claims to struct
		currentUser := dto.CurrentUser {
			ID:       claims.ID,
			Username: claims.Username,
			Email:    claims.Email,
			Role:     claims.Role,
		}
		
		// Set current user to context
		ctx.Set("currentUser", currentUser)
		ctx.Next()
	}
}