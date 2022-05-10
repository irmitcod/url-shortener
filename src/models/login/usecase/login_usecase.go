package usecase

import (
	"context"

	"time"
	"url-shortener/src/models/users"
)

type loginUsecase struct {
	userRepo       users.UserRepository
	contextTimeout time.Duration
}

func (l loginUsecase) GetUser(ctx context.Context, username string, password string) (*users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	res, err := l.userRepo.GetByCredential(ctx, username, password)
	if err != nil {
		return res, err
	}

	return res, nil
}

func NewLoginUsecase(u users.UserRepository, to time.Duration) users.LoginUsecase {
	return &loginUsecase{
		userRepo:       u,
		contextTimeout: to,
	}
}
