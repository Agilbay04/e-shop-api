package middleware

import (
	"e-shop-api/internal/pkg/logger"
	"e-shop-api/internal/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiter middleware for limiting requests
func RateLimiter(rdb *redis.Client, limitKey string, duration time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get IP address for key rate limit
		key := "limit:" + limitKey + ":" + ctx.ClientIP()

		if util.IsRateLimited(rdb, key, duration) {
			logger.Log.Warn("To many request, please try again later.")
			ctx.Error(util.ToManyRequestException("To many request, please try again later."))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}