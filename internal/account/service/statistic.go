package service

//
// import (
// 	"context"
// 	"time"
//
// 	"kaixiao7/account/internal/account/model"
// 	"kaixiao7/account/internal/account/store"
// 	"kaixiao7/account/internal/pkg/constant"
// 	"kaixiao7/account/internal/pkg/errno"
// 	"kaixiao7/account/internal/pkg/timex"
// )
//
// type StatisticsSrv interface {
// 	QueryBill(ctx context.Context, bookId, userId int, beginTime, endTime int64) ([]model.Bill, error)
//
// 	// QueryByYear 根据指定年份查询账单统计信息
// 	QueryByYear(ctx context.Context, bookId, userId int, year int) (*model.YearStatistic, error)
//
// 	// StatisticByCategory 根据分类进行统计
// 	StatisticByCategory(ctx context.Context, bookId, userId int, beginTime, endTime time.Time) (model.CategoryStatistic, error)
// }
//
// type statisticService struct {
// 	bookStore store.BookStore
// 	billStore store.BillStore
// }
//
// func NewStatisticsSrv() StatisticsSrv {
// 	return &statisticService{
// 		bookStore: store.NewBookStore(),
// 		billStore: store.NewBillStore(),
// 	}
// }
//
// func (s *statisticService) QueryBill(ctx context.Context, bookId, userId int, beginTime, endTime int64) ([]model.Bill, error) {
// 	if err := s.checkBookAndUser(ctx, bookId, userId); err != nil {
// 		return nil, err
// 	}
//
// 	return s.billStore.QueryByTime(ctx, bookId, beginTime, endTime)
// }
//
// // QueryByYear 根据指定年份查询账单统计信息
// func (s *statisticService) QueryByYear(ctx context.Context, bookId, userId int, year int) (*model.YearStatistic, error) {
// 	if err := s.checkBookAndUser(ctx, bookId, userId); err != nil {
// 		return nil, err
// 	}
//
// 	count := 12
// 	now := time.Now()
// 	if year > now.Year() {
// 		return nil, errno.New(errno.ErrValidation)
// 	}
// 	if now.Year() == year {
// 		count = int(now.Month())
// 	}
//
// 	var monthCount []model.MonthCount
//
// 	for i := count; i > 0; i-- {
// 		t := timex.GetTime(year, i, 1)
// 		begin := timex.GetFirstDateOfMonth(t)
// 		end := timex.GetLastDateTimeOfMonth(t)
//
// 		monthBills, err := s.billStore.QueryByTime(ctx, bookId, begin.Unix(), end.Unix())
// 		if err != nil {
// 			return nil, err
// 		}
// 		m := model.MonthCount{
// 			Month: i,
// 		}
// 		for _, bill := range monthBills {
// 			if *bill.Type == constant.AccountTypeIncome {
// 				m.Income += bill.Cost
// 			} else if *bill.Type == constant.AccountTypeExpense {
// 				m.Expense += bill.Cost
// 			}
// 		}
// 		monthCount = append(monthCount, m)
// 	}
//
// 	ret := model.YearStatistic{
// 		Months: monthCount,
// 	}
// 	for _, m := range monthCount {
// 		ret.Income += m.Income
// 		ret.Expense += m.Expense
// 	}
//
// 	return &ret, nil
// }
//
// // StatisticByCategory 根据分类进行统计
// func (s *statisticService) StatisticByCategory(ctx context.Context, bookId, userId int, beginTime, endTime time.Time) (model.CategoryStatistic, error) {
// 	var ret = model.CategoryStatistic{}
// 	if err := s.checkBookAndUser(ctx, bookId, userId); err != nil {
// 		return ret, err
// 	}
//
// 	begin := timex.GetFirstDateOfMonth(beginTime)
// 	end := timex.GetLastDateTimeOfMonth(endTime)
//
// 	ss, err := s.billStore.StatisticByCategory(ctx, bookId, begin.Unix(), end.Unix())
// 	if err != nil {
// 		return ret, err
// 	}
//
// 	for _, cs := range ss {
// 		if cs.Type == constant.AccountTypeIncome {
// 			ret.IncomeCost += cs.Cost
// 		} else if cs.Type == constant.AccountTypeExpense {
// 			ret.ExpenseCost += cs.Cost
// 		}
// 	}
// 	var incomes []model.CategoryStatisticInfo
// 	var expenses []model.CategoryStatisticInfo
// 	for _, cs := range ss {
// 		if cs.Type == constant.AccountTypeIncome {
// 			incomes = append(incomes, model.CategoryStatisticInfo{
// 				CategoryId: cs.CategoryId,
// 				Cost:       cs.Cost,
// 				Count:      cs.Count,
// 				Percent:    cs.Cost / ret.IncomeCost,
// 			})
// 		} else if cs.Type == constant.AccountTypeExpense {
// 			expenses = append(expenses, model.CategoryStatisticInfo{
// 				CategoryId: cs.CategoryId,
// 				Cost:       cs.Cost,
// 				Count:      cs.Count,
// 				Percent:    cs.Cost / ret.ExpenseCost,
// 			})
// 		}
// 	}
//
// 	ret.Incomes = incomes
// 	ret.Expenses = expenses
//
// 	return ret, nil
// }
//
// func (s *statisticService) checkBookAndUser(ctx context.Context, bookId, userId int) error {
// 	members, err := s.bookStore.QueryBookMember(ctx, bookId)
// 	if err != nil {
// 		return err
// 	}
// 	for _, member := range members {
// 		if member == userId {
// 			return nil
// 		}
// 	}
// 	return errno.New(errno.ErrIllegalOperate)
// }
