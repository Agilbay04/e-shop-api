package dto

type QueryProductRequest struct {
	PaginationParam
	StoreID *string `form:"store_id" binding:"required"`
	ID *string `form:"id"`
	MinPrice *int `form:"min_price"`
	MaxPrice *int `form:"max_price"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       *int 	`json:"price" binding:"required,gt=0"`
	Stock       *int    `json:"stock" binding:"required,min=0"`
	Unit        string  `json:"unit" binding:"required"`
}