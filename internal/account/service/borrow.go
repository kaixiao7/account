package service

import (
	"context"
	"time"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/errno"
)

type BorrowSrv interface {
	// QueryTotal 查询总借入借出
	QueryTotal(ctx context.Context, userId int) (*model.BorrowTotal, error)
	// AddBorrowFlow 添加借入借出流水
	AddBorrowFlow(ctx context.Context, bf *model.BorrowFlow, userId int) error
	// UpdateBorrowFlow 更新借入借出流水
	UpdateBorrowFlow(ctx context.Context, bf *model.BorrowFlow, userId int) error
	// DeleteBorrowFlow 删除借入借出流水
	DeleteBorrowFlow(ctx context.Context, bfId, userId int) error
	// QueryBorrowFlowList 查询借入借出流水列表
	QueryBorrowFlowList(ctx context.Context, assetFlowId, userId int) ([]model.BorrowFlow, error)
	// QueryBorrowList 查询借入借出列表
	QueryBorrowList(ctx context.Context, userId, borrowType int) ([]model.AssetFlow, error)
}

type borrowService struct {
	borrowFlowStore store.BorrowFlowStore
	assetStore      store.AssetStore
	assetFlowStore  store.AssetFlowStore
}

func NewBorrowSrv() BorrowSrv {
	return &borrowService{
		borrowFlowStore: store.NewBorrowFlowStore(),
		assetFlowStore:  store.NewAssetFlowStore(),
		assetStore:      store.NewAssetStore(),
	}
}

// QueryTotal 查询总借入借出
func (b *borrowService) QueryTotal(ctx context.Context, userId int) (*model.BorrowTotal, error) {
	borrowInFlows, err := b.assetFlowStore.QueryByUserIdAndType(ctx, userId, constant.AssetTypeBorrowIn)
	if err != nil {
		return nil, err
	}
	borrowOutFlows, err := b.assetFlowStore.QueryByUserIdAndType(ctx, userId, constant.AssetTypeBorrowOut)
	if err != nil {
		return nil, err
	}
	var borrowIn float64
	var borrowOut float64
	for _, flow := range borrowInFlows {
		borrowIn = borrowIn + flow.Cost
	}
	for _, flow := range borrowOutFlows {
		borrowOut = borrowOut + flow.Cost
	}

	return &model.BorrowTotal{
		BorrowIn:  borrowIn,
		BorrowOut: borrowOut,
	}, nil
}

// AddBorrowFlow 添加借入借出流水
func (b *borrowService) AddBorrowFlow(ctx context.Context, bf *model.BorrowFlow, userId int) error {
	assetFlow, err := b.checkAssetFlow(ctx, bf.AssetFlowId, userId)
	if err != nil {
		return err
	}
	if assetFlow.Type != bf.Type {
		return errno.New(errno.ErrIllegalOperate)
	}

	if _, err := b.checkAsset(ctx, bf.AssetId, userId); err != nil {
		return err
	}

	cost := bf.Cost
	if bf.Type == constant.AssetTypeBorrowIn {
		// 类型为借入，那么流水表示还款，需要取反
		cost = -cost
	}

	now := time.Now().Unix()
	bf.CreateTime = now
	bf.UpdateTime = now

	return WithTransaction(ctx, func(ctx context.Context) error {
		// 修改账户余额
		if err := b.assetStore.ModifyBalance(ctx, bf.AssetId, cost); err != nil {
			return err
		}
		// 添加借入借出流水
		if err := b.borrowFlowStore.Add(ctx, bf); err != nil {
			return err
		}
		// 更新资产流水中的借入借出的金额
		assetFlow.Cost = assetFlow.Cost - bf.Cost
		if err := b.assetFlowStore.Update(ctx, assetFlow); err != nil {
			return err
		}

		return nil
	})
}

// UpdateBorrowFlow 更新借入借出流水
func (b *borrowService) UpdateBorrowFlow(ctx context.Context, bf *model.BorrowFlow, userId int) error {
	assetFlow, err := b.checkAssetFlow(ctx, bf.AssetFlowId, userId)
	if err != nil {
		return err
	}
	if assetFlow.Type != bf.Type {
		return errno.New(errno.ErrIllegalOperate)
	}
	if _, err := b.checkAsset(ctx, bf.AssetId, userId); err != nil {
		return err
	}
	// 所属的资产流水不能修改
	if assetFlow.Id != bf.AssetFlowId {
		return errno.New(errno.ErrIllegalOperate)
	}

	borrowFlowBefore, err := b.checkBorrowFlow(ctx, bf.Id)
	if err != nil {
		return err
	}
	bf.UpdateTime = time.Now().Unix()

	return WithTransaction(ctx, func(ctx context.Context) error {
		diff := borrowFlowBefore.Cost - bf.Cost
		// 借出取反
		if bf.Type == constant.AssetTypeBorrowOut {
			diff = -diff
		}

		// 更新借入借出流水
		if err := b.borrowFlowStore.Update(ctx, bf); err != nil {
			return err
		}
		// 更新资产流水中的借入借出的金额
		assetFlow.Cost = assetFlow.Cost + diff
		if err := b.assetFlowStore.Update(ctx, assetFlow); err != nil {
			return err
		}

		if bf.AssetId == borrowFlowBefore.AssetId {
			// 修改账户余额
			if err := b.assetStore.ModifyBalance(ctx, bf.AssetId, diff); err != nil {
				return err
			}
		} else {
			// 修改前面的账户金额
			cost := borrowFlowBefore.Cost
			// 借出取反
			if bf.Type == constant.AssetTypeBorrowOut {
				cost = -cost
			}
			if err := b.assetStore.ModifyBalance(ctx, borrowFlowBefore.AssetId, cost); err != nil {
				return err
			}

			// 修改当前账户的金额
			cost = bf.Cost
			// 借入取反
			if bf.Type == constant.AssetTypeBorrowIn {
				cost = -cost
			}
			if err := b.assetStore.ModifyBalance(ctx, bf.AssetId, cost); err != nil {
				return err
			}
		}

		return nil
	})
}

// DeleteBorrowFlow 删除借入借出流水
func (b *borrowService) DeleteBorrowFlow(ctx context.Context, bfId, userId int) error {
	bf, err := b.checkBorrowFlow(ctx, bfId)
	if err != nil {
		return err
	}
	assetFlow, err := b.checkAssetFlow(ctx, bf.AssetFlowId, userId)
	if err != nil {
		return err
	}
	if assetFlow.Type != bf.Type {
		return errno.New(errno.ErrIllegalOperate)
	}

	return WithTransaction(ctx, func(ctx context.Context) error {
		// 更新资产流水中的借入借出的金额
		assetFlow.Cost = assetFlow.Cost + bf.Cost
		if err := b.assetFlowStore.Update(ctx, assetFlow); err != nil {
			return err
		}

		// 修改账户余额
		cost := bf.Cost
		if bf.Type == constant.AssetTypeBorrowOut {
			// 借出取反
			cost = -cost
		}
		if err := b.assetStore.ModifyBalance(ctx, bf.AssetId, cost); err != nil {
			return err
		}

		// 删除
		if err := b.borrowFlowStore.Delete(ctx, bf.Id); err != nil {
			return err
		}

		return nil
	})

}

// QueryBorrowFlowList 查询借入借出流水列表
func (b *borrowService) QueryBorrowFlowList(ctx context.Context, assetFlowId, userId int) ([]model.BorrowFlow, error) {
	if _, err := b.checkAssetFlow(ctx, assetFlowId, userId); err != nil {
		return nil, err
	}

	return b.borrowFlowStore.QueryByAssetFlowId(ctx, assetFlowId)
}

// QueryBorrowList 查询借入借出列表
func (b *borrowService) QueryBorrowList(ctx context.Context, userId, borrowType int) ([]model.AssetFlow, error) {
	return b.assetFlowStore.QueryByUserIdAndType(ctx, userId, borrowType)
}

func (b *borrowService) checkAssetFlow(ctx context.Context, assetFlowId, userId int) (*model.AssetFlow, error) {
	assetFlow, err := b.assetFlowStore.QueryById(ctx, assetFlowId)
	if err != nil {
		return nil, err
	}
	if assetFlow == nil {
		return nil, errno.New(errno.ErrAssetFlowNotFound)
	}

	if assetFlow.UserId != userId {
		return nil, errno.New(errno.ErrIllegalOperate)
	}

	return assetFlow, err
}

func (b *borrowService) checkAsset(ctx context.Context, assetId, userId int) (*model.Asset, error) {
	asset, err := b.assetStore.QueryById(ctx, assetId)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errno.New(errno.ErrAssetNotFound)
	}
	if asset.UserId != userId {
		return nil, errno.New(errno.ErrIllegalOperate)
	}

	return asset, err
}

func (b *borrowService) checkBorrowFlow(ctx context.Context, borrowFlowId int) (*model.BorrowFlow, error) {
	borrowFlow, err := b.borrowFlowStore.QueryById(ctx, borrowFlowId)
	if err != nil {
		return nil, err
	}
	if borrowFlow == nil {
		return nil, errno.New(errno.ErrBorrowFlowNotFound)
	}

	return borrowFlow, err
}
