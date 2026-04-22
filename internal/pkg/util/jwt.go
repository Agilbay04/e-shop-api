package util

import (
	"e-shop-api/internal/model"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	ID       uuid.UUID 		`json:"id"`
	Username string    		`json:"username"`
	Email    string    		`json:"email"`
	Role     model.UserRole	`json:"role"`
	Picture  string			`json:"picture"`
	jwt.RegisteredClaims
}

func GenerateToken(id uuid.UUID, username, email, picture string, role model.UserRole) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	ttl, err := strconv.Atoi(os.Getenv("JWT_TTL"))
	if err != nil {
		ttl = 3600
	} 

	claims := CustomClaims{
		ID:       id,
		Username: username,
		Email:    email,
		Role:     role,
		Picture:  picture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttl) * time.Second)),
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
