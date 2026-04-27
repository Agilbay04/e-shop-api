package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Set Redis Cache
func SetCache[T any](rdb *redis.Client, key string, data T, ttl time.Duration) error {
	ctx := context.Background()
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, jsonData, ttl).Err()
}

// Get Redis Cache
func GetCache[T any](rdb *redis.Client, key string) (T, error) {
	ctx := context.Background()
	var data T

	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return data, err // return nil if key not found
	}

	err = json.Unmarshal([]byte(val), &data)
	return data, err
}

// Delete Redis Cache
func DeleteCache(rdb *redis.Client, key string) error {
	return rdb.Del(context.Background(), key).Err()
}

// Delete Cache By Pattern
func DeleteCacheByPattern(rdb *redis.Client, pattern string) error {
	ctx := context.Background()
	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := rdb.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

func IsRateLimited(ctx context.Context, rdb *redis.Client, key string, limit int, window time.Duration) (bool, int, error) {
	windowEpoch := time.Now().Unix() / int64(window.Seconds())
	rateKey := fmt.Sprintf("%s:%d", key, windowEpoch)

	count, err := rdb.Incr(ctx, rateKey).Result()
	if err != nil {
		return true, 0, fmt.Errorf("redis incr failed: %w", err)
	}

	if count == 1 {
		rdb.Expire(ctx, rateKey, window*3/2)
	}

	return count > int64(limit), int(count), nil
}
