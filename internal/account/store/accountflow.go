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

	// QueryByBookSyncTimeCount 根据同步时间查询账本的流水记录总数
	QueryByBookSyncTimeCount(ctx context.Context, bookId int64, syncTime int64) (int, error)
	// QueryByBookSyncTime 根据同步时间查询账本的流水记录
	QueryByBookSyncTime(ctx context.Context, bookId int64, syncTime int64, pageNum, pageSize int) ([]*model.AccountFlow, error)
	// QueryByUserIdSyncTime 根据用户id及同步时间查询流水记录，不包括账本的记录
	QueryByUserIdSyncTime(ctx context.Context, userId int64, syncTime int64) ([]*model.AccountFlow, error)
	// QueryByBookIdPull 根据账本id同步指定时间范围内的数据
	QueryByBookIdPull(ctx context.Context, bookId, startTime, endTime, syncTime int64) ([]*model.AccountFlow, error)

	// Delete 删除账户流水
	// 逻辑删除，将字段del置为1
	Delete(ctx context.Context, id int64) error
	// DeleteByAccountId 根据资产删除其对应的流水
	// 逻辑删除，将字段del置为1
	DeleteByAccountId(ctx context.Context, accountId int64) error

	// FinishedBorrow 结束债务
	FinishedBorrow(ctx context.Context, accountId int64) error

	// QueryById 根据id查询
	QueryById(ctx context.Context, id int64) (*model.AccountFlow, error)

	QueryByBookIdCount(ctx context.Context, bookId int64) (int, error)
	QueryByBookIdPage(ctx context.Context, bookId int64, pageNum, pageSize int) ([]model.AccountFlow, error)

	// QueryByBorrowLendId 根据借贷id查询流水
	QueryByBorrowLendId(ctx context.Context, borrowLendId int64) ([]model.AccountFlow, error)

	// QueryByUserIdAndType 根据userId与类型查询
	QueryByUserIdAndType(ctx context.Context, userId int64, blType int) ([]model.AccountFlow, error)

	// QueryBillTag 查询账单的标签备注
	QueryBillTag(ctx context.Context, bookId int64) ([]model.BillTag, error)

	QueryByAccountId(ctx context.Context, accountId int64) ([]model.AccountFlow, error)
}

type accountFlow struct {
}

func NewAccountFlowStore() AccountFlow {
	return &accountFlow{}
}

// Add 添加账户流水
func (af *accountFlow) Add(ctx context.Context, flow *model.AccountFlow) error {
	db := getDBFromContext(ctx)

	insertSql := db.Rebind(`insert into account_flow values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	_, err := db.Exec(insertSql, flow.Id, flow.UserId, flow.Username, flow.AccountId, flow.Type, flow.Cost, flow.RecordTime,
		flow.DelFlag, flow.BookId, flow.CategoryId, flow.Remark, flow.TargetAccountId, flow.AssociateName, flow.Finished,
		flow.BorrowLendId, flow.Profit, flow.Reimburse, flow.SyncState, flow.SyncTime, flow.CreateTime, flow.UpdateTime)

	if err != nil {
		return errors.Wrap(err, "account flow add store")
	}

	return nil
}

// Update 更新账户流水
func (af *accountFlow) Update(ctx context.Context, flow *model.AccountFlow) error {
	db := getDBFromContext(ctx)

	updateSql := db.Rebind(`update account_flow set user_id=?,username=?,account_id=?,type=?,cost=?,record_time=?,del_flag=?,
                        book_id=?,category_id=?,remark=?,target_account_id=?,associate_name=?,finished=?,
                        borrow_lend_id=?,profit=?,reimburse=?,sync_state=?,sync_time=?,create_time=?,update_time=? where id = ?`)
	_, err := db.Exec(updateSql, flow.UserId, flow.Username, flow.AccountId, flow.Type, flow.Cost, flow.RecordTime,
		flow.DelFlag, flow.BookId, flow.CategoryId, flow.Remark, flow.TargetAccountId, flow.AssociateName, flow.Finished,
		flow.BorrowLendId, flow.Profit, flow.Reimburse, flow.SyncState, flow.SyncTime, flow.CreateTime, flow.UpdateTime, flow.Id)
	if err != nil {
		return errors.Wrap(err, "account flow update store")
	}

	return nil
}

// QueryByBookSyncTimeCount 根据同步时间查询账本的流水记录总数
func (af *accountFlow) QueryByBookSyncTimeCount(ctx context.Context, bookId int64, syncTime int64) (int, error) {
	db := getDBFromContext(ctx)

	// 查询总记录数
	querySql := db.Rebind("select count(1) from account_flow where book_id = ? and sync_time > ?")
	var count int
	if err := db.Get(&count, querySql, bookId, syncTime); err != nil {
		return 0, errors.Wrap(err, "QueryByBookSyncTimeCount err.")
	}
	return count, nil
}

// QueryByBookIdPull 根据账本id同步指定时间范围内的数据
func (af *accountFlow) QueryByBookIdPull(ctx context.Context, bookId, startTime, endTime, syncTime int64) ([]*model.AccountFlow, error) {
	db := getDBFromContext(ctx)

	querySql := db.Rebind("select * from account_flow where book_id = ? and record_time >= ? and record_time < ? and sync_time > ?")

	var ret []*model.AccountFlow
	if err := db.Select(&ret, querySql, bookId, startTime, endTime, syncTime); err != nil {
		return nil, errors.Wrap(err, "QueryByBookSyncTime error.")
	}
	if ret == nil {
		ret = []*model.AccountFlow{}
	}

	return ret, nil
}

// QueryByBookSyncTime 根据同步时间分页查询账本的流水记录
func (af *accountFlow) QueryByBookSyncTime(ctx context.Context, bookId int64, syncTime int64, pageNum, pageSize int) ([]*model.AccountFlow, error) {
	db := getDBFromContext(ctx)

	querySql := db.Rebind("select * from account_flow where book_id = ? and sync_time > ? limit ? offset ?")

	var ret []*model.AccountFlow
	if err := db.Select(&ret, querySql, bookId, syncTime, pageSize, (pageNum-1)*pageSize); err != nil {
		return nil, errors.Wrap(err, "QueryByBookSyncTime error.")
	}
	if ret == nil {
		ret = []*model.AccountFlow{}
	}

	return ret, nil
}

// QueryByUserIdSyncTime 根据用户id及同步时间查询流水记录，不包括账本的记录
func (af *accountFlow) QueryByUserIdSyncTime(ctx context.Context, userId int64, syncTime int64) ([]*model.AccountFlow, error) {
	db := getDBFromContext(ctx)

	querySql := db.Rebind("select * from account_flow where user_id = ? and book_id = 0 and sync_time > ? ")

	var ret []*model.AccountFlow
	if err := db.Select(&ret, querySql, userId, syncTime); err != nil {
		return nil, errors.Wrap(err, "QueryByUserIdSyncTime error.")
	}
	if ret == nil {
		ret = []*model.AccountFlow{}
	}

	return ret, nil
}

// Delete 删除账户流水
// 逻辑删除，将字段del置为1
func (af *accountFlow) Delete(ctx context.Context, id int64) error {
	db := getDBFromContext(ctx)

	deleteSql := db.Rebind("update account_flow set del_flag = ? where id = ?")
	_, err := db.Exec(deleteSql, constant.DelTrue, id)
	if err != nil {
		return errors.Wrap(err, "delete account flow store")
	}
	return nil
}

// DeleteByAccountId 根据资产删除其对应的流水
// 逻辑删除，将字段del置为1
func (af *accountFlow) DeleteByAccountId(ctx context.Context, accountId int64) error {
	db := getDBFromContext(ctx)

	deleteSql := db.Rebind("update account_flow set del_flag = ? where account_id =?")
	_, err := db.Exec(deleteSql, constant.DelTrue, accountId)
	if err != nil {
		return errors.Wrap(err, "delete account flow store")
	}
	return nil
}

// QueryById 根据id查询
func (af *accountFlow) QueryById(ctx context.Context, id int64) (*model.AccountFlow, error) {
	db := getDBFromContext(ctx)

	querySql := db.Rebind("select * from account_flow where id = ? and del_flag =?")
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
func (af *accountFlow) QueryByBookIdCount(ctx context.Context, bookId int64) (int, error) {
	db := getDBFromContext(ctx)

	// 查询总记录数
	querySql := db.Rebind("select count(1) from account_flow where book_id = ? and del_flag = ?")
	var count int
	if err := db.Get(&count, querySql, bookId, constant.DelFalse); err != nil {
		return 0, errors.Wrap(err, "QueryByBookIdCount err.")
	}
	return count, nil
}

// QueryByBookIdPage 根据bookId分页查询
func (af *accountFlow) QueryByBookIdPage(ctx context.Context, bookId int64, pageNum, pageSize int) ([]model.AccountFlow, error) {
	db := getDBFromContext(ctx)

	querySql := db.Rebind("select * from account_flow where book_id = ? and del_flag = ? limit ?, ?")

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
func (af *accountFlow) QueryByAccountId(ctx context.Context, accountId int64) ([]model.AccountFlow, error) {
	db := getDBFromContext(ctx)

	querySql := db.Rebind("select * from account_flow where (account_id = ? or target_account_id = ?) and del_flag = ? ")

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
func (af *accountFlow) QueryByBorrowLendId(ctx context.Context, borrowLendId int64) ([]model.AccountFlow, error) {
	db := getDBFromContext(ctx)
	querySql := db.Rebind("select * from account_flow where borrow_lend_id = ? and del_flag = ?")

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
func (af *accountFlow) QueryByUserIdAndType(ctx context.Context, userId int64, blType int) ([]model.AccountFlow, error) {
	db := getDBFromContext(ctx)
	querySql := db.Rebind("select * from account_flow where user_id = ? and type = ? and del_flag = ?")

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
func (af *accountFlow) QueryBillTag(ctx context.Context, bookId int64) ([]model.BillTag, error) {
	db := getDBFromContext(ctx)

	t := strconv.Itoa(constant.AccountTypeExpense) + "," + strconv.Itoa(constant.AccountTypeIncome)
	querySql := db.Rebind("SELECT category_id, group_concat(distinct remark) as remark from account_flow where book_id = ? " +
		"and type in (" + t + ") group by category_id order by record_time desc")

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
func (af *accountFlow) FinishedBorrow(ctx context.Context, accountId int64) error {
	db := getDBFromContext(ctx)

	sql := db.Rebind("update account_flow set finished = ? where account_id = ?")
	_, err := db.Exec(sql, 1, accountId)
	if err != nil {
		return errors.Wrap(err, "finished borrowlend store")
	}
	return nil
}
