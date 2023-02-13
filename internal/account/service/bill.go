package service

import (
	"context"
	"time"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/errno"
)

type BillSrv interface {
	// Add 添加账单
	Add(ctx context.Context, bill *model.AccountFlow) error

	// Update 更新账单
	Update(ctx context.Context, bill *model.AccountFlow) error

	// Delete 删除账单
	Delete(ctx context.Context, billId, userId, bookId int) error

	// QueryByPage 分页查询账单列表（支出、收入）
	QueryByPage(ctx context.Context, bookId, userId, pageSize, pageNum int) (model.Page, error)

	// QueryTag 查询账单标签/备注
	QueryTag(ctx context.Context, bookId, userId int) ([]model.BillTag, error)
}

type billService struct {
	categoryStore    store.CategoryStore
	bookStore        store.BookStore
	userStore        store.UserStore
	accountStore     store.AccountStore
	accountFlowStore store.AccountFlow
}

func NewBillSrv() BillSrv {
	return &billService{
		categoryStore:    store.NewCategoryStore(),
		bookStore:        store.NewBookStore(),
		userStore:        store.NewUserStore(),
		accountStore:     store.NewAccountStore(),
		accountFlowStore: store.NewAccountFlowStore(),
	}
}

// QueryByPage 分页查询账单列表（支出、收入）
func (b *billService) QueryByPage(ctx context.Context, bookId, userId, pageSize, pageNum int) (model.Page, error) {
	// 校验账本是否存在
	_, err := b.checkBook(ctx, bookId)
	if err != nil {
		return model.Page{}, err
	}
	// 校验用户是否是账本的成员
	if err = b.checkUserInBook(ctx, bookId, userId); err != nil {
		return model.Page{}, err
	}
	count, err := b.accountFlowStore.QueryByBookIdCount(ctx, bookId)
	if err != nil {
		return model.Page{}, err
	}
	// 计算总页数
	var total int
	if count%pageSize == 0 {
		total = count / pageSize
	} else {
		total = count/pageSize + 1
	}
	ret := model.Page{
		Total:    total,
		PageSize: pageSize,
		PageNum:  pageNum,
	}
	if total < pageNum {
		return ret, nil
	}

	result, err := b.accountFlowStore.QueryByBookIdPage(ctx, bookId, pageNum, pageSize)
	if err != nil {
		return model.Page{}, err
	}

	ret.Data = result

	return ret, nil
}

// QueryTag 查询账单标签/备注
func (b *billService) QueryTag(ctx context.Context, bookId, userId int) ([]model.BillTag, error) {
	// 校验账本是否存在
	_, err := b.checkBook(ctx, bookId)
	if err != nil {
		return nil, err
	}
	// 校验用户是否是账本的成员
	if err := b.checkUserInBook(ctx, bookId, userId); err != nil {
		return nil, err
	}

	return b.accountFlowStore.QueryBillTag(ctx, bookId)
}

// Add 添加账单
func (b *billService) Add(ctx context.Context, bill *model.AccountFlow) error {
	return b.save(ctx, bill)
}

// Update 更新账单
func (b *billService) Update(ctx context.Context, bill *model.AccountFlow) error {
	return b.save(ctx, bill)
}

// Delete 删除账单
func (b *billService) Delete(ctx context.Context, billId, userId, bookId int) error {
	bill, err := b.queryBillById(ctx, billId)
	if err != nil {
		return err
	}
	// 只有账单的拥有者才可以删除
	if bill.UserId != userId || *bill.BookId != bookId {
		return errno.New(errno.ErrBillNotDelete)
	}
	return WithTransaction(ctx, func(ctx context.Context) error {
		err = b.accountFlowStore.Delete(ctx, billId)
		if err != nil {
			return err
		}
		// 指定账户加上金额
		if bill.AccountId > 0 {
			cost := bill.Cost
			// 如果删除的是收入，则取反
			if bill.Type == constant.AccountTypeIncome {
				cost = -cost
			}
			if err := b.accountStore.ModifyBalance(ctx, bill.AccountId, cost); err != nil {
				return err
			}
		}
		return nil
	})
}

// 保存账单，根据 bill.Id 是否为0判断是新增还是更新
func (b *billService) save(ctx context.Context, bill *model.AccountFlow) error {
	// 校验分类是否存在
	err := b.checkCategory(ctx, *bill.CategoryId, *bill.BookId)
	if err != nil {
		return err
	}

	// 校验账本是否存在
	_, err = b.checkBook(ctx, *bill.BookId)
	if err != nil {
		return err
	}

	// 校验用户是否是账本的成员
	if err := b.checkUserInBook(ctx, *bill.BookId, bill.UserId); err != nil {
		return err
	}

	// 如果是修改账单，只能自己修改自己的
	if bill.Id > 0 {
		if billBefore, err := b.queryBillById(ctx, bill.Id); err != nil {
			return err
		} else if billBefore.UserId != bill.UserId {
			// 不允许修改他人账单
			return errno.New(errno.ErrBillNotModify)
		}
	}

	// 校验账户是否存在
	if bill.AccountId > 0 {
		if _, err := b.checkAccount(ctx, bill.AccountId, bill.UserId); err != nil {
			return err
		}
	}

	user, err := b.userStore.GetById(ctx, bill.UserId)
	if err != nil {
		return err
	}
	bill.Username = user.Username

	now := time.Now().Unix()
	bill.CreateTime = now
	bill.UpdateTime = now

	return WithTransaction(ctx, func(ctx context.Context) error {
		// 需要向账户加/减的金额
		var diff float64 = 0

		if bill.Id > 0 {
			// 查询更新之前的记录，用于计算差值
			billBefore, err := b.queryBillById(ctx, bill.Id)
			if err != nil {
				return err
			}
			diff = billBefore.Cost - bill.Cost
			// 修改时，收入为差值取反
			if bill.Type == constant.AccountTypeIncome {
				diff = -diff
			}
			// 更新账单
			err = b.accountFlowStore.Update(ctx, bill)
		} else {
			// 插入账单
			err = b.accountFlowStore.Add(ctx, bill)
			diff = bill.Cost
			// 支出
			if bill.Type == constant.AccountTypeExpense {
				diff = -diff
			}
		}
		if err != nil {
			return err
		}

		// 账户加减指定金额
		if bill.AccountId > 0 {

			if err := b.accountStore.ModifyBalance(ctx, bill.AccountId, diff); err != nil {
				return err
			}
		}

		return nil
	})
}

func (b *billService) queryBillById(ctx context.Context, billId int) (*model.AccountFlow, error) {
	bill, err := b.accountFlowStore.QueryById(ctx, billId)
	if err != nil {
		return nil, err
	}
	if bill == nil {
		return nil, errno.New(errno.ErrBillNotFound)
	}
	return bill, nil
}

// 校验分类是否存在
// 存在返回nil，否则返回具体错误信息
func (b *billService) checkCategory(ctx context.Context, categoryId, bookId int) error {
	category, err := b.categoryStore.QueryById(ctx, categoryId)
	if err != nil {
		return err
	}

	if category == nil || category.BookId != bookId {
		return errno.New(errno.ErrCategoryNotFound)
	}

	return nil
}

// 校验账本是否存在
// 存在返回账本信息，否则返回具体错误信息
func (b *billService) checkBook(ctx context.Context, bookId int) (*model.Book, error) {
	book, err := b.bookStore.QueryById(ctx, bookId)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, errno.New(errno.ErrBookNotFound)
	}

	return book, nil
}

// 校验账本成员是否包含指定用户
// 如果包含，则返回nil，否则返回具体错误信息
func (b *billService) checkUserInBook(ctx context.Context, bookId, userId int) error {
	// memberIds, err := b.bookStore.QueryBookMember(ctx, bookId)
	// if err != nil {
	// 	return err
	// }
	// for _, id := range memberIds {
	// 	if id == userId {
	// 		return nil
	// 	}
	// }
	//
	// log.Errorf("用户 %d 不是账本 %d 的成员", userId, bookId)
	return errno.New(errno.ErrIllegalOperate)
}

func (b *billService) checkAccount(ctx context.Context, accountId, userId int) (*model.Account, error) {
	account, err := b.accountStore.QueryById(ctx, accountId)
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
