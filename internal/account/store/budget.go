package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type BudgetStore interface {
	AddBudget(ctx context.Context, budget *model.Budget) error
	QueryBudget(ctx context.Context, bookId int) (*model.Budget, error)
	UpdateBudget(ctx context.Context, budgetId int, budget float64, updateTime int64) error
}

type budget struct {
}

func NewBudgetStore() BudgetStore {
	return &budget{}
}

// AddBudget 添加账本总预算
func (b *budget) AddBudget(ctx context.Context, budget *model.Budget) error {
	db := getDBFromContext(ctx)
	sql := "insert into budget_setting(budget, book_id, type, category_id, create_id, create_time, update_time) values(?, ?, ?, ?, ?, ?, ?)"

	_, err := db.Exec(sql, budget.Budget, budget.BookId, budget.Type, budget.CategoryId, budget.CreateId, budget.CreateTime, budget.UpdateTime)
	if err != nil {
		return errors.Wrap(err, "add budget store")
	}

	return nil
}

// UpdateBudget 更新账本预算
func (b *budget) UpdateBudget(ctx context.Context, budgetId int, budget float64, updateTime int64) error {
	db := getDBFromContext(ctx)
	sql := "update budget_setting set budget=?, update_time=? where id = ?"

	_, err := db.Exec(sql, budget, updateTime, budgetId)
	if err != nil {
		return errors.Wrap(err, "update budget store")
	}

	return nil
}

// QueryBudget 查询账本总预算
func (b *budget) QueryBudget(ctx context.Context, bookId int) (*model.Budget, error) {
	db := getDBFromContext(ctx)
	querySql := "select * from budget_setting where book_id = ? and type = 0"

	var budget model.Budget
	err := db.Get(&budget, querySql, bookId)
	if err != nil {
		// 没有查询到记录，则返回nil
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "query budget store")
	}

	return &budget, nil
}
