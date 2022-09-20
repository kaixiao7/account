package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/constant"

	"github.com/pkg/errors"
)

type AssetFlowStore interface {
	// Add 添加资产流水
	Add(ctx context.Context, flow *model.AssetFlow) error
	// Update 更新资产流水
	Update(ctx context.Context, flow *model.AssetFlow) error
	// Delete 删除资产流水
	// 逻辑删除，将字段del置为1
	Delete(ctx context.Context, id int) error
	// DeleteByAssetId 根据资产删除其对应的流水
	// 逻辑删除，将字段del置为1
	DeleteByAssetId(ctx context.Context, assetId int) error

	// QueryByUserIdAndAssetId 根据用户id和资产id查询资产流水
	QueryByUserIdAndAssetId(ctx context.Context, userId, assetId int) ([]model.AssetFlow, error)
	// QueryByUserIdAndType 根据用户id和类型查询资产流水
	QueryByUserIdAndType(ctx context.Context, userId, assetFlowType int) ([]model.AssetFlow, error)
	// QueryById 根据id查询
	QueryById(ctx context.Context, id int) (*model.AssetFlow, error)
	// QueryReverseFlow 查询转入、转出反转流水
	QueryReverseFlow(ctx context.Context, assetId, targetAssetId, flowType int, time int64) (*model.AssetFlow, error)
}

type assetFlow struct {
}

func NewAssetFlowStore() AssetFlowStore {
	return &assetFlow{}
}

// Add 添加资产流水
func (af *assetFlow) Add(ctx context.Context, flow *model.AssetFlow) error {
	db := getDBFromContext(ctx)

	field := "user_id, asset_id, type, cost, record_time, remark, associate_name, create_time, update_time"
	values := "?,?,?,?,?,?,?,?,?"
	v := []any{flow.UserId, flow.AssetId, flow.Type, flow.Cost, flow.RecordTime, flow.Remark, flow.AssociateName, flow.CreateTime, flow.UpdateTime}
	if flow.CategoryId != nil {
		field = field + ", category_id"
		values = values + ", ?"
		v = append(v, flow.CategoryId)
	}
	if flow.TargetAssetId != nil {
		field = field + ", target_asset_id"
		values = values + ", ?"
		v = append(v, flow.TargetAssetId)
	}
	if flow.Finished != nil {
		field = field + ", finished"
		values = values + ", ?"
		v = append(v, flow.Finished)
	}

	sql := "insert into asset_flow(" + field + ") values(" + values + ")"
	_, err := db.Exec(sql, v...)
	if err != nil {
		return errors.Wrap(err, "asset flow add store")
	}

	return nil
}

// Update 更新资产流水
func (af *assetFlow) Update(ctx context.Context, flow *model.AssetFlow) error {
	db := getDBFromContext(ctx)

	sql := "update asset_flow set user_id=?,asset_id=?,type=?,cost=?,record_time=?,remark=?,associate_name=?,update_time=?"
	v := []any{flow.UserId, flow.AssetId, flow.Type, flow.Cost, flow.RecordTime, flow.Remark, flow.AssociateName, flow.UpdateTime}
	if flow.CategoryId != nil {
		sql = sql + ",category_id=?"
		v = append(v, flow.CategoryId)
	}
	if flow.TargetAssetId != nil {
		sql = sql + ",target_asset_id=?"
		v = append(v, flow.TargetAssetId)
	}
	if flow.Finished != nil {
		sql = sql + ",finished=?"
		v = append(v, flow.Finished)
	}
	sql = sql + " where id = ?"
	v = append(v, flow.Id)

	_, err := db.Exec(sql, v...)
	if err != nil {
		return errors.Wrap(err, "asset flow update store")
	}
	return nil
}

// Delete 删除资产流水
// 逻辑删除，将字段del置为1
func (af *assetFlow) Delete(ctx context.Context, id int) error {
	db := getDBFromContext(ctx)

	sql := "update asset_flow set del = ? where id = ?"
	_, err := db.Exec(sql, constant.DelTrue, id)
	if err != nil {
		return errors.Wrap(err, "delete asset flow store")
	}
	return nil
}

// DeleteByAssetId 根据资产删除其对应的流水
// 逻辑删除，将字段del置为1
func (af *assetFlow) DeleteByAssetId(ctx context.Context, assetId int) error {
	db := getDBFromContext(ctx)

	sql := "update asset_flow set del = ? where asset_id =?"
	_, err := db.Exec(sql, constant.DelTrue, assetId)
	if err != nil {
		return errors.Wrap(err, "delete asset flow store")
	}
	return nil
}

// QueryByUserIdAndAssetId 根据用户id和资产id查询资产流水
func (af *assetFlow) QueryByUserIdAndAssetId(ctx context.Context, userId, assetId int) ([]model.AssetFlow, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from asset_flow where user_id=? and asset_id=? and del=?"
	var assetFlows = []model.AssetFlow{}
	err := db.Select(&assetFlows, querySql, userId, assetId, constant.DelFalse)
	if err != nil {
		return nil, errors.Wrap(err, "query asset flow by userId and assetId store")
	}

	return assetFlows, nil
}

// QueryByUserIdAndType 根据用户id和类型查询资产流水
func (af *assetFlow) QueryByUserIdAndType(ctx context.Context, userId, assetFlowType int) ([]model.AssetFlow, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from asset_flow where user_id=? and type=? and del=?"
	var assetFlows = []model.AssetFlow{}
	err := db.Select(&assetFlows, querySql, userId, assetFlowType, constant.DelFalse)
	if err != nil {
		return nil, errors.Wrap(err, "query asset flow by userId and assetId store")
	}

	return assetFlows, nil
}

// QueryById 根据id查询
func (af *assetFlow) QueryById(ctx context.Context, id int) (*model.AssetFlow, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from asset_flow where id = ? and del =?"
	var assetFlow model.AssetFlow
	err := db.Get(&assetFlow, querySql, id, constant.DelFalse)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "asset flow query by id store")
	}

	return &assetFlow, nil
}

// QueryReverseFlow 查询转入、转出反转流水
func (af *assetFlow) QueryReverseFlow(ctx context.Context, assetId, targetAssetId, flowType int, time int64) (*model.AssetFlow, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from asset_flow where type=? and asset_id = ? and target_asset_id = ? and record_time=? and del=?"
	var assetFlow model.AssetFlow
	err := db.Get(&assetFlow, querySql, flowType, assetId, targetAssetId, time, constant.DelFalse)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "asset flow query by id store")
	}

	return &assetFlow, nil
}
