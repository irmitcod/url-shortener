package mongo

import (
	"context"
	"fmt"
	"url-shortener/config"
	"url-shortener/src/repository/redis"
)

const (
	URLSNAMESAPCE       = "url"
	INVALIDURLNAMESAPCE = "invalid"
)

func NewUrlRepository(database *config.MongoClient) redis.UrlRepository {
	return &urlRepositoryImpl{
		mongoCleint: database,
	}
}

type urlRepositoryImpl struct {
	mongoCleint *config.MongoClient
}

func (i *urlRepositoryImpl) GetInvalidUrl(ctx context.Context, url string) (string, error) {
	key := fmt.Sprintf("%s:%s", INVALIDURLNAMESAPCE, url)
	return i.mongoCleint.Client.Connect()
}

func (i *urlRepositoryImpl) CacheInvalidUrl(url string) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", INVALIDURLNAMESAPCE, url)
	err := i.redis.Client.Set(ctx, key, "invalid", 0).Err()
	fmt.Println(err)
}

func (i *urlRepositoryImpl) HasUrl(ctx context.Context, url string) (bool, error) {
	key := fmt.Sprintf("%s:%s", URLSNAMESAPCE, url)
	return i.redis.Client.Get(ctx, key).Bool()
}

func (i *urlRepositoryImpl) GetUrl(ctx context.Context, url string) (string, error) {
	key := fmt.Sprintf("%s:%s", URLSNAMESAPCE, url)
	return i.redis.Client.Get(ctx, key).Result()
}

func (i *urlRepositoryImpl) CacheUrl(url string, buffer []byte) {
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s", URLSNAMESAPCE, url)
	i.redis.Client.Set(ctx, key, buffer, 0).Err()
}
