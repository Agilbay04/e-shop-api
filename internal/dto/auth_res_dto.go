package dto

import (
	"e-shop-api/internal/model"

	"github.com/google/uuid"
)

type LoginResponse struct {
	Token string      	`json:"token,omitempty"`
	User  UserResponse	`json:"user"`
}

type UserResponse struct {
	ID       uuid.UUID 			`json:"id"`
	Username string    			`json:"username"`
	Email    string    			`json:"email"`
	Role     model.UserRole    `json:"role"`
}