package dtos

import (
	"github.com/google/uuid"
)

type CreateStoreResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
}

type StoreResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	DeletedAt   string    `json:"deleted_at,omitempty"`
}

type StoreListResponse struct {
	Items interface{} `json:"items"`
	Meta  MetaData    `json:"meta"`
}
