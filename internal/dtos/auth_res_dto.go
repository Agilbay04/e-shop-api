package dto

import "e-shop-api/internal/constants"

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type UserResponse struct {
	ID       string				`json:"id"`
	Username string    			`json:"username"`
	Email    string    			`json:"email"`
	Role     constant.UserRole	`json:"role"`
	Picture  string				`json:"picture"`
}