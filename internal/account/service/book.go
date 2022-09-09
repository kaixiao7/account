package service

import (
	"context"

	"kaixiao7/account/internal/account/model"
	"kaixiao7/account/internal/account/store"
)

type BookSrv interface {
	QueryBookList(ctx context.Context, userId int) ([]*model.Book, error)
}

type bookService struct {
	bookStore store.BookStore
}

func NewBookSrv() BookSrv {
	return &bookService{bookStore: store.NewBookStore()}
}

// QueryBookList 查询用户账本列表
func (b *bookService) QueryBookList(ctx context.Context, userId int) ([]*model.Book, error) {
	return b.bookStore.QueryBookList(ctx, userId)
}
