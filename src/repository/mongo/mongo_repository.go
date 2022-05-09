package mongo

import (
	"context"
	"url-shortener/src/models/urls"
)

type MongoRepository interface {
	CacheUrl(url urls.Service, buffer []byte)
	CacheInvalidUrl(url string)
	GetInvalidUrl(ctx context.Context, url string) (string, error)
	GetUrl(ctx context.Context, url string) (string, error)
	HasUrl(ctx context.Context, url string) (bool, error)
}
