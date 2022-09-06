package service

import (
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
)

type UserSrv interface {
	Get(username string) (*model.User, error)
	GetById(id int) (*model.User, error)
}

type UserService struct {
	userStore store.UserStore
}

func NewUserSrv() UserSrv {
	return &UserService{userStore: store.NewUserStore()}
}

func (u *UserService) Get(username string) (*model.User, error) {
	return u.userStore.GetByUsername(username)
}

func (u *UserService) GetById(id int) (*model.User, error) {
	return u.userStore.GetById(id)
}
