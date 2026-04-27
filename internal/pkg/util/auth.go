package util

import (
	"crypto/rand"
	"e-shop-api/internal/dto"
	"math/big"

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