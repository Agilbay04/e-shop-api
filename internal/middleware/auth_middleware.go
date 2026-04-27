package middleware

import (
	"e-shop-api/internal/constant"
	"e-shop-api/internal/dto"
	"e-shop-api/internal/pkg/util"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get Authorization header
		authHeader := ctx.GetHeader(constant.AuthorizationHeader)
		if authHeader == "" {
			ctx.Error(util.UnauthorizedException("Missing authorization header"))
			ctx.Abort()
			return
		}

		// Format Authorization header must be "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != constant.BearerPrefix {
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
			Picture:  claims.Picture,
		}
		
		// Set current user to context
		ctx.Set("currentUser", currentUser)
		ctx.Next()
	}
}

func RoleMiddleware(roles ...constant.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
        // Get user from context
        val, exists := ctx.Get("currentUser")
        if !exists {
            ctx.Error(util.UnauthorizedException("Unauthorized"))
			ctx.Abort()
            return
        }

        user := val.(dto.CurrentUser)
        
        // Check if user role is allowed
        isAllowed := slices.Contains(roles, user.Role)

        if !isAllowed {
			ctx.Error(util.ForbiddenException("Access is denied for your role"))
			ctx.Abort()
            return
        }

        ctx.Next()
    }
}
