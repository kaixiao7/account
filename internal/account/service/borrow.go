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
	// DeleteBorrowFlow 删除借入借出流水
	DeleteBorrowFlow(ctx context.Context, bfId, userId int) error
	// QueryBorrowFlowList 查询借入借出流水列表
	QueryBorrowFlowList(ctx context.Context, assetFlowId, userId int) ([]model.AssetFlow, error)
	// QueryBorrowList 查询借入借出列表
	QueryBorrowList(ctx context.Context, userId, borrowType int) ([]model.AssetFlow, error)
	// FinishedBorrow 结束债务
	FinishedBorrow(ctx context.Context, borrowId, userId int) error
}

type borrowService struct {
	assetStore     store.AssetStore
	assetFlowStore store.AssetFlowStore
	assetFlowSrv   AssetFlowSrv
}

func NewBorrowSrv() BorrowSrv {
	return &borrowService{
		assetFlowStore: store.NewAssetFlowStore(),
		assetStore:     store.NewAssetStore(),
		assetFlowSrv:   NewAssertFlowSrv(),
	}
}

// QueryTotal 查询总借入借出
func (b *borrowService) QueryTotal(ctx context.Context, userId int) (*model.BorrowTotal, error) {
	// 借入流水
	borrowInFlows, err := b.assetFlowStore.QueryByUserIdAndType(ctx, userId, constant.AssetTypeBorrowIn)
	if err != nil {
		return nil, err
	}
	// 借出流水
	borrowOutFlows, err := b.assetFlowStore.QueryByUserIdAndType(ctx, userId, constant.AssetTypeBorrowOut)
	if err != nil {
		return nil, err
	}
	// 还款流水
	// stills, err := b.assetFlowStore.QueryByUserIdAndType(ctx, userId, constant.AssetTypeStill)
	// if err != nil {
	// 	return nil, err
	// }
	// // 收款流水
	// harvests, err := b.assetFlowStore.QueryByUserIdAndType(ctx, userId, constant.AssetTypeHarvest)
	// if err != nil {
	// 	return nil, err
	// }

	var borrowIn float64
	var borrowOut float64
	for _, flow := range borrowInFlows {
		if *flow.Finished == 1 {
			continue
		}

		borrowIn = borrowIn + flow.Cost
		stills, err := b.assetFlowStore.QueryByUserIdAndBorrowId(ctx, userId, flow.Id)
		if err != nil {
			return nil, err
		}
		for _, still := range stills {
			borrowIn = borrowIn - still.Cost
		}
	}

	for _, flow := range borrowOutFlows {
		if *flow.Finished == 1 {
			continue
		}

		borrowOut = borrowOut + flow.Cost
		stills, err := b.assetFlowStore.QueryByUserIdAndBorrowId(ctx, userId, flow.Id)
		if err != nil {
			return nil, err
		}
		for _, still := range stills {
			borrowOut = borrowOut - still.Cost
		}
	}

	return &model.BorrowTotal{
		BorrowIn:  borrowIn,
		BorrowOut: borrowOut,
	}, nil
}

// AddBorrowFlow 添加借入借出流水
func (b *borrowService) AddBorrowFlow(ctx context.Context, bf *model.BorrowFlow, userId int) error {

	if bf.Type != constant.AssetTypeStill && bf.Type != constant.AssetTypeHarvest {
		return errno.New(errno.ErrIllegalOperate)
	}

	flow, err := b.checkAssetFlow(ctx, bf.BorrowId, userId)
	if err != nil {
		return err
	}
	// 校验类型，还款 -> 借入， 收款 -> 借出
	if (bf.Type == constant.AssetTypeStill && flow.Type != constant.AssetTypeBorrowIn) ||
		(bf.Type == constant.AssetTypeHarvest && flow.Type != constant.AssetTypeBorrowOut) {
		return errno.New(errno.ErrValidation)
	}

	af := model.AssetFlow{
		UserId:     userId,
		AssetId:    bf.AssetId,
		Type:       bf.Type,
		Cost:       bf.Cost,
		RecordTime: bf.RecordTime,
		Remark:     bf.Remark,
		BorrowId:   &bf.BorrowId,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	}

	return b.assetFlowSrv.Add(ctx, &af)
}

// DeleteBorrowFlow 删除借入借出流水(还款、收款)
func (b *borrowService) DeleteBorrowFlow(ctx context.Context, bfId, userId int) error {
	return b.assetFlowSrv.Delete(ctx, bfId, userId)
}

// QueryBorrowFlowList 查询借入借出流水列表(还款、收款)
func (b *borrowService) QueryBorrowFlowList(ctx context.Context, assetFlowId, userId int) ([]model.AssetFlow, error) {
	if _, err := b.checkAssetFlow(ctx, assetFlowId, userId); err != nil {
		return nil, err
	}

	return b.assetFlowStore.QueryByUserIdAndBorrowId(ctx, userId, assetFlowId)
}

// QueryBorrowList 查询借入借出列表
func (b *borrowService) QueryBorrowList(ctx context.Context, userId, borrowType int) ([]model.AssetFlow, error) {
	if borrowType != constant.AssetTypeBorrowIn && borrowType != constant.AssetTypeBorrowOut {
		return nil, errno.New(errno.ErrValidation)
	}
	return b.assetFlowStore.QueryByUserIdAndType(ctx, userId, borrowType)
}

// FinishedBorrow 结束债务
func (b *borrowService) FinishedBorrow(ctx context.Context, borrowId, userId int) error {
	if _, err := b.checkAssetFlow(ctx, borrowId, userId); err != nil {
		return err
	}
	return b.assetFlowStore.FinishedBorrow(ctx, borrowId)
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
