package mocks

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	"url-shortener/src/models/url_shortener"
)

type UrlUserRepository struct {
	mock.Mock
}

func (_m *UrlUserRepository) CacheUrl(url string, buffer string) {
	//TODO implement me
	panic("implement me")
}

func (_m *UrlUserRepository) CacheInvalidUrl(url string) {
	//TODO implement me
	panic("implement me")
}

func (_m *UrlUserRepository) GetInvalidUrl(ctx context.Context, url string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (_m *UrlUserRepository) GetUrl(ctx context.Context, url string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (_m *UrlUserRepository) HasUrl(ctx context.Context, url string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

// InsertOne provides a mock function with given fields: ctx, u
func (_m *UrlUserRepository) InsertOne(ctx context.Context, u *url_shortener.UrlShortener) (*url_shortener.UrlShortener, error) {
	ret := _m.Called(ctx, u)

	var r0 *url_shortener.UrlShortener
	if rf, ok := ret.Get(0).(func(context.Context, *url_shortener.UrlShortener) *url_shortener.UrlShortener); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*url_shortener.UrlShortener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *url_shortener.UrlShortener) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
