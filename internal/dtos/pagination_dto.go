package dto

import "regexp"

var sortByRegex = regexp.MustCompile(`^[a-z_]+$`)

func ValidateSortByPattern(sortBy string) bool {
	return sortByRegex.MatchString(sortBy)
}

func IsAllowedSortBy(sortBy string, allowed map[string]bool) bool {
	return allowed[sortBy]
}

type PaginationParam struct {
	Page    int    `form:"page,default=1" binding:"min=1"`
	Limit  	int    `form:"limit,default=10" binding:"min=1,max=1000"`
	Search 	string `form:"search"`
	SortBy 	string `form:"sort_by,default=created_at"`
	OrderBy string `form:"order_by,default=desc" binding:"oneof=asc desc"`
}

type MetaData struct {
	CurrentPage int   `json:"current_page"`
	TotalPage   int   `json:"total_page"`
	TotalData   int64 `json:"total_data"`
	Limit       int   `json:"limit"`
}

type PaginationResponse struct {
	Items interface{} `json:"items"`
	Meta MetaData     `json:"meta"`
}