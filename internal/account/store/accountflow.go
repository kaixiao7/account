package store

import (
	"context"
	"database/sql"
	"strconv"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/pkg/constant"

	"github.com/pkg/errors"
)

// AccountFlow 账户流水
type AccountFlow interface {
	// Add 添加账户流水
	Add(ctx context.Context, flow *model.AccountFlow) error
	// Update 更新账户流水
	Update(ctx context.Context, flow *model.AccountFlow) error
	// Delete 删除账户流水
	// 逻辑删除，将字段del置为1
	Delete(ctx context.Context, id int) error
	// DeleteByAccountId 根据资产删除其对应的流水
	// 逻辑删除，将字段del置为1
	DeleteByAccountId(ctx context.Context, accountId int) error

	// FinishedBorrow 结束债务
	FinishedBorrow(ctx context.Context, accountId int) error

	// QueryById 根据id查询
	QueryById(ctx context.Context, id int) (*model.AccountFlow, error)

	QueryByBookIdCount(ctx context.Context, bookId int) (int, error)
	QueryByBookIdPage(ctx context.Context, bookId, pageNum, pageSize int) ([]model.AccountFlow, error)

	// QueryByBorrowLendId 根据借贷id查询流水
	QueryByBorrowLendId(ctx context.Context, borrowLendId int) ([]model.AccountFlow, error)

	// QueryByUserIdAndType 根据userId与类型查询
	QueryByUserIdAndType(ctx context.Context, userId, blType int) ([]model.AccountFlow, error)

	// QueryBillTag 查询账单的标签备注
	QueryBillTag(ctx context.Context, bookId int) ([]model.BillTag, error)

	QueryByAccountId(ctx context.Context, accountId int) ([]model.AccountFlow, error)
}

type accountFlow struct {
}

func NewAccountFlowStore() AccountFlow {
	return &accountFlow{}
}

// Add 添加账户流水
func (af *accountFlow) Add(ctx context.Context, flow *model.AccountFlow) error {
	db := getDBFromContext(ctx)

	field := "user_id, username, account_id, type, cost, record_time, remark, associate_name, create_time, update_time"
	values := "?,?,?,?,?,?,?,?,?,?"
	v := []any{flow.UserId, flow.Username, flow.AccountId, flow.Type, flow.Cost, flow.RecordTime, flow.Remark, flow.AssociateName,
		flow.CreateTime, flow.UpdateTime}

	if flow.BookId != nil {
		field = field + ", book_id"
		values = values + ", ?"
		v = append(v, flow.BookId)
	}
	if flow.CategoryId != nil {
		field = field + ", category_id"
		values = values + ", ?"
		v = append(v, flow.CategoryId)
	}
	if flow.TargetAccountId != nil {
		field = field + ", target_account_id"
		values = values + ", ?"
		v = append(v, flow.TargetAccountId)
	}
	if flow.Finished != nil {
		field = field + ", finished"
		values = values + ", ?"
		v = append(v, flow.Finished)
	}
	if flow.BorrowLendId != nil {
		field = field + ", borrow_lend_id"
		values = values + ", ?"
		v = append(v, flow.BorrowLendId)
	}
	if flow.Profit != nil {
		field = field + ", profit"
		values = values + ", ?"
		v = append(v, flow.Profit)
	}

	insertSql := "insert into account_flow(" + field + ") values(" + values + ")"
	_, err := db.Exec(insertSql, v...)
	if err != nil {
		return errors.Wrap(err, "account flow add store")
	}

	return nil
}

// Update 更新账户流水
func (af *accountFlow) Update(ctx context.Context, flow *model.AccountFlow) error {
	db := getDBFromContext(ctx)

	field := " user_id=?, username=?, account_id=?, type=?, cost=?, record_time=?, remark=?, associate_name=?, create_time=?, update_time=?"
	v := []any{flow.UserId, flow.Username, flow.AccountId, flow.Type, flow.Cost, flow.RecordTime, flow.Remark, flow.AssociateName,
		flow.CreateTime, flow.UpdateTime}

	if flow.BookId != nil {
		field = field + ", book_id=?"
		v = append(v, flow.BookId)
	}
	if flow.CategoryId != nil {
		field = field + ", category_id=?"
		v = append(v, flow.CategoryId)
	}
	if flow.TargetAccountId != nil {
		field = field + ", target_account_id=?"
		v = append(v, flow.TargetAccountId)
	}
	if flow.Finished != nil {
		field = field + ", finished=?"
		v = append(v, flow.Finished)
	}
	if flow.BorrowLendId != nil {
		field = field + ", borrow_id=?"
		v = append(v, flow.BorrowLendId)
	}
	if flow.Profit != nil {
		field = field + ", profit=?"
		v = append(v, flow.Profit)
	}
	v = append(v, flow.Id)

	updateSql := "update account_flow set" + field + " where id =?"
	_, err := db.Exec(updateSql, v...)
	if err != nil {
		return errors.Wrap(err, "account flow update store")
	}

	return nil
}

// Delete 删除账户流水
// 逻辑删除，将字段del置为1
func (af *accountFlow) Delete(ctx context.Context, id int) error {
	db := getDBFromContext(ctx)

	deleteSql := "update account_flow set del = ? where id = ?"
	_, err := db.Exec(deleteSql, constant.DelTrue, id)
	if err != nil {
		return errors.Wrap(err, "delete account flow store")
	}
	return nil
}

// DeleteByAccountId 根据资产删除其对应的流水
// 逻辑删除，将字段del置为1
func (af *accountFlow) DeleteByAccountId(ctx context.Context, accountId int) error {
	db := getDBFromContext(ctx)

	deleteSql := "update account_flow set del = ? where account_id =?"
	_, err := db.Exec(deleteSql, constant.DelTrue, accountId)
	if err != nil {
		return errors.Wrap(err, "delete account flow store")
	}
	return nil
}

// QueryById 根据id查询
func (af *accountFlow) QueryById(ctx context.Context, id int) (*model.AccountFlow, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from account_flow where id = ? and del =?"
	var accountFlow model.AccountFlow
	err := db.Get(&accountFlow, querySql, id, constant.DelFalse)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "account flow query by id store")
	}

	return &accountFlow, nil
}

// QueryByBookIdCount 根据bookId分页查询的总记录数
func (af *accountFlow) QueryByBookIdCount(ctx context.Context, bookId int) (int, error) {
	db := getDBFromContext(ctx)

	// 查询总记录数
	querySql := "select count(1) from account_flow where book_id = ? and del = ?"
	var count int
	if err := db.Get(&count, querySql, bookId, constant.DelFalse); err != nil {
		return 0, errors.Wrap(err, "QueryByBookIdCount err.")
	}
	return count, nil
}

// QueryByBookIdPage 根据bookId分页查询
func (af *accountFlow) QueryByBookIdPage(ctx context.Context, bookId, pageNum, pageSize int) ([]model.AccountFlow, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from account_flow where book_id = ? and del = ? limit ?, ?"

	var ret []model.AccountFlow
	if err := db.Select(&ret, querySql, bookId, constant.DelFalse, (pageNum-1)*pageSize, pageSize); err != nil {
		return nil, errors.Wrap(err, "QueryByBookIdPage error.")
	}
	if ret == nil {
		ret = []model.AccountFlow{}
	}

	return ret, nil
}

// QueryByAccountId 根据accountId查询
func (af *accountFlow) QueryByAccountId(ctx context.Context, accountId int) ([]model.AccountFlow, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from account_flow where (account_id = ? or target_account_id = ?) and del = ? "

	var ret []model.AccountFlow
	if err := db.Select(&ret, querySql, accountId, accountId, constant.DelFalse); err != nil {
		return nil, errors.Wrap(err, "QueryByAccountId error.")
	}
	if ret == nil {
		ret = []model.AccountFlow{}
	}

	return ret, nil
}

// QueryByBorrowLendId 根据借贷id查询流水
func (af *accountFlow) QueryByBorrowLendId(ctx context.Context, borrowLendId int) ([]model.AccountFlow, error) {
	db := getDBFromContext(ctx)
	querySql := "select * from account_flow where borrow_lend_id = ? and del = ?"

	var flows []model.AccountFlow
	if err := db.Select(&flows, querySql, borrowLendId, constant.DelFalse); err != nil {
		return nil, errors.Wrap(err, "QueryByBorrowLendId")
	}
	if flows == nil {
		flows = []model.AccountFlow{}
	}

	return flows, nil
}

// QueryByUserIdAndType 根据userId与类型查询
func (af *accountFlow) QueryByUserIdAndType(ctx context.Context, userId, blType int) ([]model.AccountFlow, error) {
	db := getDBFromContext(ctx)
	querySql := "select * from account_flow where user_id = ? and type = ? and del = ?"

	var flows []model.AccountFlow
	if err := db.Select(&flows, querySql, userId, blType, constant.DelFalse); err != nil {
		return nil, errors.Wrap(err, "QueryByUserIdAndType err.")
	}
	if flows == nil {
		flows = []model.AccountFlow{}
	}

	return flows, nil
}

// QueryBillTag 查询账单的标签备注
func (af *accountFlow) QueryBillTag(ctx context.Context, bookId int) ([]model.BillTag, error) {
	db := getDBFromContext(ctx)

	t := strconv.Itoa(constant.AccountTypeExpense) + "," + strconv.Itoa(constant.AccountTypeIncome)
	querySql := "SELECT category_id, group_concat(distinct remark) as remark from account_flow where book_id = ? " +
		"and type in (" + t + ") group by category_id order by record_time desc"

	var tags []model.BillTag
	err := db.Select(&tags, querySql, bookId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "query bill tag store.")
	}
	if tags == nil {
		tags = []model.BillTag{}
	}

	return tags, nil
}

// FinishedBorrow 结束债务
func (af *accountFlow) FinishedBorrow(ctx context.Context, accountId int) error {
	db := getDBFromContext(ctx)

	sql := "update account_flow set finished = ? where account_id = ?"
	_, err := db.Exec(sql, 1, accountId)
	if err != nil {
		return errors.Wrap(err, "finished borrowlend store")
	}
	return nil
}
