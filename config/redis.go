package config

import (
	"errors"
	"github.com/go-redis/redis/v8"
)

var ErrNon200 = errors.New("received non 200 response code")
var ErrImageNotFound = errors.New("urls not found")

type MemoryClient struct {
	Client              *redis.Client
	MaxWidth, MaxHeight int
}

func NewMemoryClient(c *Config) *MemoryClient {

	// Creating MemoryClient
	mc := MemoryClient{
		Client: redis.NewClient(&redis.Options{
			Addr:     c.RedisAddress,
			Password: c.RedisPassword,
			DB:       c.RedisDB,
		}),
		MaxWidth:  c.MaxWidth,
		MaxHeight: c.MaxHeight,
	}
	return &mc
}
