package service

import (
	"context"
	"fmt"
	"time"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/timex"

	"github.com/pkg/errors"
)

type AssetFlowSrv interface {
	// Add 添加流水
	Add(ctx context.Context, assetFlow *model.AssetFlow) error
	// Update 修改流水
	// 仅修改基本信息，不会修改类型及账户
	Update(ctx context.Context, assetFlow *model.AssetFlow) error
	// Delete 删除流水
	Delete(ctx context.Context, assertFlowId, userId int) error

	// QueryByAssetId 根据账户id查询其下的所有流水信息
	QueryByAssetId(ctx context.Context, assetId, userId int) ([]model.AssetFlow, error)
	// QueryByAssetIdAndTime 根据资产id查询时间范围内的流水和账单
	QueryByAssetIdAndTime(ctx context.Context, assetId, userId int, date time.Time) ([]model.AssetFlowVO, error)
	// QueryBorrowIn 查询借入记录
	QueryBorrowIn(ctx context.Context, userId int) ([]model.AssetFlow, error)
	// QueryBorrowOut 查询借出记录
	QueryBorrowOut(ctx context.Context, userId int) ([]model.AssetFlow, error)
}

type assertFlowService struct {
	assetStore     store.AssetStore
	assetFlowStore store.AssetFlowStore
	billStore      store.BillStore
}

func NewAssertFlowSrv() AssetFlowSrv {
	return &assertFlowService{
		assetStore:     store.NewAssetStore(),
		assetFlowStore: store.NewAssetFlowStore(),
		billStore:      store.NewBillStore(),
	}
}

// Add 添加流水
func (af *assertFlowService) Add(ctx context.Context, assetFlow *model.AssetFlow) error {
	if err := af.saveCheck(ctx, assetFlow); err != nil {
		return err
	}

	return WithTransaction(ctx, func(ctx context.Context) error {
		diff := assetFlow.Cost
		// 转出、借出 将金额变为负数
		if assetFlow.Type == constant.AssetTypeTransferOut || assetFlow.Type == constant.AssetTypeBorrowOut {
			diff = -diff
		}
		// 修改账户余额
		if err := af.assetStore.ModifyBalance(ctx, assetFlow.AssetId, diff); err != nil {
			return err
		}
		// 插入流水记录
		if err := af.assetFlowStore.Add(ctx, assetFlow); err != nil {
			return err
		}

		// 转入、转出插入反转流水
		if assetFlow.Type == constant.AssetTypeTransferIn || assetFlow.Type == constant.AssetTypeTransferOut {
			if err := af.transferReverse(ctx, *assetFlow); err != nil {
				return err
			}
		}

		return nil
	})
}

// Update 修改流水
// 仅修改基本信息，不会修改类型及账户
func (af *assertFlowService) Update(ctx context.Context, assetFlow *model.AssetFlow) error {
	// 仅修改基本信息，不会修改类型及账户
	assetFlowBefore, err := af.checkAssetFlow(ctx, assetFlow.Id, assetFlow.UserId)
	if err != nil {
		return err
	}

	if assetFlowBefore.Type != assetFlow.Type || assetFlowBefore.AssetId != assetFlow.AssetId {
		return errno.New(errno.ErrIllegalOperate)
	}

	if err := af.saveCheck(ctx, assetFlow); err != nil {
		return err
	}

	// 前后差值
	diff := assetFlowBefore.Cost - assetFlow.Cost
	return WithTransaction(ctx, func(ctx context.Context) error {
		// 转入、借入 将金额变为负数
		if assetFlow.Type == constant.AssetTypeTransferIn || assetFlow.Type == constant.AssetTypeBorrowIn {
			diff = -diff
		}
		// 修改账户金额
		if err := af.assetStore.ModifyBalance(ctx, assetFlow.AssetId, diff); err != nil {
			return err
		}
		// 修改流水信息
		if err := af.assetFlowStore.Update(ctx, assetFlow); err != nil {
			return err
		}

		// 转入、转出反转流水更新
		if assetFlow.Type == constant.AssetTypeTransferIn || assetFlow.Type == constant.AssetTypeTransferOut {
			// 修改目标账户金额
			if err := af.assetStore.ModifyBalance(ctx, *assetFlow.TargetAssetId, -diff); err != nil {
				return err
			}
			// 查询反转流水信息并修改
			reverseFlow, err := af.transferReverseQuery(ctx, assetFlowBefore)
			if err != nil {
				return err
			}
			if reverseFlow == nil {
				return errors.New(fmt.Sprintf("未查询到反转流水, asset flow id: %d", assetFlow.Id))
			}
			reverseFlow.Cost = assetFlow.Cost
			reverseFlow.RecordTime = assetFlow.RecordTime
			reverseFlow.Remark = assetFlow.Remark
			if err = af.assetFlowStore.Update(ctx, reverseFlow); err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete 删除流水
func (af *assertFlowService) Delete(ctx context.Context, assertFlowId, userId int) error {
	assetFlow, err := af.checkAssetFlow(ctx, assertFlowId, userId)
	if err != nil {
		return err
	}
	// 修改账户余额类型不允许删除
	if assetFlow.Type == constant.AssetTypeModify {
		return errno.New(errno.ErrIllegalOperate)
	}

	return WithTransaction(ctx, func(ctx context.Context) error {
		if err := af.moneyRegain(ctx, assetFlow); err != nil {
			return err
		}

		// 删除记录
		if err := af.assetFlowStore.Delete(ctx, assertFlowId); err != nil {
			return err
		}

		// 转入、转出反转流水删除
		if assetFlow.Type == constant.AssetTypeTransferIn || assetFlow.Type == constant.AssetTypeTransferOut {
			reverseFlow, err := af.transferReverseQuery(ctx, assetFlow)
			if err != nil {
				return err
			}
			if reverseFlow == nil {
				return errors.New(fmt.Sprintf("未查询到反转流水: asset flow id: %d", assetFlow.Id))
			}

			// 金额恢复
			if err := af.moneyRegain(ctx, reverseFlow); err != nil {
				return err
			}

			// 删除记录
			if err := af.assetFlowStore.Delete(ctx, reverseFlow.Id); err != nil {
				return err
			}
		}

		return nil
	})
}

// 转账操作的反转插入
// 如果是转入：那么向对方账户插入转出记录
// 如果是转出：那么向对方账户插入转入记录
func (af *assertFlowService) transferReverse(ctx context.Context, assetFlow model.AssetFlow) error {
	diff := assetFlow.Cost
	// 反转账户id
	assetId := assetFlow.AssetId
	assetFlow.AssetId = *assetFlow.TargetAssetId
	assetFlow.TargetAssetId = &assetId

	if assetFlow.Type == constant.AssetTypeTransferOut {
		assetFlow.Type = constant.AssetTypeTransferIn
	}

	if assetFlow.Type == constant.AssetTypeTransferIn {
		assetFlow.Type = constant.AssetTypeTransferOut
		diff = -diff
	}

	// 插入流水
	if err := af.assetFlowStore.Add(ctx, &assetFlow); err != nil {
		return err
	}

	// 修改账户余额
	if err := af.assetStore.ModifyBalance(ctx, assetFlow.AssetId, diff); err != nil {
		return err
	}

	return nil
}

// 查询转入、转出反转流水
func (af *assertFlowService) transferReverseQuery(ctx context.Context, assetFlow *model.AssetFlow) (*model.AssetFlow, error) {
	flowType := constant.AssetTypeTransferIn
	if assetFlow.Type == constant.AssetTypeTransferIn {
		flowType = constant.AssetTypeTransferOut
	}
	if assetFlow.Type == constant.AssetTypeTransferOut {
		flowType = constant.AssetTypeTransferIn
	}

	ret, err := af.assetFlowStore.QueryReverseFlow(ctx, *assetFlow.TargetAssetId, assetFlow.AssetId, flowType,
		assetFlow.RecordTime)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// 账户余额恢复
func (af *assertFlowService) moneyRegain(ctx context.Context, assetFlow *model.AssetFlow) error {
	cost := assetFlow.Cost
	// 转入、借入 将金额变为负数
	if assetFlow.Type == constant.AssetTypeTransferIn || assetFlow.Type == constant.AssetTypeBorrowIn {
		cost = -cost
	}
	// 账户余额恢复
	if err := af.assetStore.ModifyBalance(ctx, assetFlow.AssetId, cost); err != nil {
		return err
	}

	return nil
}

func (af *assertFlowService) checkAsset(ctx context.Context, assetId, userId int) (*model.Asset, error) {
	asset, err := af.assetStore.QueryById(ctx, assetId)
	if err != nil {
		return nil, err
	}

	if asset == nil {
		return nil, errno.New(errno.ErrAssetNotFound)
	}

	if asset.UserId != userId {
		return nil, errno.New(errno.ErrIllegalOperate)
	}

	return asset, nil
}

func (af *assertFlowService) checkAssetFlow(ctx context.Context, assetFlowId, userId int) (*model.AssetFlow, error) {
	assetFlow, err := af.assetFlowStore.QueryById(ctx, assetFlowId)
	if err != nil {
		return nil, err
	}

	if assetFlow == nil {
		return nil, errno.New(errno.ErrAssetFlowNotFound)
	}

	if assetFlow.UserId != userId {
		return nil, errno.New(errno.ErrIllegalOperate)
	}

	return assetFlow, nil
}

// 插入、更新操作的前置校验
func (af *assertFlowService) saveCheck(ctx context.Context, assetFlow *model.AssetFlow) error {
	// 收入、支出、修改余额类型不应该出现在这里
	if assetFlow.Type == constant.AssetTypeIncome || assetFlow.Type == constant.AssetTypeExpense ||
		assetFlow.Type == constant.AssetTypeModify {
		return errno.New(errno.ErrIllegalOperate)
	}

	_, err := af.checkAsset(ctx, assetFlow.AssetId, assetFlow.UserId)
	if err != nil {
		return err
	}

	// 转入、转出校验目标账户
	if assetFlow.Type == constant.AssetTypeTransferIn || assetFlow.Type == constant.AssetTypeTransferOut {
		_, err = af.checkAsset(ctx, *assetFlow.TargetAssetId, assetFlow.UserId)
		if err != nil {
			return err
		}
	}

	// 借入、借出校验对方名称
	if assetFlow.Type == constant.AssetTypeBorrowIn || assetFlow.Type == constant.AssetTypeBorrowOut {
		if assetFlow.AssociateName == "" {
			return errno.New(errno.ErrAssetFlowAssociateNil)
		}
		finished := 0
		assetFlow.Finished = &finished
	}
	return nil
}

// QueryByAssetId 根据账户id查询其下的所有流水信息
func (af *assertFlowService) QueryByAssetId(ctx context.Context, assetId, userId int) ([]model.AssetFlow, error) {
	_, err := af.checkAsset(ctx, assetId, userId)
	if err != nil {
		return nil, err
	}
	return af.assetFlowStore.QueryByUserIdAndAssetId(ctx, userId, assetId)
}

// QueryByAssetIdAndTime 根据资产id查询时间范围内的流水和账单
func (af *assertFlowService) QueryByAssetIdAndTime(ctx context.Context, assetId, userId int, date time.Time) ([]model.AssetFlowVO, error) {
	asset, err := af.checkAsset(ctx, assetId, userId)
	if err != nil {
		return nil, err
	}
	// 判断传入时间是否在资产账户创建时间之后，如果是，则返回没有更多数据了
	assetTime := time.Unix(asset.CreateTime, 0)
	if date.Year() < assetTime.Year() || date.Month() < assetTime.Month() {
		return nil, errno.New(errno.ErrBillNotMore)
	}

	begin := timex.GetFirstDateOfMonth(date)
	end := timex.GetLastDateOfMonth(date)

	var ret = []model.AssetFlowVO{}

	flows, err := af.assetFlowStore.QueryByUserIdAndAssetId(ctx, userId, assetId)
	if err != nil {
		return nil, err
	}

	ret = append(ret, model.AssetFlow2VO(flows)...)

	bills, err := af.billStore.QueryByAssetIdAndTime(ctx, assetId, begin.Unix(), end.Unix())
	if err != nil {
		return nil, err
	}
	ret = append(ret, model.Bill2VO(bills)...)

	return ret, nil
}

// QueryBorrowIn 查询借入记录
func (af *assertFlowService) QueryBorrowIn(ctx context.Context, userId int) ([]model.AssetFlow, error) {
	return af.queryBorrow(ctx, userId, constant.AssetTypeBorrowIn)
}

// QueryBorrowOut 查询借出记录
func (af *assertFlowService) QueryBorrowOut(ctx context.Context, userId int) ([]model.AssetFlow, error) {
	return af.queryBorrow(ctx, userId, constant.AssetTypeBorrowOut)
}

func (af *assertFlowService) queryBorrow(ctx context.Context, userId, borrowType int) ([]model.AssetFlow, error) {
	return af.assetFlowStore.QueryByUserIdAndType(ctx, userId, borrowType)
}
