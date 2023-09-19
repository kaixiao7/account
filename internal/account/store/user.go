package store

import (
	"context"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type UserStore interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetById(ctx context.Context, id int64) (*model.User, error)
}

type user struct {
}

var base_field = "id, username, phone, wx_id, gender, password, avatar, register_time, update_time"

func NewUserStore() UserStore {
	return &user{}
}

func (u *user) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	db := getDBFromContext(ctx)

	// sql := fmt.Sprintf("select * from users where username = ?", base_field)
	sql := db.Rebind("select * from users where username = ?")
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

func (u *user) GetById(ctx context.Context, id int64) (*model.User, error) {
	db := getDBFromContext(ctx)

	// sql := fmt.Sprintf("select * from users where id = ?", base_field)
	sql := db.Rebind("select * from users where id = ?")
	user := model.User{}
	err := db.Get(&user, sql, id)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
