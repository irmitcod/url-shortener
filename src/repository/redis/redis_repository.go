package redis

import "context"

type UrlRepository interface {
	CacheUrl(url string, buffer string)
	CacheInvalidUrl(url string)
	GetInvalidUrl(ctx context.Context, url string) (string, error)
	GetUrl(ctx context.Context, url string) (string, error)
	HasUrl(ctx context.Context, url string) (bool, error)
}
