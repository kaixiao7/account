package service

import (
	"context"
	"time"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/errno"
)

type BorrowLendSrv interface {
	// QueryTotal 查询总借入借出
	QueryTotal(ctx context.Context, userId int64) (*model.BorrowLendTotal, error)
	// AddBorrowLendFlow 添加借入借出流水
	AddBorrowLendFlow(ctx context.Context, bf *model.BorrowLendFlow, userId int64) error
	// DeleteBorrowLendFlow 删除借入借出流水
	DeleteBorrowLendFlow(ctx context.Context, bfId, userId int64) error
	// QueryBorrowLendFlowList 查询借入借出流水列表
	QueryBorrowLendFlowList(ctx context.Context, accountFlowId, userId int64) ([]model.AccountFlow, error)
	// QueryBorrowLendList 查询借入借出列表
	QueryBorrowLendList(ctx context.Context, userId int64, borrowType int) ([]model.AccountFlow, error)
	// FinishedBorrowLend 结束债务
	FinishedBorrowLend(ctx context.Context, borrowId, userId int64) error
}

type borrowLendService struct {
	accountStore     store.AccountStore
	accountFlowStore store.AccountFlow
	accountFlowSrv   AccountFlowSrv
}

func NewBorrowLendSrv() BorrowLendSrv {
	return &borrowLendService{
		accountFlowStore: store.NewAccountFlowStore(),
		accountStore:     store.NewAccountStore(),
		accountFlowSrv:   NewAccountFlowSrv(),
	}
}

// QueryTotal 查询总借入借出
func (b *borrowLendService) QueryTotal(ctx context.Context, userId int64) (*model.BorrowLendTotal, error) {
	// 借入流水
	borrowFlows, err := b.accountFlowStore.QueryByUserIdAndType(ctx, userId, constant.AccountTypeBorrow)
	if err != nil {
		return nil, err
	}
	// 借出流水
	lendFlows, err := b.accountFlowStore.QueryByUserIdAndType(ctx, userId, constant.AccountTypeLend)
	if err != nil {
		return nil, err
	}

	var borrow float64
	var lend float64
	for _, flow := range borrowFlows {
		if *flow.Finished == 1 {
			continue
		}

		borrow = borrow + flow.Cost
		// 已还记录
		stills, err := b.accountFlowStore.QueryByBorrowLendId(ctx, flow.Id)
		if err != nil {
			return nil, err
		}
		for _, still := range stills {
			borrow = borrow - still.Cost
		}
	}

	for _, flow := range lendFlows {
		if *flow.Finished == 1 {
			continue
		}

		lend = lend + flow.Cost
		// 已还记录
		stills, err := b.accountFlowStore.QueryByBorrowLendId(ctx, flow.Id)
		if err != nil {
			return nil, err
		}
		for _, still := range stills {
			lend = lend - still.Cost
		}
	}

	return &model.BorrowLendTotal{
		Borrow: borrow,
		Lend:   lend,
	}, nil
}

// AddBorrowLendFlow 添加借入借出流水
func (b *borrowLendService) AddBorrowLendFlow(ctx context.Context, bf *model.BorrowLendFlow, userId int64) error {

	if bf.Type != constant.AccountTypeStill && bf.Type != constant.AccountTypeHarvest {
		return errno.New(errno.ErrIllegalOperate)
	}

	flow, err := b.checkAccountFlow(ctx, bf.BorrowId, userId)
	if err != nil {
		return err
	}
	// 校验类型，还款 -> 借入， 收款 -> 借出
	if (bf.Type == constant.AccountTypeStill && flow.Type != constant.AccountTypeBorrow) ||
		(bf.Type == constant.AccountTypeHarvest && flow.Type != constant.AccountTypeLend) {
		return errno.New(errno.ErrValidation)
	}

	af := model.AccountFlow{
		UserId:       userId,
		AccountId:    bf.AccountId,
		Type:         bf.Type,
		Cost:         bf.Cost,
		RecordTime:   bf.RecordTime,
		Remark:       bf.Remark,
		BorrowLendId: &bf.BorrowId,
		CreateTime:   time.Now().Unix(),
		UpdateTime:   time.Now().Unix(),
	}

	return b.accountFlowSrv.Add(ctx, &af)
}

// DeleteBorrowLendFlow 删除借入借出流水(还款、收款)
func (b *borrowLendService) DeleteBorrowLendFlow(ctx context.Context, bfId, userId int64) error {
	return b.accountFlowSrv.Delete(ctx, bfId, userId)
}

// QueryBorrowLendFlowList 查询借入借出流水列表(还款、收款)
func (b *borrowLendService) QueryBorrowLendFlowList(ctx context.Context, accountFlowId, userId int64) ([]model.AccountFlow, error) {
	if _, err := b.checkAccountFlow(ctx, accountFlowId, userId); err != nil {
		return nil, err
	}

	return b.accountFlowStore.QueryByBorrowLendId(ctx, accountFlowId)
}

// QueryBorrowLendList 查询借入借出列表
func (b *borrowLendService) QueryBorrowLendList(ctx context.Context, userId int64, borrowType int) ([]model.AccountFlow, error) {
	if borrowType != constant.AccountTypeBorrow && borrowType != constant.AccountTypeLend {
		return nil, errno.New(errno.ErrValidation)
	}
	return b.accountFlowStore.QueryByUserIdAndType(ctx, userId, borrowType)
}

// FinishedBorrowLend 结束债务
func (b *borrowLendService) FinishedBorrowLend(ctx context.Context, borrowId, userId int64) error {
	if _, err := b.checkAccountFlow(ctx, borrowId, userId); err != nil {
		return err
	}
	return b.accountFlowStore.FinishedBorrow(ctx, borrowId)
}

func (b *borrowLendService) checkAccountFlow(ctx context.Context, accountFlowId, userId int64) (*model.AccountFlow, error) {
	accountFlow, err := b.accountFlowStore.QueryById(ctx, accountFlowId)
	if err != nil {
		return nil, err
	}
	if accountFlow == nil {
		return nil, errno.New(errno.ErrAccountFlowNotFound)
	}

	if accountFlow.UserId != userId {
		return nil, errno.New(errno.ErrIllegalOperate)
	}

	return accountFlow, err
}
