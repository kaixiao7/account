package service

import (
	"context"
	"time"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/errno"
)

type AssetSrv interface {
	// Add 添加资产账户
	Add(ctx context.Context, asset *model.Asset) error
	// Update 修改资产账户
	Update(ctx context.Context, asset *model.Asset) error
	// Delete 删除资产账户
	Delete(ctx context.Context, assetId, userId int) error
	// QueryByUserId 根据用户id查询所有资产账户
	QueryByUserId(ctx context.Context, userId int) ([]model.Asset, error)
}

type assetService struct {
	assetStore     store.AssetStore
	assetFlowStore store.AssetFlowStore
}

func NewAssetSrv() AssetSrv {
	return &assetService{
		assetStore:     store.NewAssetStore(),
		assetFlowStore: store.NewAssetFlowStore(),
	}
}

// Add 添加资产账户
func (a *assetService) Add(ctx context.Context, asset *model.Asset) error {
	asset.Init = asset.Balance
	return a.assetStore.Add(ctx, asset)
}

// Update 修改资产账户
func (a *assetService) Update(ctx context.Context, asset *model.Asset) error {
	assetBefore, err := a.assetStore.QueryById(ctx, asset.Id)
	if err != nil {
		return err
	}

	// 非法操作
	if assetBefore.UserId != asset.UserId {
		return errno.New(errno.ErrIllegalOperate)
	}

	now := time.Now().Unix()
	asset.UpdateTime = now

	return WithTransaction(ctx, func(ctx context.Context) error {
		if e := a.assetStore.Update(ctx, asset); e != nil {
			return e
		}

		// 资产账户金额不等，则需要添加一条流水
		diff := asset.Balance - assetBefore.Balance
		if diff != 0 {
			assetFlow := model.AssetFlow{
				UserId:     asset.UserId,
				AssetId:    asset.Id,
				Type:       constant.AssetTypeModify,
				Cost:       diff,
				RecordTime: now,
				Remark:     "",
				CreateTime: now,
				UpdateTime: now,
			}

			if e := a.assetFlowStore.Add(ctx, &assetFlow); e != nil {
				return e
			}
		}
		return nil
	})
}

// Delete 删除资产账户
func (a *assetService) Delete(ctx context.Context, assetId, userId int) error {
	asset, err := a.assetStore.QueryById(ctx, assetId)
	if err != nil {
		return err
	}

	// 非法操作
	if asset.UserId != userId {
		return errno.New(errno.ErrIllegalOperate)
	}

	return WithTransaction(ctx, func(ctx context.Context) error {
		if e := a.assetStore.Delete(ctx, assetId); e != nil {
			return e
		}

		if e := a.assetFlowStore.DeleteByAssetId(ctx, assetId); e != nil {
			return e
		}

		return nil
	})
}

// QueryByUserId 根据用户id查询所有资产账户
func (a *assetService) QueryByUserId(ctx context.Context, userId int) ([]model.Asset, error) {
	return a.assetStore.QueryAllByUserId(ctx, userId)
}
