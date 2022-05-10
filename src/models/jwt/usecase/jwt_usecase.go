package usecase

import (
	"context"
	"time"
	"url-shortener/config"
	"url-shortener/src/models/users"
)

type JwtUsecase struct {
	Configuration  *config.Config
	UserRepo       users.UserRepository
	ContextTimeout time.Duration
}

func NewJwtUsecase(u users.UserRepository, to time.Duration, configuration *config.Config) *JwtUsecase {
	return &JwtUsecase{
		Configuration:  configuration,
		UserRepo:       u,
		ContextTimeout: to,
	}
}

func (h *JwtUsecase) getOneUser(c context.Context, id string) (*users.User, error) {

	ctx, cancel := context.WithTimeout(c, h.ContextTimeout)
	defer cancel()

	res, err := h.UserRepo.FindOne(ctx, id)
	if err != nil {
		return res, err
	}

	return res, nil
}
