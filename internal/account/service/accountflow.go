package service

import (
	"context"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/errno"
)

type AccountFlowSrv interface {
	// Add 添加流水
	Add(ctx context.Context, accountFlow *model.AccountFlow) error
	// Update 修改流水
	// 仅修改基本信息，不会修改类型及账户
	Update(ctx context.Context, accountFlow *model.AccountFlow) error
	// Delete 删除流水
	Delete(ctx context.Context, accountFlowId, userId int) error

	// QueryByAccountId 根据账户id查询其下的所有流水信息
	QueryByAccountId(ctx context.Context, accountId, userId int) ([]model.AccountFlow, error)
}

type accountFlowService struct {
	accountStore     store.AccountStore
	accountFlowStore store.AccountFlow
}

func NewAccountFlowSrv() AccountFlowSrv {
	return &accountFlowService{
		accountStore:     store.NewAccountStore(),
		accountFlowStore: store.NewAccountFlowStore(),
	}
}

// Add 添加流水
func (af *accountFlowService) Add(ctx context.Context, accountFlow *model.AccountFlow) error {
	if err := af.saveCheck(ctx, accountFlow); err != nil {
		return err
	}

	return WithTransaction(ctx, func(ctx context.Context) error {
		diff := accountFlow.Cost
		// 转出、借出、还款 将金额变为负数，因为修改账户余额中的sql计算方式为加法
		if accountFlow.Type == constant.AccountTypeTransferOut || accountFlow.Type == constant.AccountTypeLend ||
			accountFlow.Type == constant.AccountTypeStill {
			diff = -diff
		}
		// 修改账户余额
		if err := af.accountStore.ModifyBalance(ctx, accountFlow.AccountId, diff); err != nil {
			return err
		}
		if accountFlow.TargetAccountId != nil {
			// 修改目标账户余额
			if err := af.accountStore.ModifyBalance(ctx, *accountFlow.TargetAccountId, -diff); err != nil {
				return err
			}
		}
		// 插入流水记录
		if err := af.accountFlowStore.Add(ctx, accountFlow); err != nil {
			return err
		}

		return nil
	})
}

// Update 修改流水
// 仅修改基本信息，不会修改类型及账户
// 其实流水不应该被修改，只能增加、删除
func (af *accountFlowService) Update(ctx context.Context, accountFlow *model.AccountFlow) error {
	// 仅修改基本信息，不会修改类型及账户
	accountFlowBefore, err := af.checkAccountFlow(ctx, accountFlow.Id, accountFlow.UserId)
	if err != nil {
		return err
	}

	if accountFlowBefore.Type != accountFlow.Type || accountFlowBefore.AccountId != accountFlow.AccountId {
		return errno.New(errno.ErrIllegalOperate)
	}

	if err := af.saveCheck(ctx, accountFlow); err != nil {
		return err
	}

	// 前后差值
	diff := accountFlowBefore.Cost - accountFlow.Cost
	return WithTransaction(ctx, func(ctx context.Context) error {
		// 转入、借入、收款 将金额变为负数
		if accountFlow.Type == constant.AccountTypeTransferIn || accountFlow.Type == constant.AccountTypeBorrow ||
			accountFlow.Type == constant.AccountTypeHarvest {
			diff = -diff
		}
		// 修改账户金额
		if err := af.accountStore.ModifyBalance(ctx, accountFlow.AccountId, diff); err != nil {
			return err
		}
		// 修改目标账户金额
		if err := af.accountStore.ModifyBalance(ctx, *accountFlow.TargetAccountId, -diff); err != nil {
			return err
		}
		// 修改流水信息
		if err := af.accountFlowStore.Update(ctx, accountFlow); err != nil {
			return err
		}

		return nil
	})
}

// Delete 删除流水
func (af *accountFlowService) Delete(ctx context.Context, accountFlowId, userId int) error {
	accountFlow, err := af.checkAccountFlow(ctx, accountFlowId, userId)
	if err != nil {
		return err
	}
	// 修改账户余额类型不允许删除
	if accountFlow.Type == constant.AccountTypeModify {
		return errno.New(errno.ErrIllegalOperate)
	}

	return WithTransaction(ctx, func(ctx context.Context) error {
		// 账户金额恢复
		if err := af.moneyRegain(ctx, accountFlow); err != nil {
			return err
		}

		// 删除流水记录
		if err := af.accountFlowStore.Delete(ctx, accountFlowId); err != nil {
			return err
		}

		return nil
	})
}

// 账户余额恢复
func (af *accountFlowService) moneyRegain(ctx context.Context, accountFlow *model.AccountFlow) error {
	cost := accountFlow.Cost
	// 转入、借入、收款 将金额变为负数
	if accountFlow.Type == constant.AccountTypeTransferIn || accountFlow.Type == constant.AccountTypeBorrow ||
		accountFlow.Type == constant.AccountTypeHarvest {
		cost = -cost
	}
	// 账户余额恢复
	if err := af.accountStore.ModifyBalance(ctx, accountFlow.AccountId, cost); err != nil {
		return err
	}
	if accountFlow.TargetAccountId != nil {
		// 目标账户余额恢复
		if err := af.accountStore.ModifyBalance(ctx, *accountFlow.TargetAccountId, -cost); err != nil {
			return err
		}
	}

	return nil
}

func (af *accountFlowService) checkAccount(ctx context.Context, accountId, userId int) (*model.Account, error) {
	account, err := af.accountStore.QueryById(ctx, accountId)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errno.New(errno.ErrAccountNotFound)
	}

	if account.UserId != userId {
		return nil, errno.New(errno.ErrIllegalOperate)
	}

	return account, nil
}

func (af *accountFlowService) checkAccountFlow(ctx context.Context, accountFlowId, userId int) (*model.AccountFlow, error) {
	accountFlow, err := af.accountFlowStore.QueryById(ctx, accountFlowId)
	if err != nil {
		return nil, err
	}

	if accountFlow == nil {
		return nil, errno.New(errno.ErrAccountFlowNotFound)
	}

	if accountFlow.UserId != userId {
		return nil, errno.New(errno.ErrIllegalOperate)
	}

	return accountFlow, nil
}

// 插入、更新操作的前置校验
func (af *accountFlowService) saveCheck(ctx context.Context, accountFlow *model.AccountFlow) error {
	// 收入、支出、修改余额类型不应该出现在这里
	if accountFlow.Type == constant.AccountTypeIncome || accountFlow.Type == constant.AccountTypeExpense ||
		accountFlow.Type == constant.AccountTypeModify {
		return errno.New(errno.ErrIllegalOperate)
	}

	_, err := af.checkAccount(ctx, accountFlow.AccountId, accountFlow.UserId)
	if err != nil {
		return err
	}

	// 转入、转出校验目标账户
	if accountFlow.Type == constant.AccountTypeTransferIn || accountFlow.Type == constant.AccountTypeTransferOut {
		if accountFlow.AccountId == *accountFlow.TargetAccountId {
			return errno.New(errno.ErrIllegalOperate)
		}
		_, err = af.checkAccount(ctx, *accountFlow.TargetAccountId, accountFlow.UserId)
		if err != nil {
			return err
		}
	}

	// 借入、借出校验对方名称
	if accountFlow.Type == constant.AccountTypeBorrow || accountFlow.Type == constant.AccountTypeLend {
		if accountFlow.AssociateName == "" {
			return errno.New(errno.ErrAccountFlowAssociateNil)
		}
		finished := 0
		accountFlow.Finished = &finished
	}
	return nil
}

// QueryByAccountId 根据账户id查询其下的所有流水信息
func (af *accountFlowService) QueryByAccountId(ctx context.Context, accountId, userId int) ([]model.AccountFlow, error) {
	_, err := af.checkAccount(ctx, accountId, userId)
	if err != nil {
		return nil, err
	}
	return af.accountFlowStore.QueryByAccountId(ctx, accountId)
}
