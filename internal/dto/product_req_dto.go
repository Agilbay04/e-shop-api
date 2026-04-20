package dto

type QueryProductRequest struct {
	PaginationParam
	IsActive *bool   `form:"is_active"`
	StoreID  *string `form:"store_id"`
	ID       *string `form:"id"`
	MinPrice *int    `form:"min_price"`
	MaxPrice *int    `form:"max_price"`
}

type ActivateProductRequest struct {
	ID       string `json:"id" binding:"required"`
	IsActive *bool 	`json:"is_active" binding:"required"`
}

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       *int   `json:"price" binding:"required,gt=0"`
	Stock       *int   `json:"stock" binding:"required,min=0"`
	Unit        string `json:"unit" binding:"required"`
}

type UpdateProductRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Price       *int    `json:"price" binding:"omitempty,gt=0"`
	Stock       *int    `json:"stock" binding:"omitempty,min=0"`
	Unit        *string `json:"unit"`
	IsActive    *bool   `json:"is_active"`
}
