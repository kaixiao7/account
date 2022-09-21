package service

import (
	"context"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/errno"
)

type StatisticsSrv interface {
	QueryBill(ctx context.Context, bookId, userId int, beginTime, endTime int64) ([]model.Bill, error)
}

type statisticService struct {
	bookStore store.BookStore
	billStore store.BillStore
}

func NewStatisticsSrv() StatisticsSrv {
	return &statisticService{
		bookStore: store.NewBookStore(),
		billStore: store.NewBillStore(),
	}
}

func (s *statisticService) QueryBill(ctx context.Context, bookId, userId int, beginTime, endTime int64) ([]model.Bill, error) {
	if err := s.checkBookAndUser(ctx, bookId, userId); err != nil {
		return nil, err
	}

	return s.billStore.QueryByTime(ctx, bookId, beginTime, endTime)
}

func (s *statisticService) checkBookAndUser(ctx context.Context, bookId, userId int) error {
	members, err := s.bookStore.QueryBookMember(ctx, bookId)
	if err != nil {
		return err
	}
	for _, member := range members {
		if member == userId {
			return nil
		}
	}
	return errno.New(errno.ErrIllegalOperate)
}
