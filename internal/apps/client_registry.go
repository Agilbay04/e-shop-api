package apps

import "github.com/redis/go-redis/v9"

type ClientRegistry struct {
	Redis *redis.Client
}


func NewClientRegistry(rdb *redis.Client) *ClientRegistry {
	return &ClientRegistry{Redis: rdb}
}