package service

import (
	"context"
	"time"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/errno"
)

type AccountSrv interface {
	// Add 添加资产账户
	Add(ctx context.Context, account *model.Account) error
	// Update 修改资产账户
	Update(ctx context.Context, account *model.Account) error
	// Delete 删除资产账户
	Delete(ctx context.Context, accountId, userId int) error
	// QueryByUserId 根据用户id查询所有资产账户
	QueryByUserId(ctx context.Context, userId int) ([]model.Account, error)
}

type accountService struct {
	accountStore     store.AccountStore
	accountFlowStore store.AccountFlow
}

func NewAccountSrv() AccountSrv {
	return &accountService{
		accountStore:     store.NewAccountStore(),
		accountFlowStore: store.NewAccountFlowStore(),
	}
}

// Add 添加资产账户
func (a *accountService) Add(ctx context.Context, account *model.Account) error {
	account.Init = account.Balance
	return a.accountStore.Add(ctx, account)
}

// Update 修改资产账户
func (a *accountService) Update(ctx context.Context, account *model.Account) error {
	accountBefore, err := a.accountStore.QueryById(ctx, account.Id)
	if err != nil {
		return err
	}

	// 非法操作
	if accountBefore.UserId != account.UserId {
		return errno.New(errno.ErrIllegalOperate)
	}

	now := time.Now().Unix()
	account.UpdateTime = now

	return WithTransaction(ctx, func(ctx context.Context) error {
		if e := a.accountStore.Update(ctx, account); e != nil {
			return e
		}

		// 资产账户金额不等，则需要添加一条流水
		diff := account.Balance - accountBefore.Balance
		if diff != 0 {
			accountFlow := model.AccountFlow{
				UserId:     account.UserId,
				AccountId:  account.Id,
				Type:       constant.AccountTypeModify,
				Cost:       diff,
				RecordTime: now,
				Remark:     "",
				CreateTime: now,
				UpdateTime: now,
			}

			if e := a.accountFlowStore.Add(ctx, &accountFlow); e != nil {
				return e
			}
		}
		return nil
	})
}

// Delete 删除资产账户
func (a *accountService) Delete(ctx context.Context, accountId, userId int) error {
	account, err := a.accountStore.QueryById(ctx, accountId)
	if err != nil {
		return err
	}

	// 非法操作
	if account.UserId != userId {
		return errno.New(errno.ErrIllegalOperate)
	}

	return WithTransaction(ctx, func(ctx context.Context) error {
		if e := a.accountStore.Delete(ctx, accountId); e != nil {
			return e
		}

		if e := a.accountFlowStore.DeleteByAccountId(ctx, accountId); e != nil {
			return e
		}

		return nil
	})
}

// QueryByUserId 根据用户id查询所有资产账户
func (a *accountService) QueryByUserId(ctx context.Context, userId int) ([]model.Account, error) {
	return a.accountStore.QueryAllByUserId(ctx, userId)
}
