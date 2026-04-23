package util

import (
	"context"
	"encoding/json"
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

func IsRateLimited(rdb *redis.Client, key string, duration time.Duration) bool {
	ctx := context.Background()
	
	// Check if key exists
	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false // If redis error, set false
	}

	if exists > 0 {
		return true // Limit hit
	}

	// Set limit
	rdb.Set(ctx, key, "1", duration)
	return false
}
