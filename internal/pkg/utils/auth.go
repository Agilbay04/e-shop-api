package utils

import (
	"crypto/rand"
	"e-shop-api/internal/dtos"
	"math/big"

	"github.com/gin-gonic/gin"
)

const UserContextKey = "currentUser"

// GetCurrentUser to extract CurrentUser from gin.Context
func GetCurrentUser(ctx *gin.Context) dtos.CurrentUser {
	val, exists := ctx.Get(UserContextKey)
	if !exists {
		return dtos.CurrentUser{}
	}

	// type assertion
	user, ok := val.(dtos.CurrentUser)
	if !ok {
		return dtos.CurrentUser{}
	}
	
	return user
}

func GenerateRandomString(length int) (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		b[i] = chars[num.Int64()]
	}
	return string(b), nil
}