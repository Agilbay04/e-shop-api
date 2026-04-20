package util

import (
	"e-shop-api/internal/model"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	ID       uuid.UUID 		`json:"id"`
	Username string    		`json:"username"`
	Email    string    		`json:"email"`
	Role     model.UserRole	`json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(id uuid.UUID, username, email string, role model.UserRole) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	claims := CustomClaims{
		ID:       id,
		Username: username,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // JWT TTL 1 hour
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
