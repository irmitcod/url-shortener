package mocks

import (
	"context"
	mock "github.com/stretchr/testify/mock"
)

type RedisRepository struct {
	mock.Mock
}

func (r RedisRepository) CacheUrl(url string, buffer string) {
	//TODO implement me
	panic("implement me")
}

func (r RedisRepository) CacheInvalidUrl(url string) {
	//TODO implement me
	panic("implement me")
}

func (r RedisRepository) GetInvalidUrl(ctx context.Context, url string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisRepository) GetUrl(ctx context.Context, url string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisRepository) HasUrl(ctx context.Context, url string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
