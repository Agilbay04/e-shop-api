package utils

import (
	"os"
	"strconv"
	"time"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvInt(key string, fallback string) int {
	value := GetEnv(key, fallback)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return intValue
}

func GetEnvTime(key string, fallback string) time.Duration {
	value := GetEnv(key, fallback)
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0
	}
	return duration
}

func GetEnvBool(key string, fallback string) bool {
	value := GetEnv(key, fallback)
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return boolValue
}
