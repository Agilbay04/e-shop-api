package middlewares

import (
	"strconv"

	"e-shop-api/internal/pkg/logger"
	"e-shop-api/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func RateLimiter(rdb *redis.Client, limitKey string, limit int, window time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := "limit:" + limitKey + ":" + ctx.ClientIP()

		limited, count, err := utils.IsRateLimited(ctx, rdb, key, limit, window)
		if err != nil {
			logger.L.Error("Rate limiter error", zap.Error(err))
			errResp := utils.InternalServerErrorException("Service unavailable")
			ctx.Error(errResp)
			ctx.Abort()
			return
		}

		retryAfter := window.Seconds()
		if limited {
			retryAfter = window.Seconds() - float64(count-1)*window.Seconds()/float64(limit+1)
			logger.L.Warn("Rate limit exceeded", zap.String("key", key), zap.Int("count", count))
			errResp := utils.ToManyRequestException("Too many requests, please try again later.")
			ctx.Error(errResp)
			ctx.Abort()
			return
		}

		remaining := limit - int(count) + 1
		ctx.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		ctx.Header("Retry-After", strconv.FormatFloat(retryAfter, 'f', 0, 64))
		ctx.Next()
	}
}