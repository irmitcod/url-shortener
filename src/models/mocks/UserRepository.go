package mocks

import (
	"context"
	mock "github.com/stretchr/testify/mock"
	"url-shortener/src/models/users"
)

type UserRepository struct {
	mock.Mock
}

func (_m *UserRepository) FindOne(ctx context.Context, id string) (*users.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *users.User
	if rf, ok := ret.Get(0).(func(context.Context, string) *users.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
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

// InsertOne provides a mock function with given fields: ctx, u
func (_m *UserRepository) InsertOne(ctx context.Context, u *users.User) (*users.User, error) {
	ret := _m.Called(ctx, u)

	var r0 *users.User
	if rf, ok := ret.Get(0).(func(context.Context, *users.User) *users.User); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *users.User) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOne provides a mock function with given fields: ctx, user, id
func (_m *UserRepository) UpdateOne(ctx context.Context, user *users.User, id string) (*users.User, error) {
	ret := _m.Called(ctx, user, id)

	var r0 *users.User
	if rf, ok := ret.Get(0).(func(context.Context, *users.User, string) *users.User); ok {
		r0 = rf(ctx, user, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *users.User, string) error); ok {
		r1 = rf(ctx, user, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByCredential provides a mock function with given fields: ctx, username, password
func (_m *UserRepository) GetByCredential(ctx context.Context, username string, password string) (*users.User, error) {
	ret := _m.Called(ctx, username, password)

	var r0 *users.User
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *users.User); ok {
		r0 = rf(ctx, username, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*users.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
