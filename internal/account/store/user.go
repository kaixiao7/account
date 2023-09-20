package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/errno"

	"github.com/pkg/errors"
)

type UserStore interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetById(ctx context.Context, id int64) (*model.User, error)

	Update(ctx context.Context, user *model.User) error

	ChangePassword(ctx context.Context, pwd string, userId int64) error
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
	querySql := db.Rebind("select * from users where username = ?")
	user := model.User{}
	err := db.Get(&user, querySql, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errno.New(errno.ErrPasswordIncorrect)
		}
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

func (u *user) Update(ctx context.Context, user *model.User) error {
	db := getDBFromContext(ctx)

	updateSql := db.Rebind("update users set username=?, avatar=?, phone=?, gender=?, update_time=? where id = ?")

	_, err := db.Exec(updateSql, user.Username, user.Avatar, user.Phone, user.Gender, user.UpdateTime, user.Id)

	if err != nil {
		return errors.Wrap(err, "user update store")
	}
	return nil
}

func (u *user) ChangePassword(ctx context.Context, pwd string, userId int64) error {
	db := getDBFromContext(ctx)

	updateSql := db.Rebind("update users set password = ? where id = ?")

	_, err := db.Exec(updateSql, pwd, userId)

	if err != nil {
		return errors.Wrap(err, "user change password store")
	}
	return nil
}
