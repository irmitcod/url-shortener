package mocks

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	"url-shortener/src/models/url_shortener"
)

type UrlRepository struct {
	mock.Mock
}

func (_m *UrlRepository) FindOne(ctx context.Context, id string) (*url_shortener.UrlShortener, error) {
	ret := _m.Called(ctx, id)

	var r0 *url_shortener.UrlShortener
	if rf, ok := ret.Get(0).(func(context.Context, string) *url_shortener.UrlShortener); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*url_shortener.UrlShortener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *UrlRepository) FindOneByKey(ctx context.Context, id string) (*url_shortener.UrlShortener, error) {
	ret := _m.Called(ctx, id)

	var r0 *url_shortener.UrlShortener
	if rf, ok := ret.Get(0).(func(context.Context, string) *url_shortener.UrlShortener); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*url_shortener.UrlShortener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *UrlRepository) UpdateOne(ctx context.Context, cat *url_shortener.UrlShortener, id string) (*url_shortener.UrlShortener, error) {
	//TODO implement me
	panic("implement me")
}

// InsertOne provides a mock function with given fields: ctx, u
func (_m *UrlRepository) InsertOne(ctx context.Context, u *url_shortener.UrlShortener) (*url_shortener.UrlShortener, error) {
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
