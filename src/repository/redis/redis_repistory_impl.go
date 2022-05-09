package redis

import (
	"context"
	"fmt"
	"url-shortener/config"
)

const (
	URLSNAMESAPCE       = "url"
	INVALIDURLNAMESAPCE = "invalid"
)

func NewUrlRepository(database *config.MemoryClient) UrlRepository {
	return &imageRepositoryImpl{
		redis: database,
	}
}

type imageRepositoryImpl struct {
	redis *config.MemoryClient
}

func (i *imageRepositoryImpl) GetInvalidUrl(ctx context.Context, url string) (string, error) {
	key := fmt.Sprintf("%s:%s", INVALIDURLNAMESAPCE, url)
	return i.redis.Client.Get(ctx, key).Result()
}

func (i *imageRepositoryImpl) CacheInvalidUrl(url string) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", INVALIDURLNAMESAPCE, url)
	err := i.redis.Client.Set(ctx, key, "invalid", 0).Err()
	fmt.Println(err)
}

func (i *imageRepositoryImpl) HasUrl(ctx context.Context, url string) (bool, error) {
	key := fmt.Sprintf("%s:%s", URLSNAMESAPCE, url)
	return i.redis.Client.Get(ctx, key).Bool()
}

func (i *imageRepositoryImpl) GetUrl(ctx context.Context, url string) (string, error) {
	key := fmt.Sprintf("%s:%s", URLSNAMESAPCE, url)
	return i.redis.Client.Get(ctx, key).Result()
}

func (i *imageRepositoryImpl) CacheUrl(url string, buffer []byte) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", URLSNAMESAPCE, url)
	i.redis.Client.Set(ctx, key, buffer, 0).Err()
}
