package middleware

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/handler"
	"e-shop-api/internal/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there is any error
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
    
			// Check if err is implement IAppError
			if appErr, ok := err.(util.IAppError); ok {
				// Cause error implement IAppError, 
				// Get message and status code.
				statusCode := appErr.GetStatusCode()
				
				// Try get error message from err
				// If err is *util.CustomError
				if customErr, ok := err.(*util.CustomError); ok {
					// If status code is 400
					if statusCode == 400 && customErr.Errors != nil {
						// If errors is implement validator.ValidationErrors
						if rawErr, isErr := customErr.Errors.(error); isErr {
							customErr.Errors = util.FormatValidationError(rawErr)
						}
					}

					handler.RespondError(c, statusCode, customErr.Message, customErr.Errors, "")
					return
				}
			}

			// Fallback to internal server error
			handler.RespondError(c, http.StatusInternalServerError, "Internal Server Error", err.Error(), "")
			return
		}

		// If there is no error get payload
		if data, exists := c.Get("payload"); exists {
			message, _ := c.Get("message")
			isPagination, _ := c.Get("is_pagination")
    
			// Get status code, if nil set to 200
			status, exists := c.Get("status")
			statusCode := http.StatusOK
			if exists {
				statusCode = status.(int)
			}

			msgString := "Success"
			if message != nil {
				msgString = message.(string)
			}

			if isPagination == true {
				total, _ := c.Get("total_data")
        		filter, _ := c.Get("pagination_filter")
        		f := filter.(dto.PaginationParam)
        
        		handler.RespondPagination(c, data, total.(int64), f.Page, f.Limit, message.(string))
        		return
			}

			handler.RespondSuccess(c, statusCode, msgString, data)
		}
	}
}