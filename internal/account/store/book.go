package store

import (
	"context"
	"database/sql"

	"kaixiao7/account/internal/account/model"

	"github.com/pkg/errors"
)

type BookStore interface {
	QueryBookList(ctx context.Context, userId int) ([]*model.Book, error)

	// QueryById 根据主键id查询
	QueryById(ctx context.Context, id int) (*model.Book, error)

	// QueryBookMember 根据账本id查询该账本下的所有成员
	QueryBookMember(ctx context.Context, bookId int) ([]int, error)
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

// QueryById 根据主键id查询
func (b *book) QueryById(ctx context.Context, id int) (*model.Book, error) {
	db := getDBFromContext(ctx)

	querySql := "select * from account_book where id = ?"
	var book model.Book
	err := db.Get(&book, querySql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "query book by id store")
	}

	return &book, nil
}

// QueryBookMember 根据账本id查询该账本下的所有成员
func (b *book) QueryBookMember(ctx context.Context, bookId int) ([]int, error) {
	db := getDBFromContext(ctx)

	querySql := "select user_id from book_member where book_id = ?"
	var memberIds []int
	err := db.Select(&memberIds, querySql, bookId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "query book member store")
	}

	return memberIds, nil
}
