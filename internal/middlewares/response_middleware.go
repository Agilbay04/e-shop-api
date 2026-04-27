package middlewares

import (
	"e-shop-api/internal/dtos"
	"e-shop-api/internal/handlers"
	"e-shop-api/internal/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResponseMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		// Check if there is any error
		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err
    
			// Check if err is implement IAppError
			if appErr, ok := err.(utils.IAppError); ok {
				// Cause error implement IAppError, 
				// Get message and status code.
				statusCode := appErr.GetStatusCode()
				
				// Try get error message from err
				// If err is *utils.CustomError
				if customErr, ok := err.(*utils.CustomError); ok {
					// If status code is 400
					if statusCode == 400 && customErr.Errors != nil {
						// If errors is implement validator.ValidationErrors
						if rawErr, isErr := customErr.Errors.(error); isErr {
							customErr.Errors = utils.FormatValidationError(rawErr)
						}
					}

					handlers.RespondError(ctx, statusCode, customErr.Message, customErr.Errors, "")
					return
				}
			}

			// Fallback to internal server error
			handlers.RespondError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error(), "")
			return
		}

		// If there is no error get payload
		if data, exists := ctx.Get("payload"); exists {
			message, _ := ctx.Get("message")
			isPagination, _ := ctx.Get("is_pagination")
    
			// Get status code, if nil set to 200
			status, exists := ctx.Get("status")
			statusCode := http.StatusOK
			if exists {
				statusCode = status.(int)
			}

			msgString := "Success"
			if message != nil {
				msgString = message.(string)
			}

			if isPagination == true {
				total, _ := ctx.Get("total_data")
        		filter, _ := ctx.Get("pagination_filter")
        		f := filter.(dto.PaginationParam)
        
        		handlers.RespondPagination(ctx, data, total.(int64), f.Page, f.Limit, message.(string))
        		return
			}

			handlers.RespondSuccess(ctx, statusCode, msgString, data)
		}
	}
}