package service

import (
	"context"
	"fmt"
	"sort"
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
	// QueryByAssetIdAndTime 根据资产id查询距离指定时间之前最近一个月的数据
	QueryByAssetIdAndTime(ctx context.Context, assetId, userId int, date time.Time) ([]model.AssetFlowVO, error)
}

type assetFlowService struct {
	assetStore     store.AssetStore
	assetFlowStore store.AssetFlowStore
	billStore      store.BillStore
}

func NewAssertFlowSrv() AssetFlowSrv {
	return &assetFlowService{
		assetStore:     store.NewAssetStore(),
		assetFlowStore: store.NewAssetFlowStore(),
		billStore:      store.NewBillStore(),
	}
}

// Add 添加流水
func (af *assetFlowService) Add(ctx context.Context, assetFlow *model.AssetFlow) error {
	if err := af.saveCheck(ctx, assetFlow); err != nil {
		return err
	}

	return WithTransaction(ctx, func(ctx context.Context) error {
		diff := assetFlow.Cost
		// 转出、借出、还款 将金额变为负数，因为修改账户余额中的sql计算方式为加法
		if assetFlow.Type == constant.AssetTypeTransferOut || assetFlow.Type == constant.AssetTypeBorrowOut ||
			assetFlow.Type == constant.AssetTypeStill {
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
// 其实流水不应该被修改，只能增加、删除
func (af *assetFlowService) Update(ctx context.Context, assetFlow *model.AssetFlow) error {
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
		// 转入、借入、收款 将金额变为负数
		if assetFlow.Type == constant.AssetTypeTransferIn || assetFlow.Type == constant.AssetTypeBorrowIn ||
			assetFlow.Type == constant.AssetTypeHarvest {
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
func (af *assetFlowService) Delete(ctx context.Context, assertFlowId, userId int) error {
	assetFlow, err := af.checkAssetFlow(ctx, assertFlowId, userId)
	if err != nil {
		return err
	}
	// 修改账户余额类型不允许删除
	if assetFlow.Type == constant.AssetTypeModify {
		return errno.New(errno.ErrIllegalOperate)
	}

	return WithTransaction(ctx, func(ctx context.Context) error {
		// 账户金额恢复
		if err := af.moneyRegain(ctx, assetFlow); err != nil {
			return err
		}

		// 删除流水记录
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
func (af *assetFlowService) transferReverse(ctx context.Context, assetFlow model.AssetFlow) error {
	diff := assetFlow.Cost
	// 反转账户id
	assetId := assetFlow.AssetId
	assetFlow.AssetId = *assetFlow.TargetAssetId
	assetFlow.TargetAssetId = &assetId

	if assetFlow.Type == constant.AssetTypeTransferOut {
		assetFlow.Type = constant.AssetTypeTransferIn
	} else if assetFlow.Type == constant.AssetTypeTransferIn {
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
func (af *assetFlowService) transferReverseQuery(ctx context.Context, assetFlow *model.AssetFlow) (*model.AssetFlow, error) {
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
func (af *assetFlowService) moneyRegain(ctx context.Context, assetFlow *model.AssetFlow) error {
	cost := assetFlow.Cost
	// 转入、借入、收款 将金额变为负数
	if assetFlow.Type == constant.AssetTypeTransferIn || assetFlow.Type == constant.AssetTypeBorrowIn ||
		assetFlow.Type == constant.AssetTypeHarvest {
		cost = -cost
	}
	// 账户余额恢复
	if err := af.assetStore.ModifyBalance(ctx, assetFlow.AssetId, cost); err != nil {
		return err
	}

	return nil
}

func (af *assetFlowService) checkAsset(ctx context.Context, assetId, userId int) (*model.Asset, error) {
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

func (af *assetFlowService) checkAssetFlow(ctx context.Context, assetFlowId, userId int) (*model.AssetFlow, error) {
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
func (af *assetFlowService) saveCheck(ctx context.Context, assetFlow *model.AssetFlow) error {
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
		if assetFlow.AssetId == *assetFlow.TargetAssetId {
			return errno.New(errno.ErrIllegalOperate)
		}
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
func (af *assetFlowService) QueryByAssetId(ctx context.Context, assetId, userId int) ([]model.AssetFlow, error) {
	_, err := af.checkAsset(ctx, assetId, userId)
	if err != nil {
		return nil, err
	}
	return af.assetFlowStore.QueryByUserIdAndAssetId(ctx, userId, assetId)
}

func (af *assetFlowService) getQueryDate(assetCreateTs int64, pageNum int) (time.Time, int, error) {
	now := time.Now()
	assetTime := time.Unix(assetCreateTs, 0)
	total := timex.CalcMonths(assetTime, now)
	if pageNum > total {
		return now, 0, errno.New(errno.ErrTokenInvalid)
	}

	queryDate := timex.SubMonth(now, pageNum)

	return queryDate, total, nil
}

// QueryByAssetIdAndTime 根据资产id查询距离指定时间之前最近一个月的数据
// 因为指定时间的月份有可能没有数据
func (af *assetFlowService) QueryByAssetIdAndTime(ctx context.Context, assetId, userId int, date time.Time) ([]model.AssetFlowVO, error) {
	asset, err := af.checkAsset(ctx, assetId, userId)
	if err != nil {
		return nil, err
	}
	// 判断传入时间是否在资产账户创建时间之后，如果是，则返回没有更多数据了
	assetTime := time.Unix(asset.CreateTime, 0)
	if date.Year() < assetTime.Year() || date.Month() < assetTime.Month() {
		return nil, errno.New(errno.ErrBillNotMore)
	}

	end := timex.GetLastDateTimeOfMonth(date)

	// 查询流水距离end之前的最近一条记录，获取其record_time
	firstFlow, err := af.assetFlowStore.QueryOneByUserIdAndAssetIdAndTime(ctx, userId, assetId, end.Unix())
	if err != nil {
		return nil, err
	}
	// 查询账单距离end之前的最近一条记录，获取其record_time
	firstBill, err := af.billStore.QueryOneByAssetIdAndTime(ctx, assetId, end.Unix())
	if err != nil {
		return nil, err
	}

	// 没有更多了
	if firstFlow == nil && firstBill == nil {
		return nil, errno.New(errno.ErrBillNotMore)
	}

	flowDate := assetTime
	if firstFlow != nil {
		flowDate = time.Unix(firstFlow.RecordTime, 0)
	}
	billDate := assetTime
	if firstBill != nil {
		billDate = time.Unix(firstBill.RecordTime, 0)
	}

	var ret []model.AssetFlowVO
	if timex.IsSameMonth(flowDate, billDate) || flowDate.After(billDate) {
		queryBegin := timex.GetFirstDateOfMonth(flowDate).Unix()
		queryEnd := timex.GetLastDateTimeOfMonth(flowDate).Unix()
		// 查询流水
		flows, err := af.assetFlowStore.QueryByUserIdAndAssetIdAndTime(ctx, userId, assetId, queryBegin, queryEnd)
		if err != nil {
			return nil, err
		}
		ret = append(ret, model.AssetFlow2VO(flows)...)
	}
	if timex.IsSameMonth(flowDate, billDate) || billDate.After(flowDate) {
		queryBegin := timex.GetFirstDateOfMonth(billDate).Unix()
		queryEnd := timex.GetLastDateTimeOfMonth(billDate).Unix()
		// 查询账户下的账单记录
		bills, err := af.billStore.QueryByAssetIdAndTime(ctx, assetId, queryBegin, queryEnd)
		if err != nil {
			return nil, err
		}
		ret = append(ret, model.Bill2VO(bills)...)
	}

	// 根据记录时间倒序排序
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].RecordTime > ret[j].RecordTime
	})

	return ret, nil
}
