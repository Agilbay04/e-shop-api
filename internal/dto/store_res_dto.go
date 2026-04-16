package dto

import "github.com/google/uuid"

type CreateStoreResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
}
