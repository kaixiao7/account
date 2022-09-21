package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"
)

type BorrowFlowStore interface {
	// Add 添加借入借出流水
	Add(ctx context.Context, bf *model.BorrowFlow) error
	// Update 修改借入借出流水
	Update(ctx context.Context, bf *model.BorrowFlow) error
	// Delete 删除借入借出流水
	Delete(ctx context.Context, id int) error
	// QueryByAssetFlowId 根据资产流水id查询借入借出流水
	QueryByAssetFlowId(ctx context.Context, assetFlowId int) ([]model.BorrowFlow, error)
	// QueryById 根据id查询
	QueryById(ctx context.Context, id int) (*model.BorrowFlow, error)
}

type borrowFlow struct {
}

func NewBorrowFlowStore() BorrowFlowStore {
	return &borrowFlow{}
}

// Add 添加借入借出流水
func (b *borrowFlow) Add(ctx context.Context, bf *model.BorrowFlow) error {
	db := getDBFromContext(ctx)

	sql := "insert into asset_borrow_flow(asset_flow_id, asset_id, cost, record_time, type, remark, create_time, update_time) values (?,?,?,?,?,?,?,?)"
	_, err := db.Exec(sql, bf.AssetFlowId, bf.AssetId, bf.Cost, bf.RecordTime, bf.Type, bf.Remark, bf.CreateTime, bf.UpdateTime)
	if err != nil {
		return err
	}

	return nil
}

// Update 修改借入借出流水
func (b *borrowFlow) Update(ctx context.Context, bf *model.BorrowFlow) error {
	db := getDBFromContext(ctx)

	sql := "update asset_borrow_flow set asset_flow_id=?,asset_id=?,cost=?,record_time=?,type=?,remark=?,update_time=? where id = ?"
	_, err := db.Exec(sql, bf.AssetFlowId, bf.AssetId, bf.Cost, bf.RecordTime, bf.Type, bf.Remark, bf.UpdateTime, bf.Id)
	if err != nil {
		return err
	}

	return nil
}

// Delete 删除借入借出流水
func (b *borrowFlow) Delete(ctx context.Context, id int) error {
	db := getDBFromContext(ctx)

	sql := "delete from asset_borrow_flow where id=?"
	_, err := db.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}

// QueryByAssetFlowId 根据资产流水id查询借入借出流水
func (b *borrowFlow) QueryByAssetFlowId(ctx context.Context, assetFlowId int) ([]model.BorrowFlow, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from asset_borrow_flow where asset_flow_id=?"
	var bfs = []model.BorrowFlow{}
	err := db.Select(&bfs, querySql, assetFlowId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return bfs, nil
}

// QueryById 根据id查询
func (b *borrowFlow) QueryById(ctx context.Context, id int) (*model.BorrowFlow, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from asset_borrow_flow where id=?"
	var bf model.BorrowFlow
	err := db.Get(&bf, querySql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &bf, nil
}
