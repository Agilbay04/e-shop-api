package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const AppVersion = "1.0.0"

// Response standar untuk sukses
type APIResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Success    bool        `json:"success"`
	Version    string      `json:"version"`
}

// Response standar untuk error
type ErrorResponse struct {
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"error_message"`
	StackTrace   string      `json:"stack_trace,omitempty"` // omitempty: sembunyikan jika kosong
	Errors       interface{} `json:"errors"`
	StatusCode   int         `json:"status_code"`
	Version      string      `json:"version"`
}

// RespondSuccess mengirim response sukses yang seragam
func RespondSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		StatusCode: code,
		Message:    message,
		Data:       data,
		Success:    true,
		Version:    AppVersion,
	})
}

// RespondError mengirim response error yang seragam
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

func Success(c *gin.Context, data interface{}, message string) {
    c.Set("payload", data)
    c.Set("message", message)
    c.Set("status", http.StatusOK)
}