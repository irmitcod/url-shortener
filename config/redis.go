package config

import (
	"context"
	"errors"
	_redisCach "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var ErrNon200 = errors.New("received non 200 response code")
var ErrImageNotFound = errors.New("urls not found")

type MemoryClient struct {
	Client              *_redisCach.Cache
	MaxWidth, MaxHeight int
}

func NewMemoryClient(c *Config) *MemoryClient {

	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"localhost": ":6379",
		},
	})
	err := ring.ForEachShard(context.Background(), func(ctx context.Context, shard *redis.Client) error {
		//shard.FlushDB(ctx)
		return shard.Ping(ctx).Err()
	})
	if err != nil {
		log.Panicln(err)
	}

	mycache := _redisCach.New(&_redisCach.Options{
		Redis:      ring,
		LocalCache: _redisCach.NewTinyLFU(1000, time.Hour),
	})

	// Creating MemoryClient
	mc := MemoryClient{
		Client: mycache,
	}
	return &mc
}
