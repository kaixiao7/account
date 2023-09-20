package service

import (
	"context"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/auth"
	"kaixiao7/account/internal/pkg/errno"
)

type UserSrv interface {
	Get(ctx context.Context, username string) (*model.User, error)
	GetById(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	ChangePassword(ctx context.Context, oldPwd, newPwd string, userId int64) error
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

func (u *UserService) GetById(ctx context.Context, id int64) (*model.User, error) {
	return u.userStore.GetById(ctx, id)
}

func (u *UserService) Update(ctx context.Context, user *model.User) error {
	return u.userStore.Update(ctx, user)
}

func (u *UserService) ChangePassword(ctx context.Context, oldPwd, newPwd string, userId int64) error {
	user, err := u.userStore.GetById(ctx, userId)
	if err != nil {
		return nil
	}

	// 比较旧密码
	if err := auth.Compare(user.Password, oldPwd); err != nil {
		return errno.New(errno.ErrOldPasswordIncorrect)
	}

	encrypt, err := auth.Encrypt(newPwd)
	if err != nil {
		return err
	}

	return u.userStore.ChangePassword(ctx, encrypt, userId)
}
