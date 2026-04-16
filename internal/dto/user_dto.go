package dto

import "github.com/google/uuid"

type CurrentUser struct {
    ID       uuid.UUID `json:"id"`
    Username string    `json:"username"`
    Email    string    `json:"email"`
    Role     string    `json:"role"`
}