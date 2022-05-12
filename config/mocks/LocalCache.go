package mocks

import (
	"github.com/stretchr/testify/mock"
	"time"
	"url-shortener/config"
)

type LocalCache struct {
	mock.Mock
}

func (_m *LocalCache) Get(key string) (interface{}, bool) {
	return _m.Get(key)
}

func (_m *LocalCache) Set(key string, val interface{}) bool {
	return _m.Set(key, val)
}

func (_m *LocalCache) SetWithTTL(key, value interface{}, cost int64, ttl time.Duration) bool {
	return _m.SetWithTTL(key, value, cost, ttl)
}

func (_m *LocalCache) Client() config.LocalCache {
	ret := _m.Called()
	var r0 config.LocalCache
	if rf, ok := ret.Get(0).(func() config.LocalCache); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(config.LocalCache)
		}
	}

	return r0
}
