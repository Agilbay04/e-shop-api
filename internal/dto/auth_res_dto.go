package dto

import "github.com/google/uuid"

type LoginResponse struct {
	Token string      	`json:"token,omitempty"`
	User  UserResponse	`json:"user"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
}