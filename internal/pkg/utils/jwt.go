package utils

import (
	"e-shop-api/internal/constants"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	ID       string 			`json:"id"`
	Username string    			`json:"username"`
	Email    string    			`json:"email"`
	Role     constant.UserRole	`json:"role"`
	Picture  string				`json:"picture"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(id, username, email, picture string, role constant.UserRole) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	ttl := GetEnvTime("JWT_TTL", constant.JwtTtl)

	claims := CustomClaims{
		ID:       id,
		Username: username,
		Email:    email,
		Role:     role,
		Picture:  picture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateAccessToken(id, username, email, picture string, role constant.UserRole) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	ttl := GetEnvTime("JWT_ACCESS_TTL", constant.JwtAccessTtl)

	claims := CustomClaims{
		ID:       id,
		Username: username,
		Email:    email,
		Role:     role,
		Picture:  picture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateRefreshToken(userID string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	ttl := GetEnvTime("JWT_REFRESH_TTL", constant.JwtRefreshTtl)

	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
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

func ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid refresh token")
}
