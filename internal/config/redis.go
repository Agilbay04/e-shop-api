package config

import (
	"context"
	"e-shop-api/internal/pkg/logger"
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

	// Test connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.L.Fatal("Failed connect to Redis:", zap.Error(err))
		panic(fmt.Sprintf("Failed connect to Redis: %v", err))
	}

	logger.L.Info("Connected to Redis!")
	return rdb
}