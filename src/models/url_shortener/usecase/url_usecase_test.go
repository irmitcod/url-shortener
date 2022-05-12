package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
	"url-shortener/config"
	"url-shortener/src/models/mocks"
	"url-shortener/src/models/url_shortener"
	"url-shortener/src/utils/lfu"
)

func TestInsertOne(t *testing.T) {
	mockRedisRepository := new(mocks.RedisRepository)
	mockUrlRepo := new(mocks.UrlRepository)
	mockLocalCaceh := config.NewCache()
	mockLfuCache := lfu.New()
	mockUser := &url_shortener.UrlShortener{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		OriginalURL: "https://google.com",
		UserID:      primitive.NewObjectID(),
	}

	mockEmptyUser := &url_shortener.UrlShortener{}

	t.Run("success", func(t *testing.T) {
		mockUrlRepo.On("InsertOne", mock.Anything, mock.AnythingOfType("*url_shortener.UrlShortener")).Return(mockUser, nil).Once()
		u := NewUrlUsecase(mockUrlRepo, time.Second*2, mockRedisRepository, mockLocalCaceh, mockLfuCache, nil)

		a, err := u.InsertOne(context.TODO(), mockUser)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockUrlRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockUrlRepo.On("InsertOne", mock.Anything, mock.AnythingOfType("*url_shortener.UrlShortener")).Return(mockEmptyUser, errors.New("Unexpected")).Once()

		u := NewUrlUsecase(mockUrlRepo, time.Second*2, mockRedisRepository, mockLocalCaceh, mockLfuCache, nil)

		a, err := u.InsertOne(context.TODO(), mockUser)

		assert.Error(t, err)
		assert.Equal(t, mockEmptyUser, a)

		mockUrlRepo.AssertExpectations(t)
	})
}
