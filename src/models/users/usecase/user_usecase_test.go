package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
	_mocks "url-shortener/config/mocks"
	"url-shortener/src/models/mocks"
	"url-shortener/src/models/users"
)

func TestInsertOne(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUrlRepo := new(mocks.UrlRepository)
	mockLocalCaceh := new(_mocks.LocalCache)
	mockUser := &users.User{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "Behrouz",
		Username:  "master",
		Password:  "master@123",
		ApiKeys: struct {
			Key      string `json:"key" bson:"key"`
			ExpireAt int64  `json:"expire_at" bson:"expire_at"`
			CreateAt int64  `json:"create_at" bson:"create_at"`
		}{},
	}

	mockEmptyUser := &users.User{}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("InsertOne", mock.Anything, mock.AnythingOfType("*users.User")).Return(mockUser, nil).Once()
		u := NewUserUsecase(mockUserRepo, time.Second*2, mockUrlRepo, mockLocalCaceh)

		a, err := u.InsertOne(context.TODO(), mockUser)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockUserRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("InsertOne", mock.Anything, mock.AnythingOfType("*users.User")).Return(mockEmptyUser, errors.New("Unexpected")).Once()

		u := NewUserUsecase(mockUserRepo, time.Second*2, mockUrlRepo, mockLocalCaceh)

		a, err := u.InsertOne(context.TODO(), mockUser)

		assert.Error(t, err)
		assert.Equal(t, mockEmptyUser, a)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestFindOne(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockUrlRepo := new(mocks.UrlUserRepository)
	mockLocalCaceh := new(_mocks.LocalCache)
	mockUser := &users.User{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "Behrouz",
		Username:  "master",
		Password:  "master@123",
		ApiKeys: struct {
			Key      string `json:"key" bson:"key"`
			ExpireAt int64  `json:"expire_at" bson:"expire_at"`
			CreateAt int64  `json:"create_at" bson:"create_at"`
		}{},
	}

	mockEmptyUser := &users.User{}

	UserID := mock.Anything

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("FindOne", mock.Anything, mock.Anything).Return(mockUser, nil).Once()
		u := NewUserUsecase(mockUserRepo, time.Second*2, mockUrlRepo, mockLocalCaceh)

		a, err := u.FindOne(context.TODO(), UserID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockUserRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("FindOne", mock.Anything, mock.Anything).Return(mockEmptyUser, errors.New("Unexpected")).Once()

		u := NewUserUsecase(mockUserRepo, time.Second*2, mockUrlRepo, mockLocalCaceh)

		a, err := u.FindOne(context.TODO(), UserID)

		assert.Error(t, err)
		assert.Equal(t, mockEmptyUser, a)

		mockUserRepo.AssertExpectations(t)
	})

}
