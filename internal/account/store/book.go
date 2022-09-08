package store

import (
	"context"
	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type BookStore interface {
	QueryBookList(ctx context.Context, userId int) ([]*model.Book, error)
}

type book struct {
}

func NewBookStore() BookStore {
	return &book{}
}

// QueryBookList 查询用户账本列表
func (b *book) QueryBookList(ctx context.Context, userId int) ([]*model.Book, error) {
	db := getDBFromContext(ctx)
	sql := `
		select *
		from account_book
		where del_flag = 0
		and id in (
			select book_id
			from book_member
			where user_id = ?
		)
	`
	var bookList = []*model.Book{}
	err := db.Select(&bookList, sql, userId)

	if err != nil {
		return nil, errors.Wrap(err, "query book list store")
	}

	return bookList, nil
}
