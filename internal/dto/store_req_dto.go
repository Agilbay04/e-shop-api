package dto

import "github.com/google/uuid"

type CreateStoreRequest struct {
	Name string `json:"name" binding:"required"`
	Description string `json:"description"`
	UserID uuid.UUID `json:"user_id" binding:"required"`
}


