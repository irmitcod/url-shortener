package config

import (
	"context"
	"errors"
	_redisCach "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
	"time"
)

var ErrNon200 = errors.New("received non 200 response code")
var ErrImageNotFound = errors.New("urls not found")

type MemoryClient struct {
	Client              *_redisCach.Cache
	MaxWidth, MaxHeight int
}

func NewMemoryClient(c *Config) *MemoryClient {

	o := &redis.Options{
		Addr:     os.Getenv("DATABASE_REDIS_HOST") + ":" + os.Getenv("DATABASE_REDIS_PORT"),
		Password: os.Getenv("DATABASE_REDIS_PASSWORD"), // no password set
		DB:       0,                                    // use default DB
	}

	ring := redis.NewRing(&redis.RingOptions{

		NewClient: func(name string, opt *redis.Options) *redis.Client {
			opt.Addr = o.Addr
			opt.DB = o.DB
			opt.Username = o.Username
			opt.Password = o.Password
			return redis.NewClient(opt)
		},
	})
	//rdb := redis.NewRing(&redis.RingOptions{
	//	NewClient: func(opt *redis.Options) *redis.NewClient {
	//		user, pass := userPassForAddr(opt.Addr)
	//		opt.Username = user
	//		opt.Password = pass
	//
	//		return redis.NewClient(opt)
	//	},
	//})

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

func ConvertStringToInt(str string) int {
	strConvert, _ := strconv.Atoi(str)
	return strConvert
}
