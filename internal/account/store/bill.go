package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type BillStore interface {
	// Add 添加账单
	Add(ctx context.Context, bill *model.Bill) error

	// Update 更新账单信息
	Update(ctx context.Context, bill *model.Bill) error

	// Delete 删除账单
	Delete(ctx context.Context, billId int) error

	// QueryByTime 根据时间范围查询账单
	QueryByTime(ctx context.Context, bookId int, beginTime int64, endTime int64) ([]model.Bill, error)
	// QueryById 根据主键id查询
	QueryById(ctx context.Context, id int) (*model.Bill, error)

	// QueryBillTag 查询账单的标签备注
	QueryBillTag(ctx context.Context, bookId int) ([]model.BillTag, error)
}

type bill struct {
}

func NewBillStore() BillStore {
	return &bill{}
}

// Add 添加账单
func (b *bill) Add(ctx context.Context, bill *model.Bill) error {
	db := getDBFromContext(ctx)

	sql := "insert into bill(cost, type, remark, record_time, user_id, book_id, account_id, category_id, " +
		"create_time, update_time) " +
		"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	_, err := db.Exec(sql, bill.Cost, bill.Type, bill.Remark, bill.RecordTime, bill.UserId, bill.BookId,
		bill.AccountId, bill.CategoryId, bill.CreateTime, bill.UpdateTime)
	if err != nil {
		return errors.Wrap(err, "add bill store.")
	}

	return nil
}

// Update 更新账单信息
func (b *bill) Update(ctx context.Context, bill *model.Bill) error {
	db := getDBFromContext(ctx)

	sql := "update bill set cost=?,type=?,remark=?,record_time=?,user_id=?,account_id=?,category_id=?,update_time=? " +
		"where id = ?"

	_, err := db.Exec(sql, bill.Cost, bill.Type, bill.Remark, bill.RecordTime, bill.UserId, bill.AccountId,
		bill.CategoryId, bill.UpdateTime, bill.Id)
	if err != nil {
		return errors.Wrap(err, "update bill store.")
	}

	return nil
}

// Delete 删除账单
func (b *bill) Delete(ctx context.Context, billId int) error {
	db := getDBFromContext(ctx)

	sql := "delete from bill where id = ?"
	_, err := db.Exec(sql, billId)

	if err != nil {
		return errors.Wrap(err, "delete bill store.")
	}

	return nil
}

// QueryByTime 根据时间范围查询账单
func (b *bill) QueryByTime(ctx context.Context, bookId int, beginTime int64, endTime int64) ([]model.Bill, error) {
	db := getDBFromContext(ctx)

	sql := "select * from bill where book_id = ? and record_time >= ? and record_time <= ?"

	var bills = []model.Bill{}
	err := db.Select(&bills, sql, bookId, beginTime, endTime)
	if err != nil {
		return nil, errors.Wrap(err, "query bill store.")
	}

	return bills, nil
}

// QueryById 根据主键id查询
func (b *bill) QueryById(ctx context.Context, id int) (*model.Bill, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from bill where id = ?"

	var bill model.Bill
	err := db.Get(&bill, querySql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "query bill by id store.")
	}

	return &bill, nil
}

// QueryBillTag 查询账单的标签备注
func (b *bill) QueryBillTag(ctx context.Context, bookId int) ([]model.BillTag, error) {
	db := getDBFromContext(ctx)

	querySql := "SELECT category_id, group_concat(distinct remark) as remark from bill where book_id = ? " +
		"group by category_id order by record_time desc"

	var tags = []model.BillTag{}
	err := db.Select(&tags, querySql, bookId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "query bill tag store.")
	}

	return tags, nil
}
