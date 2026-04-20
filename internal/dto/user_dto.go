package dto

import (
	"e-shop-api/internal/model"

	"github.com/google/uuid"
)

type CurrentUser struct {
    ID       uuid.UUID      `json:"id"`
    Username string         `json:"username"`
    Email    string         `json:"email"`
    Role     model.UserRole `json:"role"`
}