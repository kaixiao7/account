package service

import (
	"context"
	"time"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
	"kaixiao7/account/internal/pkg/errno"
	"kaixiao7/account/internal/pkg/log"
	"kaixiao7/account/internal/pkg/timex"
)

type BillSrv interface {
	// Add 添加账单
	Add(ctx context.Context, bill *model.Bill) error

	// Update 更新账单
	Update(ctx context.Context, bill *model.Bill) error

	// Delete 删除账单
	Delete(ctx context.Context, billId, userId, bookId int) error

	// QueryByTime 查询账本在指定月份的账单
	QueryByTime(ctx context.Context, bookId, userId int, date time.Time) ([]model.DayBill, error)

	// QueryTag 查询账单标签/备注
	QueryTag(ctx context.Context, bookId, userId int) ([]model.BillTag, error)
}

type billService struct {
	billStore     store.BillStore
	categoryStore store.CategoryStore
	bookStore     store.BookStore
	userStore     store.UserStore
	assetStore    store.AssetStore
}

func NewBillSrv() BillSrv {
	return &billService{
		billStore:     store.NewBillStore(),
		categoryStore: store.NewCategoryStore(),
		bookStore:     store.NewBookStore(),
		userStore:     store.NewUserStore(),
		assetStore:    store.NewAssetStore(),
	}
}

// QueryByTime 查询账本在指定月份的账单
func (b *billService) QueryByTime(ctx context.Context, bookId, userId int, date time.Time) ([]model.DayBill, error) {
	// 校验账本是否存在
	book, err := b.checkBook(ctx, bookId)
	if err != nil {
		return nil, err
	}
	// 校验用户是否是账本的成员
	if err := b.checkUserInBook(ctx, bookId, userId); err != nil {
		return nil, err
	}

	// 判断传入时间是否在账本创建时间之后，如果是，则返回没有更多数据了
	bookTime := time.Unix(book.CreateTime, 0)
	if date.Year() < bookTime.Year() || date.Month() < bookTime.Month() {
		return nil, errno.New(errno.ErrBillNotMore)
	}

	begin := timex.GetFirstDateOfMonth(date)
	end := timex.GetLastDateOfMonth(date)
	end = end.Add(time.Hour * 24)

	bills, err := b.billStore.QueryByTime(ctx, bookId, begin.Unix(), end.Unix())
	if err != nil {
		return nil, err
	}
	if len(bills) == 0 {
		return []model.DayBill{}, nil
	}

	// 日统计
	diffDays := int(end.Sub(begin).Hours() / 24)
	var dayBills []model.DayBill

	for i := diffDays - 1; i >= 0; i-- {
		start := begin.AddDate(0, 0, i)
		startTs := start.Unix()
		endTs := begin.AddDate(0, 0, i+1).Unix()
		var dayBill []model.Bill
		var income float64
		var expense float64
		for _, bill := range bills {
			if bill.RecordTime >= startTs && bill.RecordTime < endTs {
				if *bill.Type == constant.AssetTypeIncome {
					income = income + bill.Cost
				} else {
					expense = expense + bill.Cost
				}
				dayBill = append(dayBill, bill)
			}
		}

		if dayBill == nil {
			continue
		}

		dayBills = append(dayBills, model.DayBill{
			Date:    timex.JsonTime(start),
			Income:  income,
			Expense: expense,
			Bills:   dayBill,
		})
	}
	return dayBills, nil
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

	return b.billStore.QueryBillTag(ctx, bookId)
}

// Add 添加账单
func (b *billService) Add(ctx context.Context, bill *model.Bill) error {
	return b.save(ctx, bill)
}

// Update 更新账单
func (b *billService) Update(ctx context.Context, bill *model.Bill) error {
	return b.save(ctx, bill)
}

// Delete 删除账单
func (b *billService) Delete(ctx context.Context, billId, userId, bookId int) error {
	bill, err := b.queryBillById(ctx, billId)
	if err != nil {
		return err
	}
	// 只有账单的拥有者才可以删除
	if bill.UserId != userId || bill.BookId != bookId {
		return errno.New(errno.ErrBillNotDelete)
	}
	return WithTransaction(ctx, func(ctx context.Context) error {
		err = b.billStore.Delete(ctx, billId)
		if err != nil {
			return err
		}
		// 指定账户加上金额
		if bill.AssetId > 0 {
			if err := b.assetStore.ModifyBalance(ctx, bill.AssetId, bill.Cost); err != nil {
				return err
			}
		}
		return nil
	})
}

// 保存账单，根据 bill.Id 是否为0判断是新增还是更新
func (b *billService) save(ctx context.Context, bill *model.Bill) error {
	// 校验分类是否存在
	err := b.checkCategory(ctx, bill.CategoryId, bill.BookId)
	if err != nil {
		return err
	}

	// 校验账本是否存在
	_, err = b.checkBook(ctx, bill.BookId)
	if err != nil {
		return err
	}

	// 校验用户是否是账本的成员
	if err := b.checkUserInBook(ctx, bill.BookId, bill.UserId); err != nil {
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
	if bill.AssetId > 0 {
		if _, err := b.checkAsset(ctx, bill.AssetId, bill.UserId); err != nil {
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
			// 更新账单
			err = b.billStore.Update(ctx, bill)
		} else {
			// 插入账单
			err = b.billStore.Add(ctx, bill)
			diff = -bill.Cost
		}
		if err != nil {
			return err
		}

		// 账户加减指定金额
		if bill.AssetId > 0 {
			if err := b.assetStore.ModifyBalance(ctx, bill.AssetId, diff); err != nil {
				return err
			}
		}

		return nil
	})
}

func (b *billService) queryBillById(ctx context.Context, billId int) (*model.Bill, error) {
	bill, err := b.billStore.QueryById(ctx, billId)
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
	memberIds, err := b.bookStore.QueryBookMember(ctx, bookId)
	if err != nil {
		return err
	}
	for _, id := range memberIds {
		if id == userId {
			return nil
		}
	}

	log.Errorf("用户 %d 不是账本 %d 的成员", userId, bookId)
	return errno.New(errno.ErrIllegalOperate)
}

func (b *billService) checkAsset(ctx context.Context, assetId, userId int) (*model.Asset, error) {
	asset, err := b.assetStore.QueryById(ctx, assetId)
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
