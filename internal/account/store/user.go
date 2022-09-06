package store

import (
	"errors"
	"gorm.io/gorm"
	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/errno"
)

type UserStore interface {
	GetByUsername(username string) (*model.User, error)
	GetById(id int) (*model.User, error)
}

type user struct {
	db *gorm.DB
}

func NewUserStore() UserStore {
	return &user{db: db}
}

func (u *user) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := u.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.New(errno.ErrUserNotFound)
		}
		return nil, err
	}

	return user, nil
}

func (u *user) GetById(id int) (*model.User, error) {
	user := &model.User{}
	err := u.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.New(errno.ErrUserNotFound)
		}
		return nil, err
	}
	return user, nil
}
