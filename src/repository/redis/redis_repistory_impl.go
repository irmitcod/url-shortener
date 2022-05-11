package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"log"
	"time"
	"url-shortener/config"
)

const (
	URLSNAMESAPCE       = "url"
	INVALIDURLNAMESAPCE = "invalid"
)

type Object struct {
	URL string
	Num int
}

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
	var value string
	if err := i.redis.Client.Get(ctx, key, &value); err != nil {
		return "", err
	}
	return value, nil
}

func (i *imageRepositoryImpl) CacheInvalidUrl(url string) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", INVALIDURLNAMESAPCE, url)

	obj := &Object{
		URL: "invalid url",
	}
	if err := i.redis.Client.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: obj,
		TTL:   time.Hour,
	}); err != nil {
		log.Println(err)
	}

}

func (i *imageRepositoryImpl) HasUrl(ctx context.Context, url string) (bool, error) {

	key := fmt.Sprintf("%s:%s", URLSNAMESAPCE, url)

	return i.redis.Client.Exists(ctx, key), nil

}

func (i *imageRepositoryImpl) GetUrl(ctx context.Context, url string) (string, error) {
	key := fmt.Sprintf("%s:%s", URLSNAMESAPCE, url)
	obj := &Object{}

	if err := i.redis.Client.Get(ctx, key, obj); err != nil {
		return "", err
	}
	return obj.URL, nil
}

func (i *imageRepositoryImpl) CacheUrl(url string, buffer string) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", URLSNAMESAPCE, url)

	obj := &Object{
		URL: buffer,
	}
	if err := i.redis.Client.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: obj,
		TTL:   time.Hour,
	}); err != nil {
		log.Println(err)
	}

}
