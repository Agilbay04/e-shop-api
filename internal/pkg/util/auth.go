package util

import (
	"e-shop-api/internal/dto"
	"github.com/gin-gonic/gin"
)

const UserContextKey = "currentUser"

// GetCurrentUser to extract CurrentUser from gin.Context
func GetCurrentUser(ctx *gin.Context) dto.CurrentUser {
	val, exists := ctx.Get(UserContextKey)
	if !exists {
		return dto.CurrentUser{}
	}

	// type assertion
	user, ok := val.(dto.CurrentUser)
	if !ok {
		return dto.CurrentUser{}
	}
	
	return user
}