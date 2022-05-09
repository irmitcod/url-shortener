package config

import (
	"github.com/dgraph-io/ristretto"
	"time"
)

func NewCache() LocalCache {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil
	}

	return &cache{c: c}
}

type cache struct {
	c *ristretto.Cache
}

func (c cache) SetWithTTL(key, value interface{}, cost int64, ttl time.Duration) bool {
	return c.c.SetWithTTL(key, value, cost, ttl)
}

func (c cache) Get(key string) (interface{}, bool) {
	return c.c.Get(key)
}

func (c cache) Set(key string, val interface{}) bool {
	return c.c.Set(key, val, 0)
}

type LocalCache interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{}) bool
	SetWithTTL(key, value interface{}, cost int64, ttl time.Duration) bool
}
