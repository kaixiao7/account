package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/constant"

	"github.com/pkg/errors"
)

// AccountStore 账户
type AccountStore interface {
	// Add 添加账户
	Add(ctx context.Context, account *model.Account) error
	// Update 更新账户
	Update(ctx context.Context, account *model.Account) error
	// Delete 删除账户
	// 执行的是逻辑删除，将字段del置为1
	Delete(ctx context.Context, id int) error

	// QueryAllByUserId 根据用户id查询其所有账户
	QueryAllByUserId(ctx context.Context, userId int) ([]model.Account, error)
	// QueryById 根据id查询账户
	QueryById(ctx context.Context, id int) (*model.Account, error)

	// ModifyBalance 修改账户余额
	ModifyBalance(ctx context.Context, id int, diff float64) error
}

type account struct {
}

func NewAccountStore() AccountStore {
	return &account{}
}

// Add 添加账户
func (a *account) Add(ctx context.Context, account *model.Account) error {
	db := getDBFromContext(ctx)

	insertSql := "insert into user_account(user_id, account_type, account_name, balance, init, icon, is_total, remark, " +
		"create_time, update_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(insertSql, account.UserId, account.AccountType, account.AccountName, account.Balance, account.Init, account.Icon,
		account.IsTotal, account.Remark, account.CreateTime, account.UpdateTime)
	if err != nil {
		return errors.Wrap(err, "account add store")
	}

	return nil
}

// Update 更新账户
func (a *account) Update(ctx context.Context, account *model.Account) error {
	db := getDBFromContext(ctx)
	updateSql := "update user_account set account_type=?,account_name=?,balance=?,icon=?,is_total=?,remark=?,update_time=? where id = ?"

	_, err := db.Exec(updateSql, account.AccountType, account.AccountName, account.Balance, account.Icon, account.IsTotal,
		account.Remark, account.UpdateTime, account.Id)
	if err != nil {
		return errors.Wrap(err, "account update store")
	}
	return nil
}

// Delete 删除账户
// 执行的是逻辑删除，将字段del置为1
func (a *account) Delete(ctx context.Context, id int) error {
	db := getDBFromContext(ctx)

	deleteSql := "update user_account set del = ? where id =?"
	_, err := db.Exec(deleteSql, constant.DelTrue, id)
	if err != nil {
		return errors.Wrap(err, "account delete store")
	}

	return nil
}

// QueryAllByUserId 根据用户id查询其所有账户
func (a *account) QueryAllByUserId(ctx context.Context, userId int) ([]model.Account, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from user_account where user_id = ? and del = ?"

	var accounts = []model.Account{}
	err := db.Select(&accounts, querySql, userId, constant.DelFalse)
	if err != nil {
		if err == sql.ErrNoRows {
			return accounts, nil
		}
		return nil, errors.Wrap(err, "account query all store")
	}

	return accounts, nil
}

// QueryById 根据id查询账户
func (a *account) QueryById(ctx context.Context, id int) (*model.Account, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from user_account where id = ? and del = ?"

	var accountModel model.Account
	err := db.Get(&accountModel, querySql, id, constant.DelFalse)
	if err != nil {
		if err == sql.ErrNoRows {
			return &accountModel, nil
		}
		return nil, errors.Wrap(err, "accountModel query all store")
	}

	return &accountModel, nil
}

// ModifyBalance 修改账户余额
func (a *account) ModifyBalance(ctx context.Context, id int, diff float64) error {
	db := getDBFromContext(ctx)
	updateSql := "update user_account set balance = balance + ? where id = ?"

	_, err := db.Exec(updateSql, diff, id)

	if err != nil {
		return errors.Wrap(err, "account modify balance store")
	}

	return nil
}
