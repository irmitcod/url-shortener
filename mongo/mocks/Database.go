package mocks

import (
	"github.com/stretchr/testify/mock"
	"url-shortener/mongo"
)

type Database struct {
	mock.Mock
}

func (_m *Database) Client() mongo.Client {
	ret := _m.Called()
	var r0 mongo.Client
	if rf, ok := ret.Get(0).(func() mongo.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mongo.Client)
		}
	}

	return r0
}

// Collection provides a mock function with given fields: _a0
func (_m *Database) Collection(_a0 string) mongo.Collection {
	ret := _m.Called(_a0)

	var r0 mongo.Collection
	if rf, ok := ret.Get(0).(func(string) mongo.Collection); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mongo.Collection)
		}
	}

	return r0
}
