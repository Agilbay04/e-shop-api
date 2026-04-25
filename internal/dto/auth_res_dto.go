package dto

import "e-shop-api/internal/model"

type LoginResponse struct {
	Token string      	`json:"token,omitempty"`
}

type UserResponse struct {
	ID       string			`json:"id"`
	Username string    		`json:"username"`
	Email    string    		`json:"email"`
	Role     model.UserRole	`json:"role"`
	Picture  string			`json:"picture"`
}