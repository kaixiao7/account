package service

import (
	"context"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
)

type BudgetSrv interface {
	Push(ctx context.Context, budgets []*model.Budget, syncTime int64) error
	Pull(ctx context.Context, bookId int, lastSyncTime int64) ([]*model.Budget, error)
}

type budgetService struct {
	budgetStore store.BudgetStore
}

func NewBudgetSrv() BudgetSrv {
	return &budgetService{budgetStore: store.NewBudgetStore()}
}

func (b *budgetService) Push(ctx context.Context, budgets []*model.Budget, syncTime int64) error {
	return WithTransaction(ctx, func(ctx context.Context) error {
		for _, budget := range budgets {
			budget.SyncTime = syncTime
			if budget.SyncState == constant.SYNC_ADD {
				budget.SyncState = constant.SYNC_SUCCESS
				if e := b.budgetStore.AddBudget(ctx, budget); e != nil {
					return e
				}
			} else {
				budget.SyncState = constant.SYNC_SUCCESS
				if e := b.budgetStore.UpdateBudget(ctx, budget); e != nil {
					return e
				}
			}
		}

		return nil
	})
}

func (b *budgetService) Pull(ctx context.Context, bookId int, lastSyncTime int64) ([]*model.Budget, error) {
	return b.budgetStore.QueryBySyncTime(ctx, bookId, lastSyncTime)
}

// SetBudget 设置账本预算
// func (b *budgetService) SetBudget(ctx context.Context, budgetId, userId int, budget float64) error {
// 	return b.budgetStore.UpdateBudget(ctx, budgetId, budget, time.Now().Unix())
// }

// QueryBudget 查询指定账本的预算
// func (b *budgetService) QueryBudget(ctx context.Context, bookId int) (*model.Budget, error) {
// 	return b.budgetStore.QueryBudget(ctx, bookId)
// }
