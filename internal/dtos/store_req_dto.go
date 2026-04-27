package dtos

type CreateStoreRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	UserID      string 	  `json:"user_id" binding:"required"`
}

type UpdateStoreRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type ActivateStoreRequest struct {
	ID       string `json:"id" binding:"required"`
	IsActive bool `json:"is_active"`
}

func QueryStoreRequest() QueryStoreParam {
	return QueryStoreParam{
		PaginationParam: PaginationParam{
			Page:    1,
			Limit:   10,
			SortBy:  "created_at",
			OrderBy: "desc",
		},
	}
}

type QueryStoreParam struct {
	PaginationParam
	UserID *string `form:"user_id"`
}
