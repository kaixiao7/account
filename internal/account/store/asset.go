package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/constant"

	"github.com/pkg/errors"
)

// AssetStore 资产账户
type AssetStore interface {
	// Add 添加资产账户
	Add(ctx context.Context, asset *model.Asset) error
	// Update 更新资产账户
	Update(ctx context.Context, asset *model.Asset) error
	// Delete 删除资产账户
	// 执行的是逻辑删除，将字段del置为1
	Delete(ctx context.Context, id int) error

	// QueryAllByUserId 根据用户id查询其所有资产账户
	QueryAllByUserId(ctx context.Context, userId int) ([]model.Asset, error)
	// QueryById 根据id查询资产账户
	QueryById(ctx context.Context, id int) (*model.Asset, error)

	// ModifyBalance 修改账户余额
	ModifyBalance(ctx context.Context, id int, diff float64) error
}

type asset struct {
}

func NewAssetStore() AssetStore {
	return &asset{}
}

// Add 添加资产账户
func (a *asset) Add(ctx context.Context, asset *model.Asset) error {
	db := getDBFromContext(ctx)

	sql := "insert into asset_account(user_id, asset_type, asset_name, balance, init, icon, is_total, remark, create_time, update_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(sql, asset.UserId, asset.AssetType, asset.AssetName, asset.Balance, asset.Init, asset.Icon,
		asset.IsTotal, asset.Remark, asset.CreateTime, asset.UpdateTime)
	if err != nil {
		return errors.Wrap(err, "asset add store")
	}

	return nil
}

// Update 更新资产账户
func (a *asset) Update(ctx context.Context, asset *model.Asset) error {
	db := getDBFromContext(ctx)
	sql := "update asset_account set asset_type=?,asset_name=?,balance=?,icon=?,is_total=?,remark=?,update_time=? where id = ?"

	_, err := db.Exec(sql, asset.AssetType, asset.AssetName, asset.Balance, asset.Icon, asset.IsTotal,
		asset.Remark, asset.UpdateTime, asset.Id)
	if err != nil {
		return errors.Wrap(err, "asset update store")
	}
	return nil
}

// Delete 删除资产账户
// 执行的是逻辑删除，将字段del置为1
func (a *asset) Delete(ctx context.Context, id int) error {
	db := getDBFromContext(ctx)

	// sql := "delete from asset_account where id = ?"
	sql := "update asset_account set del = ? where id =?"
	_, err := db.Exec(sql, constant.DelTrue, id)
	if err != nil {
		return errors.Wrap(err, "asset delete store")
	}

	return nil
}

// QueryAllByUserId 根据用户id查询其所有资产账户
func (a *asset) QueryAllByUserId(ctx context.Context, userId int) ([]model.Asset, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from asset_account where user_id = ? and del = ?"

	var assets = []model.Asset{}
	err := db.Select(&assets, querySql, userId, constant.DelFalse)
	if err != nil {
		if err == sql.ErrNoRows {
			return assets, nil
		}
		return nil, errors.Wrap(err, "asset query all store")
	}

	return assets, nil
}

// QueryById 根据id查询资产账户
func (a *asset) QueryById(ctx context.Context, id int) (*model.Asset, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from asset_account where id = ? and del = ?"

	var asset model.Asset
	err := db.Get(&asset, querySql, id, constant.DelFalse)
	if err != nil {
		if err == sql.ErrNoRows {
			return &asset, nil
		}
		return nil, errors.Wrap(err, "asset query all store")
	}

	return &asset, nil
}

// ModifyBalance 修改账户余额
func (a *asset) ModifyBalance(ctx context.Context, id int, diff float64) error {
	db := getDBFromContext(ctx)
	sql := "update asset_account set balance = balance + ? where id = ?"

	_, err := db.Exec(sql, diff, id)

	if err != nil {
		return errors.Wrap(err, "asset modify balance store")
	}

	return nil
}
