package service

import (
	"context"
	"time"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
)

type BudgetSrv interface {
	SetBudget(ctx context.Context, budgetId, userId int, budget float32) error
	QueryBudget(ctx context.Context, bookId int) (*model.Budget, error)
}

type budgetService struct {
	budgetStore store.BudgetStore
}

func NewBudgetSrv() BudgetSrv {
	return &budgetService{budgetStore: store.NewBudgetStore()}
}

// SetBudget 设置账本预算
func (b *budgetService) SetBudget(ctx context.Context, budgetId, userId int, budget float32) error {
	return b.budgetStore.UpdateBudget(ctx, budgetId, budget, time.Now().Unix())
}

// QueryBudget 查询指定账本的预算
func (b *budgetService) QueryBudget(ctx context.Context, bookId int) (*model.Budget, error) {
	return b.budgetStore.QueryBudget(ctx, bookId)
}
