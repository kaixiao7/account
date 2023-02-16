package service

import (
	"context"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/constant"
)

type BookSrv interface {
	Push(ctx context.Context, books []*model.Book, syncTime int64) error

	Pull(ctx context.Context, userId int64, lastSyncTime int64) ([]*model.Book, error)

	QueryBookList(ctx context.Context, userId int64) ([]*model.Book, error)
}

type bookService struct {
	bookStore store.BookStore
}

func NewBookSrv() BookSrv {
	return &bookService{bookStore: store.NewBookStore()}
}

func (b *bookService) Push(ctx context.Context, books []*model.Book, syncTime int64) error {
	return WithTransaction(ctx, func(ctx context.Context) error {
		for _, book := range books {
			book.SyncTime = syncTime
			if book.SyncState == constant.SYNC_ADD {
				book.SyncState = constant.SYNC_SUCCESS
				if e := b.bookStore.Add(ctx, book); e != nil {
					return e
				}
			} else {
				book.SyncState = constant.SYNC_SUCCESS
				if e := b.bookStore.Update(ctx, book); e != nil {
					return e
				}
			}
		}

		return nil
	})
}

func (b *bookService) Pull(ctx context.Context, userId int64, lastSyncTime int64) ([]*model.Book, error) {
	return b.bookStore.QueryBySyncTime(ctx, userId, lastSyncTime)
}

// QueryBookList 查询用户账本列表
func (b *bookService) QueryBookList(ctx context.Context, userId int64) ([]*model.Book, error) {
	return b.bookStore.QueryBookList(ctx, userId)
}
