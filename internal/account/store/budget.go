package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type BudgetStore interface {
	AddBudget(ctx context.Context, budget *model.Budget) error
	UpdateBudget(ctx context.Context, budget *model.Budget) error
	QueryBySyncTime(ctx context.Context, bookId int64, syncTime int64) ([]*model.Budget, error)
}

type budget struct {
}

func NewBudgetStore() BudgetStore {
	return &budget{}
}

// AddBudget 添加账本总预算
func (b *budget) AddBudget(ctx context.Context, budget *model.Budget) error {
	db := getDBFromContext(ctx)
	sql := db.Rebind(`insert into book_budget(id,budget, book_id, type, category_id, create_id, sync_state, sync_time, create_time,
				update_time) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	_, err := db.Exec(sql, budget.Id, budget.Budget, budget.BookId, budget.Type, budget.CategoryId, budget.CreateId,
		budget.SyncState, budget.SyncTime, budget.CreateTime, budget.UpdateTime)
	if err != nil {
		return errors.Wrap(err, "add budget store")
	}

	return nil
}

// UpdateBudget 更新账本预算
func (b *budget) UpdateBudget(ctx context.Context, budget *model.Budget) error {
	db := getDBFromContext(ctx)
	sql := db.Rebind(`update book_budget set budget=?, book_id=?, type=?, category_id=?, create_id=?, sync_state=?, sync_time=?,
                       update_time=? where id = ?`)

	_, err := db.Exec(sql, budget.Budget, budget.BookId, budget.Type, budget.CategoryId, budget.CreateId,
		budget.SyncState, budget.SyncTime, budget.UpdateTime, budget.Id)
	if err != nil {
		return errors.Wrap(err, "update budget store")
	}

	return nil
}

func (b *budget) QueryBySyncTime(ctx context.Context, bookId int64, syncTime int64) ([]*model.Budget, error) {
	db := getDBFromContext(ctx)
	querySql := db.Rebind("select * from book_budget where book_id = ? and sync_time > ?")

	var budgetList = []*model.Budget{}
	err := db.Select(&budgetList, querySql, bookId, syncTime)
	if err != nil {
		// 没有查询到记录，则返回nil
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "query budget sync time store")
	}

	return budgetList, nil
}

// QueryByBookId 查询账本预算
func (b *budget) QueryByBookId(ctx context.Context, bookId int64) ([]*model.Budget, error) {
	db := getDBFromContext(ctx)
	querySql := db.Rebind("select * from book_budget where book_id = ?")

	var budgetList = []*model.Budget{}
	err := db.Select(&budgetList, querySql, bookId)
	if err != nil {
		// 没有查询到记录，则返回nil
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, errors.Wrap(err, "query budget store")
	}

	return budgetList, nil
}
