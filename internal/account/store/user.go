package store

import (
	"context"
	"fmt"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type UserStore interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetById(ctx context.Context, id int) (*model.User, error)
}

type user struct {
}

var base_field = "id, username, ifnull(phone, '') as phone, ifnull(wx_id, '') as wx_id, gender, password, avatar, register_time, update_time"

func NewUserStore() UserStore {
	return &user{}
}

func (u *user) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	db := getDBFromContext(ctx)

	sql := fmt.Sprintf("select %s from user where username = ?", base_field)
	user := model.User{}
	err := db.Get(&user, sql, username)
	if err != nil {
		return nil, errors.Wrap(err, "get by username")
	}

	// row := db.QueryRowx("select * from user where username = ? limit 1", username)
	// err := row.StructScan(&user)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "get by username")
	// }

	return &user, nil
}

func (u *user) GetById(ctx context.Context, id int) (*model.User, error) {
	db := getDBFromContext(ctx)

	sql := fmt.Sprintf("select %s from user where id = ?", base_field)
	user := model.User{}
	err := db.Get(&user, sql, id)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
