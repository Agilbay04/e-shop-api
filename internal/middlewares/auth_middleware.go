package middlewares

import (
	"e-shop-api/internal/constants"
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/pkg/utils"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get Authorization header
		authHeader := ctx.GetHeader(constants.AuthorizationHeader)
		if authHeader == "" {
			ctx.Error(utils.UnauthorizedException("Missing authorization header"))
			ctx.Abort()
			return
		}

		// Format Authorization header must be "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != constants.BearerPrefix {
			ctx.Error(utils.UnauthorizedException("Invalid authorization format"))
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		// Parse token and validate
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			ctx.Error(utils.UnauthorizedException("Token is invalid or expired"))
			ctx.Abort()
			return
		}

		// Map claims to struct
		currentUser := dtos.CurrentUser {
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

func RoleMiddleware(roles ...constants.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
        // Get user from context
        val, exists := ctx.Get("currentUser")
        if !exists {
            ctx.Error(utils.UnauthorizedException("Unauthorized"))
			ctx.Abort()
            return
        }

        user := val.(dtos.CurrentUser)
        
        // Check if user role is allowed
        isAllowed := slices.Contains(roles, user.Role)

        if !isAllowed {
			ctx.Error(utils.ForbiddenException("Access is denied for your role"))
			ctx.Abort()
            return
        }

        ctx.Next()
    }
}
