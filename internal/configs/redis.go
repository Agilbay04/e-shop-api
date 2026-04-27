package configs

import (
	"context"
	"e-shop-api/internal/pkg/logger"
	"e-shop-api/internal/pkg/utils"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func ConnectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,                          
	})

	// Retries to connect to redis
	err := utils.AutoRetry(func() error {
        return rdb.Ping(context.Background()).Err()
    })

	// Return error if failed connect to redis after retries
	if err != nil {
		logger.L.Fatal("Failed connect to Redis:", zap.Error(err))
		panic(fmt.Sprintf("Failed connect to Redis: %v", err))
	}

	// Success connect to redis
	logger.L.Info("Connected to Redis!")
	return rdb
}