package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
	"url-shortener/src/models/mocks"
	"url-shortener/src/models/users"
)

func TestGetUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
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
	username := "master"
	password := "master@123"
	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("GetByCredential", mock.Anything, mock.Anything, mock.Anything).Return(mockUser, nil).Once()
		u := NewLoginUsecase(mockUserRepo, time.Second*2)
		a, err := u.GetUser(context.TODO(), username, password)

		assert.NoError(t, err)
		assert.NotNil(t, a)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockUserRepo.On("GetByCredential", mock.Anything, mock.Anything, mock.Anything).Return(mockEmptyUser, errors.New("Unexpected")).Once()

		u := NewLoginUsecase(mockUserRepo, time.Second*2)

		a, err := u.GetUser(context.TODO(), username, password)

		assert.Error(t, err)
		assert.Equal(t, mockEmptyUser, a)

		mockUserRepo.AssertExpectations(t)
	})
}
