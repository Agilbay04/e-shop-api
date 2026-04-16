package util

import (
	"e-shop-api/internal/dto"
	"github.com/gin-gonic/gin"
)

const UserContextKey = "currentUser"

// GetCurrentUser to extract CurrentUser from gin.Context
func GetCurrentUser(c *gin.Context) dto.CurrentUser {
	val, exists := c.Get(UserContextKey)
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