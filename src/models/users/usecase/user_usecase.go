package usecase

import (
	"context"
	"time"
	"url-shortener/config"
	"url-shortener/src/models/users"
	redis2 "url-shortener/src/repository/redis"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUsecase struct {
	userRepo       users.UserRepository
	contextTimeout time.Duration
	redisRepo      redis2.UrlRepository
	cache          config.LocalCache
}

func NewUserUsecase(u users.UserRepository, to time.Duration, repository redis2.UrlRepository, cache config.LocalCache) users.UserUsecase {
	return &userUsecase{
		userRepo:       u,
		contextTimeout: to,
		redisRepo:      repository,
		cache:          cache,
	}
}

func (user *userUsecase) InsertOne(c context.Context, m *users.User) (*users.User, error) {

	ctx, cancel := context.WithTimeout(c, user.contextTimeout)
	defer cancel()

	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()

	res, err := user.userRepo.InsertOne(ctx, m)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (user *userUsecase) FindOne(c context.Context, id string) (*users.User, error) {

	ctx, cancel := context.WithTimeout(c, user.contextTimeout)
	defer cancel()

	res, err := user.userRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (user *userUsecase) UpdateOne(c context.Context, m *users.User, id string) (*users.User, error) {

	ctx, cancel := context.WithTimeout(c, user.contextTimeout)
	defer cancel()

	res, err := user.userRepo.UpdateOne(ctx, m, id)
	if err != nil {
		return res, err
	}

	return res, nil
}
