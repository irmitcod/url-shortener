package usecase

import (
	"context"
	log "github.com/sirupsen/logrus"

	"time"
	"url-shortener/src/models/users"
)

type loginUsecase struct {
	userRepo       users.UserRepository
	contextTimeout time.Duration
	entry          *log.Entry
}

func (l loginUsecase) GetUser(ctx context.Context, username string, password string) (*users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	res, err := l.userRepo.GetByCredential(ctx, username, password)
	if err != nil {
		l.entry.Error("error from GetByCredential ", err)
		return res, err
	}

	return res, nil
}

func NewLoginUsecase(u users.UserRepository, to time.Duration, entry *log.Entry) users.LoginUsecase {
	return &loginUsecase{
		userRepo:       u,
		contextTimeout: to,
		entry:          entry,
	}
}
