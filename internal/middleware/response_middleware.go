package middleware

import (
	"e-shop-api/internal/dto"
	"e-shop-api/internal/handler"
	"e-shop-api/internal/pkg/util"
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

					handler.RespondError(ctx, statusCode, customErr.Message, customErr.Errors, "")
					return
				}
			}

			// Fallback to internal server error
			handler.RespondError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error(), "")
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
        
        		handler.RespondPagination(ctx, data, total.(int64), f.Page, f.Limit, message.(string))
        		return
			}

			handler.RespondSuccess(ctx, statusCode, msgString, data)
		}
	}
}