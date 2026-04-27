package handlers

import (
	"e-shop-api/internal/dtos"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

const AppVersion = "1.0.0"

type APIResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Success    bool        `json:"success"`
	Version    string      `json:"version"`
}

type ErrorResponse struct {
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"error_message"`
	StackTrace   string      `json:"stack_trace,omitempty"` // omitempty: hide if empty
	Errors       interface{} `json:"errors"`
	StatusCode   int         `json:"status_code"`
	Version      string      `json:"version"`
}

func RespondSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		StatusCode: code,
		Message:    message,
		Data:       data,
		Success:    true,
		Version:    AppVersion,
	})
}

func RespondError(c *gin.Context, code int, message string, errs interface{}, stack string) {
	c.JSON(code, ErrorResponse{
		Success:      false,
		ErrorMessage: message,
		Errors:       errs,
		StackTrace:   stack,
		StatusCode:   code,
		Version:      AppVersion,
	})
}

func RespondPagination(c *gin.Context, data interface{}, totalData int64, page, limit int, message string) {
	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))

	res := dtos.PaginationResponse{
		Items: data,
		Meta: dtos.MetaData{
			CurrentPage: page,
			TotalPage:   totalPage,
			TotalData:   totalData,
			Limit:       limit,
		},
	}

	RespondSuccess(c, 200, message, res)
}

func OkPagination(c *gin.Context, data interface{}, total int64, filter dtos.PaginationParam, message string) {
	c.Set("payload", data)
	c.Set("total_data", total)
	c.Set("pagination_filter", filter)
	c.Set("message", message)
	c.Set("is_pagination", true)
}

func Ok(c *gin.Context, data interface{}, message string) {
	c.Set("payload", data)
	c.Set("message", message)
	c.Set("status", http.StatusOK)
}

func Created(c *gin.Context, data interface{}, message string) {
	c.Set("payload", data)
	c.Set("message", message)
	c.Set("status", http.StatusCreated)
}
