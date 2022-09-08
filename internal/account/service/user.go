package service

import (
	"context"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
)

type UserSrv interface {
	Get(ctx context.Context, username string) (*model.User, error)
	GetById(ctx context.Context, id int) (*model.User, error)
}

type UserService struct {
	userStore store.UserStore
}

func NewUserSrv() UserSrv {
	return &UserService{userStore: store.NewUserStore()}
}

func (u *UserService) Get(ctx context.Context, username string) (*model.User, error) {
	return u.userStore.GetByUsername(ctx, username)
}

func (u *UserService) GetById(ctx context.Context, id int) (*model.User, error) {
	return u.userStore.GetById(ctx, id)
}
